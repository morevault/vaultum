// SPDX-License-Identifier: MPL-2.0

package http

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/morevault/vaultum/vault"
)

func TestSysLeader_get(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/leader")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"ha_enabled":             false,
		"is_self":                false,
		"leader_address":         "",
		"leader_cluster_address": "",
		"active_time":            time.Time{}.UTC().Format(time.RFC3339),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v \n%#v", actual, expected)
	}
}
