package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetKv(t *testing.T) {

	var baseURL = Getenv("CONSUL_GO_BASE_URL", "http://127.0.0.1:8500/")
	var dc = Getenv("CONSUL_GO_DC", "dc1")

	client, err := NewClient(baseURL, "", "")
	assert.NoError(t, err)

	var getKvRequestQuery = &GetKvRequestQuery{
		Keys:      "",
		Dc:        dc,
		Separator: "/",
	}

	contents, response, err := client.Kv.GetKv(getKvRequestQuery)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contentsByte, err := json.Marshal(contents)
	assert.NoError(t, err)

	t.Log("Contents:", string(contentsByte))
}
