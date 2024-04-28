package api

import (
	"fmt"
	"net/http"
)

type KvService struct {
	client *Client
}

type GetKvRequestQuery struct {
	Keys      string `json:"keys,omitempty" url:"keys"`
	Dc        string `json:"dc,omitempty" url:"dc,omitempty"`
	Separator string `json:"separator,omitempty" url:"separator,omitempty"`
}

func (s *KvService) GetKv(requestQuery *GetKvRequestQuery, options ...RequestOptionFunc) ([]string, *Response, error) {

	u := "/v1/kv"

	req, err := s.client.NewRequest(http.MethodGet, u, requestQuery, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var contents []string
	resp, err := s.client.Do(req, &contents)
	if err != nil {
		return nil, resp, err
	}

	return contents, resp, nil
}

type GetKvNameRequestQuery struct {
	Dc string `json:"dc,omitempty" url:"dc"`
}

func (s *KvService) GetKvName(name string, requestQuery *GetKvNameRequestQuery, options ...RequestOptionFunc) ([]string, *Response, error) {

	u := fmt.Sprintf("/v1/kv/%s", name)

	req, err := s.client.NewRequest(http.MethodGet, u, requestQuery, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var contents []string
	resp, err := s.client.Do(req, &contents)
	if err != nil {
		return nil, resp, err
	}

	return contents, resp, nil
}
