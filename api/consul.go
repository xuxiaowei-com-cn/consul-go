package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	userAgent = "consul-go"
)

type Client struct {
	client *retryablehttp.Client

	baseURL *url.URL

	username string
	password string

	UserAgent string

	Logger interface{} // 日志配置，优先级高于：Out、Prefix、Flag
	Out    io.Writer   // 日志配置
	Prefix string      // 日志配置
	Flag   int         // 日志配置

	Autopilot *AutopilotService
	Catalog   *CatalogService
	Kv        *KvService
}

type ListOptions struct {
	Page  int `url:"page,omitempty" json:"page,omitempty"`   // 当前的页码
	Limit int `url:"limit,omitempty" json:"limit,omitempty"` // 每页的数量
}

func NewClient(baseURL string, username string, password string) (*Client, error) {
	c := &Client{UserAgent: userAgent}
	client, err := BuildClient(c, baseURL, username, password)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// BuildClient
// retryablehttp.NewClient()
func BuildClient(c *Client, baseURL string, username string, password string) (*Client, error) {
	if c.Logger == nil {
		if c.Out == nil {
			c.Out = os.Stdout
		}
		if c.Flag == 0 {
			c.Flag = log.LstdFlags
		}
		c.Logger = log.New(c.Out, c.Prefix, c.Flag)
	}

	c.client = &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultPooledClient(),
		Logger:       c.Logger,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
		RetryMax:     4,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}

	c.Autopilot = &AutopilotService{client: c}
	c.Catalog = &CatalogService{client: c}
	c.Kv = &KvService{client: c}

	err := c.SetBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	c.username = username
	c.password = password

	return c, nil
}

func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

func (c *Client) SetBaseURL(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseURL = baseURL

	return nil
}

type RequestOptionFunc func(*retryablehttp.Request) error

func (c *Client) NewRequest(method string, path string, requestQuery interface{}, requestBody interface{}, options []RequestOptionFunc) (*retryablehttp.Request, error) {

	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	reqHeaders := make(http.Header)

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	var body interface{}
	switch {
	case method == http.MethodPatch || method == http.MethodPost || method == http.MethodPut:
		reqHeaders.Set("Content-Type", "application/json")

		if requestBody != nil {
			body, err = json.Marshal(requestBody)
			if err != nil {
				return nil, err
			}
		}
	}

	q, err := query.Values(requestQuery)
	if err != nil {
		return nil, err
	}

	u.RawQuery = q.Encode()

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	if c.username != "" && c.password != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))))
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (c *Client) UploadRequest(method, path string, content io.Reader, filename string, options []RequestOptionFunc) (*retryablehttp.Request, error) {
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(fw, content); err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	reqHeaders.Set("Content-Type", w.FormDataContentType())

	req, err := retryablehttp.NewRequest(method, u.String(), b)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	if c.username != "" && c.password != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))))
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Response struct {
	*http.Response
	body string
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// Do 发送请求，仅能处理 application/json、text/html 的响应
func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*Response, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			contentType := resp.Header.Get("Content-Type")

			if strings.Contains(contentType, "application/json") {
				err = json.NewDecoder(resp.Body).Decode(v)
			} else if strings.Contains(contentType, "text/html") {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return response, err
				}
				response.body = string(bodyBytes)
			} else {
				err = fmt.Errorf("unknown response type: %s", contentType)
			}
		}
	}

	return response, err
}

// DoDownload 下载
func (c *Client) DoDownload(req *retryablehttp.Request, destination string) (*Response, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	dir := filepath.Dir(destination)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return response, err
	}

	file, err := os.Create(destination)
	if err != nil {
		return response, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return response, err
	}

	return response, err
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = fmt.Sprintf("failed to parse unknown error format: %s", data)
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

func Getenv(key string, defaultValue string) string {
	var value = os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
