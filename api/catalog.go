package api

import "net/http"

type CatalogService struct {
	client *Client
}

func (s *CatalogService) GetDatacenters(options ...RequestOptionFunc) ([]string, *Response, error) {

	u := "v1/catalog/datacenters"

	req, err := s.client.NewRequest(http.MethodGet, u, nil, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var datacenters []string
	resp, err := s.client.Do(req, &datacenters)
	if err != nil {
		return nil, resp, err
	}

	return datacenters, resp, nil
}
