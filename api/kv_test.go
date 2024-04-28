package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
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

	contents, response, err := client.Kv.GetKv("", getKvRequestQuery)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contentsByte, err := json.Marshal(contents)
	assert.NoError(t, err)

	t.Log("Contents:", string(contentsByte))
}

func TestGetKvName(t *testing.T) {
	var baseURL = Getenv("CONSUL_GO_BASE_URL", "http://127.0.0.1:8500/")
	var dc = Getenv("CONSUL_GO_DC", "dc1")
	var name = Getenv("CONSUL_GO_KV_NAME", randString(6))

	client, err := NewClient(baseURL, "", "")
	assert.NoError(t, err)

	var getKvNameRequestQuery = &GetKvNameRequestQuery{
		Dc: dc,
	}

	contents, response, err := client.Kv.GetKvName(name, getKvNameRequestQuery)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	contentsByte, err := json.Marshal(contents)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(contentsByte))
}

func TestPutKvName(t *testing.T) {
	var baseURL = Getenv("CONSUL_GO_BASE_URL", "http://127.0.0.1:8500")
	var dc = Getenv("CONSUL_GO_DC", "dc1")
	var name = Getenv("CONSUL_GO_KV_NAME", randString(6))
	var requestBody = Getenv("CONSUL_GO_KV_NAME_CONTENT", randString(32))
	PutKvName(baseURL, dc, name, requestBody, t)
}

func PutKvName(baseURL string, dc string, name string, requestBody string, t *testing.T) {
	t.Logf("requestBody: %s", requestBody)

	client, err := NewClient(baseURL, "", "")
	assert.NoError(t, err)

	var putKvNameRequestQuery = &PutKvNameRequestQuery{
		Dc: dc,
	}

	result, response, err := client.Kv.PutKvName(name, putKvNameRequestQuery, requestBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, *result)

	var getKvNameRequestQuery = &GetKvNameRequestQuery{
		Dc: dc,
	}
	responses, response, err := client.Kv.GetKvName(name, getKvNameRequestQuery)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contentsByte, err := json.Marshal(responses)
	assert.NoError(t, err)

	t.Log("Contents:", string(contentsByte))

	assert.Equal(t, 1, len(responses))
	assert.NotEmpty(t, responses[0].Value)

	decodedBytes, err := base64.StdEncoding.DecodeString(responses[0].Value)
	assert.NoError(t, err)

	decodedString := string(decodedBytes)
	t.Log(decodedString)

	assert.Equal(t, requestBody, strings.Trim(decodedString, `"`))
}

func TestGetRecursion(t *testing.T) {

	var baseURL = Getenv("CONSUL_GO_BASE_URL", "http://127.0.0.1:8500/")
	var dc = Getenv("CONSUL_GO_DC", "dc1")

	client, err := NewClient(baseURL, "", "")
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		TestPutKvName(t)
	}

	for i := 0; i < 5; i++ {
		name := randString(6) + "/" + randString(6)
		requestBody := randString(32)
		PutKvName(baseURL, dc, name, requestBody, t)
	}

	for i := 0; i < 5; i++ {
		name := randString(6) + "/" + randString(6) + "/" + randString(6)
		requestBody := randString(32)
		PutKvName(baseURL, dc, name, requestBody, t)
	}

	var getKvRequestQuery = &GetKvRequestQuery{
		Keys:      "",
		Dc:        dc,
		Separator: "/",
	}

	contents, response, err := client.Kv.GetKv("", getKvRequestQuery)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contentsByte, err := json.Marshal(contents)
	assert.NoError(t, err)

	t.Log("Contents:", string(contentsByte))

	for _, name := range contents {

		if strings.HasSuffix(name, "/") {
			folder(dc, name, client, t)
		} else {
			getKvName(dc, name, client, t)
		}

	}

}

func folder(dc string, path string, client *Client, t *testing.T) {
	var getKvRequestQuery = &GetKvRequestQuery{
		Keys:      "",
		Dc:        dc,
		Separator: "/",
	}

	_, err := os.Stat("tmp")
	if os.IsNotExist(err) {
		err := os.Mkdir("tmp", 0755)
		assert.NoError(t, err)
	}

	getFolder := "tmp/get/"

	_, err = os.Stat(getFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(getFolder, 0755)
		assert.NoError(t, err)
	}

	_, err = os.Stat(getFolder + path)
	if os.IsNotExist(err) {
		err := os.Mkdir(getFolder+path, 0755)
		assert.NoError(t, err)
	}

	contents, response, err := client.Kv.GetKv(path, getKvRequestQuery)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contentsByte, err := json.Marshal(contents)
	assert.NoError(t, err)

	t.Log("Contents:", string(contentsByte))

	for _, name := range contents {

		if name == path {
			continue
		}

		if strings.HasSuffix(name, "/") {
			folder(dc, name, client, t)
		} else {
			getKvName(dc, name, client, t)
		}
	}
}

func getKvName(dc string, name string, client *Client, t *testing.T) {
	var getKvNameRequestQuery = &GetKvNameRequestQuery{
		Dc: dc,
	}
	responses, response, err := client.Kv.GetKvName(name, getKvNameRequestQuery)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	contentsByte, err := json.Marshal(responses)
	assert.NoError(t, err)

	t.Log("Contents:", string(contentsByte))

	assert.Equal(t, 1, len(responses))
	assert.NotEmpty(t, responses[0].Value)

	decodedBytes, err := base64.StdEncoding.DecodeString(responses[0].Value)
	assert.NoError(t, err)

	decodedString := string(decodedBytes)
	t.Log(decodedString)

	err = os.WriteFile("tmp/get/"+name, []byte(strings.Trim(decodedString, `"`)), 0644)
	assert.NoError(t, err)
}
