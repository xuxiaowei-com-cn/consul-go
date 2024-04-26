package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetDatacenters(t *testing.T) {

	var baseURL = Getenv("CONSUL_GO_BASE_URL", "http://127.0.0.1:8500/")

	client, err := NewClient(baseURL, "", "")
	assert.NoError(t, err)

	datacenters, response, err := client.Catalog.GetDatacenters()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, len(datacenters) > 0)

	t.Log("Datacenters:", datacenters)
}
