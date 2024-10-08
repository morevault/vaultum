// SPDX-License-Identifier: MPL-2.0

package pprof

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	vaulthttp "github.com/morevault/vaultum/http"
	"github.com/morevault/vaultum/internalshared/configutil"
	"github.com/morevault/vaultum/sdk/v2/helper/testhelpers/schema"
	"github.com/morevault/vaultum/vault"
	"golang.org/x/net/http2"
)

func TestSysPprof(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	SysPprof_Test(t, cluster)
}

func TestSysPprof_MaxRequestDuration(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = cluster.Cores[0].TLSConfig()
	if err := http2.ConfigureTransport(transport); err != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	sec := strconv.Itoa(int(vault.DefaultMaxRequestDuration.Seconds()) + 1)

	req := client.NewRequest("GET", "/v1/sys/pprof/profile")
	req.Params.Set("seconds", sec)
	httpReq, err := req.ToHTTP()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	httpRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	httpResp := make(map[string]interface{})

	// If we error here, it means that profiling likely happened, which is not
	// what we're checking for in this case.
	if err := json.Unmarshal(httpRespBody, &httpResp); err != nil {
		t.Fatalf("expected valid error response, got: %v", err)
	}

	errs, ok := httpResp["errors"].([]interface{})
	if !ok {
		t.Fatalf("expected error response, got: %v", httpResp)
	}
	if len(errs) == 0 || !strings.Contains(errs[0].(string), "exceeds max request duration") {
		t.Fatalf("unexpected error returned: %v", errs)
	}
}

func TestSysPprof_Standby(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		DisablePerformanceStandby: true,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{
				Profiling: configutil.ListenerProfiling{
					UnauthenticatedPProfAccess: true,
				},
			},
		},
	})
	defer cluster.Cleanup()

	SysPprof_Standby_Test(t, cluster)
}
