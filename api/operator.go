package api

import "net/http"

type AutopilotService struct {
	client *Client
}

func (s *AutopilotService) GetAutopilotState(options ...RequestOptionFunc) (*interface{}, *Response, error) {

	u := "v1/operator/autopilot/state"

	req, err := s.client.NewRequest(http.MethodGet, u, nil, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var state *interface{}
	resp, err := s.client.Do(req, &state)
	if err != nil {
		return nil, resp, err
	}

	return state, resp, nil
}
