package shared

import (
	"context"
	"encoding/json"
	"fmt"
	config "github.com/BacoFoods/menu/internal"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"net/http"
)

type Request struct {
	Endpoint string
	Body     interface{}
	Query    map[string]string
	Headers  map[string]string
	Response interface{}
}

// RestClient contains methods for a rest client implementation.
type RestClient interface {
	Get(get Request) (*resty.Response, error)
	Post(post Request) (*resty.Response, error)
	Patch(patch Request) (*resty.Response, error)
	GetClient() *http.Client
}

// RestClient Implementation

const (
	LogHttpClient = "pkg/httpclient/http_client_impl"
)

type Rest struct {
	client *resty.Client
}

func NewRestClient(client *resty.Client) *Rest {
	environment := config.Config.AppEnv
	if environment == "development" {
		client.SetDebug(true)
	}
	client.OnAfterResponse(LogAfterResponse)
	client.OnBeforeRequest(LogBeforeRequest)
	return &Rest{client}
}

func (r *Rest) Get(get Request) (*resty.Response, error) {
	return r.client.R().
		SetHeaders(get.Headers).
		SetBody(get.Body).
		SetQueryParams(get.Query).
		SetResult(&get.Response).
		Get(get.Endpoint)
}

func (r *Rest) Post(post Request) (*resty.Response, error) {
	return r.client.R().
		SetResult(&post.Response).
		SetHeaders(post.Headers).
		SetQueryParams(post.Query).
		SetBody(post.Body).
		Post(post.Endpoint)
}

func (r *Rest) Patch(patch Request) (*resty.Response, error) {
	return r.client.R().
		SetResult(&patch.Response).
		SetHeaders(patch.Headers).
		SetQueryParams(patch.Query).
		SetBody(patch.Body).
		Patch(patch.Endpoint)
}

func (r *Rest) GetClient() *http.Client {
	return r.client.GetClient()
}

type ContextKey string

const (
	RequestUUID ContextKey = "request_uuid"
)

func LogAfterResponse(r *resty.Client, response *resty.Response) error {
	requestID := fmt.Sprintf("%v", response.Request.Context().Value(RequestUUID))
	LogResponse(requestID, response.Status(), string(response.Body()), response.Request.Method, response.Request.URL, response.Header())
	return nil
}

func LogBeforeRequest(_ *resty.Client, request *resty.Request) error {
	id, err := uuid.NewUUID()
	if err != nil {
		LogError("error uuid lib fail", LogHttpClient, "LogBeforeRequest", err)
	}

	ctx := context.WithValue(request.Context(), RequestUUID, id.String())
	request.SetContext(ctx)

	body, err := json.Marshal(request.Body)
	if err != nil {
		LogError("error json marshal fail", LogHttpClient, "LogBeforeRequest", err)
	}
	LogRequest(id.String(), request.Method, request.URL, string(body), request, request.Header)
	return nil
}
