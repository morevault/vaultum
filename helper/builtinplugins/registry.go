// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package builtinplugins

import (
	"context"

	credAppRole "github.com/morevault/vaultum/builtin/credential/approle"
	credCert "github.com/morevault/vaultum/builtin/credential/cert"
	credJWT "github.com/morevault/vaultum/builtin/credential/jwt"
	credKerb "github.com/morevault/vaultum/builtin/credential/kerberos"
	credKube "github.com/morevault/vaultum/builtin/credential/kubernetes"
	credLdap "github.com/morevault/vaultum/builtin/credential/ldap"
	credRadius "github.com/morevault/vaultum/builtin/credential/radius"
	credUserpass "github.com/morevault/vaultum/builtin/credential/userpass"
	logicalKube "github.com/morevault/vaultum/builtin/logical/kubernetes"
	logicalKv "github.com/morevault/vaultum/builtin/logical/kv"
	logicalLDAP "github.com/morevault/vaultum/builtin/logical/openldap"
	logicalPki "github.com/morevault/vaultum/builtin/logical/pki"
	logicalRabbit "github.com/morevault/vaultum/builtin/logical/rabbitmq"
	logicalSsh "github.com/morevault/vaultum/builtin/logical/ssh"
	logicalTotp "github.com/morevault/vaultum/builtin/logical/totp"
	logicalTransit "github.com/morevault/vaultum/builtin/logical/transit"
	dbCass "github.com/morevault/vaultum/plugins/database/cassandra"
	dbInflux "github.com/morevault/vaultum/plugins/database/influxdb"
	dbMysql "github.com/morevault/vaultum/plugins/database/mysql"
	dbPostgres "github.com/morevault/vaultum/plugins/database/postgresql"
	"github.com/morevault/vaultum/sdk/v2/framework"
	"github.com/morevault/vaultum/sdk/v2/helper/consts"
	"github.com/morevault/vaultum/sdk/v2/logical"
)

// Registry is inherently thread-safe because it's immutable.
// Thus, rather than creating multiple instances of it, we only need one.
var Registry = newRegistry()

var addExternalPlugins = addExtPluginsImpl

// BuiltinFactory is the func signature that should be returned by
// the plugin's New() func.
type BuiltinFactory func() (interface{}, error)

// There are three forms of Backends which exist in the BuiltinRegistry.
type credentialBackend struct {
	logical.Factory
	consts.DeprecationStatus
}

type databasePlugin struct {
	Factory BuiltinFactory
	consts.DeprecationStatus
}

type logicalBackend struct {
	logical.Factory
	consts.DeprecationStatus
}

type removedBackend struct {
	*framework.Backend
}

func removedFactory(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
	removedBackend := &removedBackend{}
	removedBackend.Backend = &framework.Backend{}
	return removedBackend, nil
}

func newRegistry() *registry {
	reg := &registry{
		credentialBackends: map[string]credentialBackend{
			"approle":    {Factory: credAppRole.Factory},
			"cert":       {Factory: credCert.Factory},
			"jwt":        {Factory: credJWT.Factory},
			"kerberos":   {Factory: credKerb.Factory},
			"kubernetes": {Factory: credKube.Factory},
			"ldap":       {Factory: credLdap.Factory},
			"oidc":       {Factory: credJWT.Factory},
			"radius":     {Factory: credRadius.Factory},
			"userpass":   {Factory: credUserpass.Factory},
		},
		databasePlugins: map[string]databasePlugin{
			// These four plugins all use the same mysql implementation but with
			// different username settings passed by the constructor.
			"mysql-database-plugin":        {Factory: dbMysql.New(dbMysql.DefaultUserNameTemplate)},
			"mysql-aurora-database-plugin": {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},
			"mysql-rds-database-plugin":    {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},
			"mysql-legacy-database-plugin": {Factory: dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate)},

			"cassandra-database-plugin":  {Factory: dbCass.New},
			"influxdb-database-plugin":   {Factory: dbInflux.New},
			"postgresql-database-plugin": {Factory: dbPostgres.New},
		},
		logicalBackends: map[string]logicalBackend{
			"kubernetes": {Factory: logicalKube.Factory},
			"kv":         {Factory: logicalKv.Factory},
			"openldap":   {Factory: logicalLDAP.Factory},
			"ldap":       {Factory: logicalLDAP.Factory},
			"pki":        {Factory: logicalPki.Factory},
			"rabbitmq":   {Factory: logicalRabbit.Factory},
			"ssh":        {Factory: logicalSsh.Factory},
			"totp":       {Factory: logicalTotp.Factory},
			"transit":    {Factory: logicalTransit.Factory},
		},
	}

	addExternalPlugins(reg)

	return reg
}

func addExtPluginsImpl(r *registry) {}

type registry struct {
	credentialBackends map[string]credentialBackend
	databasePlugins    map[string]databasePlugin
	logicalBackends    map[string]logicalBackend
}

// Get returns the Factory func for a particular backend plugin from the
// plugins map.
func (r *registry) Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool) {
	switch pluginType {
	case consts.PluginTypeCredential:
		if f, ok := r.credentialBackends[name]; ok {
			return toFunc(f.Factory), ok
		}
	case consts.PluginTypeSecrets:
		if f, ok := r.logicalBackends[name]; ok {
			return toFunc(f.Factory), ok
		}
	case consts.PluginTypeDatabase:
		if f, ok := r.databasePlugins[name]; ok {
			return f.Factory, ok
		}
	default:
		return nil, false
	}

	return nil, false
}

// Keys returns the list of plugin names that are considered builtin plugins.
func (r *registry) Keys(pluginType consts.PluginType) []string {
	var keys []string
	switch pluginType {
	case consts.PluginTypeDatabase:
		for key, backend := range r.databasePlugins {
			keys = appendIfNotRemoved(keys, key, backend.DeprecationStatus)
		}
	case consts.PluginTypeCredential:
		for key, backend := range r.credentialBackends {
			keys = appendIfNotRemoved(keys, key, backend.DeprecationStatus)
		}
	case consts.PluginTypeSecrets:
		for key, backend := range r.logicalBackends {
			keys = appendIfNotRemoved(keys, key, backend.DeprecationStatus)
		}
	}
	return keys
}

func (r *registry) Contains(name string, pluginType consts.PluginType) bool {
	for _, key := range r.Keys(pluginType) {
		if key == name {
			return true
		}
	}
	return false
}

// DeprecationStatus returns the Deprecation status for a builtin with type `pluginType`
func (r *registry) DeprecationStatus(name string, pluginType consts.PluginType) (consts.DeprecationStatus, bool) {
	switch pluginType {
	case consts.PluginTypeCredential:
		if f, ok := r.credentialBackends[name]; ok {
			return f.DeprecationStatus, ok
		}
	case consts.PluginTypeSecrets:
		if f, ok := r.logicalBackends[name]; ok {
			return f.DeprecationStatus, ok
		}
	case consts.PluginTypeDatabase:
		if f, ok := r.databasePlugins[name]; ok {
			return f.DeprecationStatus, ok
		}
	default:
		return consts.Unknown, false
	}

	return consts.Unknown, false
}

func toFunc(ifc interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		return ifc, nil
	}
}

func appendIfNotRemoved(keys []string, name string, status consts.DeprecationStatus) []string {
	if status != consts.Removed {
		return append(keys, name)
	}
	return keys
}
