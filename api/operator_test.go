package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAutopilotState(t *testing.T) {

	var baseURL = Getenv("CONSUL_GO_BASE_URL", "http://127.0.0.1:8500/")

	client, err := NewClient(baseURL, "", "")
	assert.NoError(t, err)

	state, response, err := client.Autopilot.GetAutopilotState()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	stateByte, err := json.Marshal(state)
	assert.NoError(t, err)

	t.Log("State:", string(stateByte))
}
