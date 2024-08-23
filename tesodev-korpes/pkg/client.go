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

	if !isLoginEndpoint(URI) {
		req.Header.Set("Authentication", token)
	}

	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.ProcessClientResponseData(req, resp, respModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *RestClient) DoPostRequest(URI string, body any, respModel any, token string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(URI)

	if !isLoginEndpoint(URI) {
		req.Header.Set("Authentication", token)
	}

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	// Encode the body to JSON and set it in the request
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		req.SetBody(bodyBytes)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Eğer respModel nil ise, sadece isteği gönder ve yanıtı çözümleme
	if respModel == nil {
		if err := c.Client.Do(req, resp); err != nil {
			return fmt.Errorf("failed to perform request: %w", err)
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("expected status code 200 but got %d", resp.StatusCode())
		}
		return nil
	}

	// Aksi takdirde, yanıtı çözümlemek için ProcessClientResponseData kullan
	err := c.ProcessClientResponseData(req, resp, respModel)
	if err != nil {
		return err
	}
	return nil
}

// ProcessClientResponseData processes the response from the clientCon and decodes it into respModel
func (c *RestClient) ProcessClientResponseData(req *fasthttp.Request, resp *fasthttp.Response, respModel any) error {

	if err := c.Client.Do(req, resp); err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("not found, expected status code 200 but got %d", resp.StatusCode())
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

func isLoginEndpoint(URI string) bool {
	// Login endpointi ile eşleşen bir kontrol ekleyin. Örneğin:
	return URI == "http://localhost:8001/login"
}
