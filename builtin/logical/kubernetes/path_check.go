// SPDX-License-Identifier: MPL-2.0

package kubesecrets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/morevault/vaultum/sdk/v2/framework"
	"github.com/morevault/vaultum/sdk/v2/logical"
)

const (
	checkPath            = "check"
	checkHelpSynopsis    = `Checks the Kubernetes configuration is valid.`
	checkHelpDescription = `Checks the Kubernetes configuration is valid, checking if required environment variables are set.`
)

var envVarsToCheck = []string{k8sServiceHostEnv, k8sServicePortEnv}

func (b *backend) pathCheck() *framework.Path {
	return &framework.Path{
		Pattern: checkPath + "/?$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKubernetes,
			OperationVerb:   "check",
			OperationSuffix: "configuration",
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathCheckRead,
			},
		},
		HelpSynopsis:    checkHelpSynopsis,
		HelpDescription: checkHelpDescription,
	}
}

func (b *backend) pathCheckRead(_ context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	var missing []string
	for _, key := range envVarsToCheck {
		val := os.Getenv(key)
		if val == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) == 0 {
		return &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPStatusCode: http.StatusNoContent,
			},
		}, nil
	}

	missingText := strings.Join(missing, ", ")
	return logical.ErrorResponse(fmt.Sprintf("Missing environment variables: %s", missingText)), nil
}
