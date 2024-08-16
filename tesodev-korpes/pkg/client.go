package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type RestClient struct {
	Client *fasthttp.Client
}

func NewRestClient() *RestClient {
	return &RestClient{
		Client: &fasthttp.Client{},
	}
}

func (c *RestClient) DoGetRequest(URI string, respModel any, token string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(URI)
	req.Header.Set("Authentication", token)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.ProcessClientResponseData(req, resp, respModel)
	if err != nil {
		return err
	}
	return nil
}

// ProcessClientResponseData processes the response from the client and decodes it into respModel
func (c *RestClient) ProcessClientResponseData(req *fasthttp.Request, resp *fasthttp.Response, respModel any) error {
	if err := c.Client.Do(req, resp); err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("Customer not found, expected status code 200 but got %d", resp.StatusCode())
	}
	contentType := resp.Header.Peek("Content-Type")
	if bytes.Index(contentType, []byte("application/json")) != 0 {
		return fmt.Errorf("expected content type application/json but got %s", contentType)
	}

	body := resp.Body()
	reader := bytes.NewReader(body)
	err := json.NewDecoder(reader).Decode(respModel)
	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	return nil
}
