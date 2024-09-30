
// SPDX-License-Identifier: MPL-2.0

package ssh

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/morevault/vaultum/sdk/v2/framework"
	"github.com/morevault/vaultum/sdk/v2/logical"
	"golang.org/x/crypto/ssh"

	"github.com/mikesmitty/edkey"
)

const (
	caPublicKey                       = "ca_public_key"
	caPrivateKey                      = "ca_private_key"
	caPublicKeyStoragePath            = "config/ca_public_key"
	caPublicKeyStoragePathDeprecated  = "public_key"
	caPrivateKeyStoragePath           = "config/ca_private_key"
	caPrivateKeyStoragePathDeprecated = "config/ca_bundle"
)

type keyStorageEntry struct {
	Key string `json:"key" structs:"key" mapstructure:"key"`
}

func pathConfigCA(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/ca",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
		},

		Fields: map[string]*framework.FieldSchema{
			"private_key": {
				Type:        framework.TypeString,
				Description: `Private half of the SSH key that will be used to sign certificates.`,
			},
			"public_key": {
				Type:        framework.TypeString,
				Description: `Public half of the SSH key that will be used to sign certificates.`,
			},
			"generate_signing_key": {
				Type:        framework.TypeBool,
				Description: `Generate SSH key pair internally rather than use the private_key and public_key fields.`,
				Default:     true,
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: `Specifies the desired key type when generating; could be a OpenSSH key type identifier (ssh-rsa, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, ecdsa-sha2-nistp521, or ssh-ed25519) or an algorithm (rsa, ec, ed25519).`,
				Default:     "ssh-rsa",
			},
			"key_bits": {
				Type:        framework.TypeInt,
				Description: `Specifies the desired key bits when generating variable-length keys (such as when key_type="ssh-rsa") or which NIST P-curve to use when key_type="ec" (256, 384, or 521).`,
				Default:     0,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigCAUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "ca",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigCADelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "ca-configuration",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigCARead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "ca-configuration",
				},
			},
		},

		HelpSynopsis: `Set the SSH private key used for signing certificates.`,
		HelpDescription: `This sets the CA information used for certificates generated by this
by this mount. The fields must be in the standard private and public SSH format.

For security reasons, the private key cannot be retrieved later.

Read operations will return the public key, if already stored/generated.`,
	}
}

func (b *backend) pathConfigCARead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	publicKeyEntry, err := caKey(ctx, req.Storage, caPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA public key: %w", err)
	}

	if publicKeyEntry == nil {
		return logical.ErrorResponse("keys haven't been configured yet"), nil
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			"public_key": publicKeyEntry.Key,
		},
	}

	return response, nil
}

func (b *backend) pathConfigCADelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, caPrivateKeyStoragePath); err != nil {
		return nil, err
	}
	if err := req.Storage.Delete(ctx, caPublicKeyStoragePath); err != nil {
		return nil, err
	}
	return nil, nil
}

func caKey(ctx context.Context, storage logical.Storage, keyType string) (*keyStorageEntry, error) {
	var path, deprecatedPath string
	switch keyType {
	case caPrivateKey:
		path = caPrivateKeyStoragePath
		deprecatedPath = caPrivateKeyStoragePathDeprecated
	case caPublicKey:
		path = caPublicKeyStoragePath
		deprecatedPath = caPublicKeyStoragePathDeprecated
	default:
		return nil, fmt.Errorf("unrecognized key type %q", keyType)
	}

	entry, err := storage.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA key of type %q: %w", keyType, err)
	}

	if entry == nil {
		// If the entry is not found, look at an older path. If found, upgrade
		// it.
		entry, err = storage.Get(ctx, deprecatedPath)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			entry, err = logical.StorageEntryJSON(path, keyStorageEntry{
				Key: string(entry.Value),
			})
			if err != nil {
				return nil, err
			}
			if err := storage.Put(ctx, entry); err != nil {
				return nil, err
			}
			if err = storage.Delete(ctx, deprecatedPath); err != nil {
				return nil, err
			}
		}
	}
	if entry == nil {
		return nil, nil
	}

	var keyEntry keyStorageEntry
	if err := entry.DecodeJSON(&keyEntry); err != nil {
		return nil, err
	}

	return &keyEntry, nil
}

func (b *backend) pathConfigCAUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error
	publicKey := data.Get("public_key").(string)
	privateKey := data.Get("private_key").(string)

	var generateSigningKey bool

	generateSigningKeyRaw, ok := data.GetOk("generate_signing_key")
	switch {
	// explicitly set true
	case ok && generateSigningKeyRaw.(bool):
		if publicKey != "" || privateKey != "" {
			return logical.ErrorResponse("public_key and private_key must not be set when generate_signing_key is set to true"), nil
		}

		generateSigningKey = true

	// explicitly set to false, or not set and we have both a public and private key
	case ok, publicKey != "" && privateKey != "":
		if publicKey == "" {
			return logical.ErrorResponse("missing public_key"), nil
		}

		if privateKey == "" {
			return logical.ErrorResponse("missing private_key"), nil
		}

		_, err := ssh.ParsePrivateKey([]byte(privateKey))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Unable to parse private_key as an SSH private key: %v", err)), nil
		}

		_, err = parsePublicSSHKey(publicKey)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Unable to parse public_key as an SSH public key: %v", err)), nil
		}

	// not set and no public/private key provided so generate
	case publicKey == "" && privateKey == "":
		generateSigningKey = true

	// not set, but one or the other supplied
	default:
		return logical.ErrorResponse("only one of public_key and private_key set; both must be set to use, or both must be blank to auto-generate"), nil
	}

	if generateSigningKey {
		keyType := data.Get("key_type").(string)
		keyBits := data.Get("key_bits").(int)

		publicKey, privateKey, err = generateSSHKeyPair(b.Backend.GetRandomReader(), keyType, keyBits)
		if err != nil {
			return nil, err
		}
	}

	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("failed to generate or parse the keys")
	}

	publicKeyEntry, err := caKey(ctx, req.Storage, caPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA public key: %w", err)
	}

	privateKeyEntry, err := caKey(ctx, req.Storage, caPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA private key: %w", err)
	}

	if (publicKeyEntry != nil && publicKeyEntry.Key != "") || (privateKeyEntry != nil && privateKeyEntry.Key != "") {
		return logical.ErrorResponse("keys are already configured; delete them before reconfiguring"), nil
	}

	entry, err := logical.StorageEntryJSON(caPublicKeyStoragePath, &keyStorageEntry{
		Key: publicKey,
	})
	if err != nil {
		return nil, err
	}

	// Save the public key
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	entry, err = logical.StorageEntryJSON(caPrivateKeyStoragePath, &keyStorageEntry{
		Key: privateKey,
	})
	if err != nil {
		return nil, err
	}

	// Save the private key
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		var mErr *multierror.Error

		mErr = multierror.Append(mErr, fmt.Errorf("failed to store CA private key: %w", err))

		// If storing private key fails, the corresponding public key should be
		// removed
		if delErr := req.Storage.Delete(ctx, caPublicKeyStoragePath); delErr != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("failed to cleanup CA public key: %w", delErr))
			return nil, mErr
		}

		return nil, err
	}

	if generateSigningKey {
		response := &logical.Response{
			Data: map[string]interface{}{
				"public_key": publicKey,
			},
		}

		return response, nil
	}

	return nil, nil
}

func generateSSHKeyPair(randomSource io.Reader, keyType string, keyBits int) (string, string, error) {
	if randomSource == nil {
		randomSource = rand.Reader
	}

	var publicKey crypto.PublicKey
	var privateBlock *pem.Block

	switch keyType {
	case ssh.KeyAlgoRSA, "rsa":
		if keyBits == 0 {
			keyBits = 4096
		}

		if keyBits < 2048 {
			return "", "", fmt.Errorf("refusing to generate weak %v key: %v bits < 2048 bits", keyType, keyBits)
		}

		privateSeed, err := rsa.GenerateKey(randomSource, keyBits)
		if err != nil {
			return "", "", err
		}

		privateBlock = &pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(privateSeed),
		}

		publicKey = privateSeed.Public()
	case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521, "ec":
		var curve elliptic.Curve
		switch keyType {
		case ssh.KeyAlgoECDSA256:
			curve = elliptic.P256()
		case ssh.KeyAlgoECDSA384:
			curve = elliptic.P384()
		case ssh.KeyAlgoECDSA521:
			curve = elliptic.P521()
		default:
			switch keyBits {
			case 0, 256:
				curve = elliptic.P256()
			case 384:
				curve = elliptic.P384()
			case 521:
				curve = elliptic.P521()
			default:
				return "", "", fmt.Errorf("unknown ECDSA key pair algorithm and bits: %v / %v", keyType, keyBits)
			}
		}

		privateSeed, err := ecdsa.GenerateKey(curve, randomSource)
		if err != nil {
			return "", "", err
		}

		marshalled, err := x509.MarshalECPrivateKey(privateSeed)
		if err != nil {
			return "", "", err
		}

		privateBlock = &pem.Block{
			Type:    "EC PRIVATE KEY",
			Headers: nil,
			Bytes:   marshalled,
		}

		publicKey = privateSeed.Public()
	case ssh.KeyAlgoED25519, "ed25519":
		_, privateSeed, err := ed25519.GenerateKey(randomSource)
		if err != nil {
			return "", "", err
		}

		marshalled := edkey.MarshalED25519PrivateKey(privateSeed)
		if marshalled == nil {
			return "", "", errors.New("unable to marshal ed25519 private key")
		}

		privateBlock = &pem.Block{
			Type:    "OPENSSH PRIVATE KEY",
			Headers: nil,
			Bytes:   marshalled,
		}

		publicKey = privateSeed.Public()
	default:
		return "", "", fmt.Errorf("unknown ssh key pair algorithm: %v", keyType)
	}

	public, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	return string(ssh.MarshalAuthorizedKey(public)), string(pem.EncodeToMemory(privateBlock)), nil
}
