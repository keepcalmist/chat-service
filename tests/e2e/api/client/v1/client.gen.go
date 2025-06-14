// Package apiclientv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package apiclientv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	ErrorCodeCreateChatError    ErrorCode = 1000
	ErrorCodeCreateProblemError ErrorCode = 1001
)

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode int

// GetHistoryRequest defines model for GetHistoryRequest.
type GetHistoryRequest struct {
	Cursor   *string `json:"cursor,omitempty"`
	PageSize *int    `json:"pageSize,omitempty"`
}

// GetHistoryResponse defines model for GetHistoryResponse.
type GetHistoryResponse struct {
	Data  *MessagesPage `json:"data,omitempty"`
	Error *Error        `json:"error,omitempty"`
}

// Message defines model for Message.
type Message struct {
	AuthorId   *types.UserID   `json:"authorId,omitempty"`
	Body       string          `json:"body"`
	CreatedAt  time.Time       `json:"createdAt"`
	Id         types.MessageID `json:"id"`
	IsBlocked  bool            `json:"isBlocked"`
	IsReceived bool            `json:"isReceived"`
	IsService  bool            `json:"isService"`
}

// MessageHeader defines model for MessageHeader.
type MessageHeader struct {
	AuthorId  *types.UserID   `json:"authorId,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessagesPage defines model for MessagesPage.
type MessagesPage struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next"`
}

// SendMessageRequest defines model for SendMessageRequest.
type SendMessageRequest struct {
	MessageBody string `json:"messageBody"`
}

// SendMessageResponse defines model for SendMessageResponse.
type SendMessageResponse struct {
	Data  *MessageHeader `json:"data,omitempty"`
	Error *Error         `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostGetHistoryParams defines parameters for PostGetHistory.
type PostGetHistoryParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostSendMessageParams defines parameters for PostSendMessage.
type PostSendMessageParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetHistoryJSONRequestBody defines body for PostGetHistory for application/json ContentType.
type PostGetHistoryJSONRequestBody = GetHistoryRequest

// PostSendMessageJSONRequestBody defines body for PostSendMessage for application/json ContentType.
type PostSendMessageJSONRequestBody = SendMessageRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostGetHistoryWithBody request with any body
	PostGetHistoryWithBody(ctx context.Context, params *PostGetHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostGetHistory(ctx context.Context, params *PostGetHistoryParams, body PostGetHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostSendMessageWithBody request with any body
	PostSendMessageWithBody(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostSendMessage(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostGetHistoryWithBody(ctx context.Context, params *PostGetHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostGetHistoryRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostGetHistory(ctx context.Context, params *PostGetHistoryParams, body PostGetHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostGetHistoryRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostSendMessageWithBody(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostSendMessageRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostSendMessage(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostSendMessageRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostGetHistoryRequest calls the generic PostGetHistory builder with application/json body
func NewPostGetHistoryRequest(server string, params *PostGetHistoryParams, body PostGetHistoryJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostGetHistoryRequestWithBody(server, params, "application/json", bodyReader)
}

// NewPostGetHistoryRequestWithBody generates requests for PostGetHistory with any type of body
func NewPostGetHistoryRequestWithBody(server string, params *PostGetHistoryParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/getHistory")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

// NewPostSendMessageRequest calls the generic PostSendMessage builder with application/json body
func NewPostSendMessageRequest(server string, params *PostSendMessageParams, body PostSendMessageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostSendMessageRequestWithBody(server, params, "application/json", bodyReader)
}

// NewPostSendMessageRequestWithBody generates requests for PostSendMessage with any type of body
func NewPostSendMessageRequestWithBody(server string, params *PostSendMessageParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/sendMessage")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, params.XRequestID)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Request-ID", headerParam0)

	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostGetHistoryWithBodyWithResponse request with any body
	PostGetHistoryWithBodyWithResponse(ctx context.Context, params *PostGetHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostGetHistoryResponse, error)

	PostGetHistoryWithResponse(ctx context.Context, params *PostGetHistoryParams, body PostGetHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*PostGetHistoryResponse, error)

	// PostSendMessageWithBodyWithResponse request with any body
	PostSendMessageWithBodyWithResponse(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error)

	PostSendMessageWithResponse(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error)
}

type PostGetHistoryResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GetHistoryResponse
}

// Status returns HTTPResponse.Status
func (r PostGetHistoryResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostGetHistoryResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostSendMessageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SendMessageResponse
}

// Status returns HTTPResponse.Status
func (r PostSendMessageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostSendMessageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostGetHistoryWithBodyWithResponse request with arbitrary body returning *PostGetHistoryResponse
func (c *ClientWithResponses) PostGetHistoryWithBodyWithResponse(ctx context.Context, params *PostGetHistoryParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostGetHistoryResponse, error) {
	rsp, err := c.PostGetHistoryWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostGetHistoryResponse(rsp)
}

func (c *ClientWithResponses) PostGetHistoryWithResponse(ctx context.Context, params *PostGetHistoryParams, body PostGetHistoryJSONRequestBody, reqEditors ...RequestEditorFn) (*PostGetHistoryResponse, error) {
	rsp, err := c.PostGetHistory(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostGetHistoryResponse(rsp)
}

// PostSendMessageWithBodyWithResponse request with arbitrary body returning *PostSendMessageResponse
func (c *ClientWithResponses) PostSendMessageWithBodyWithResponse(ctx context.Context, params *PostSendMessageParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error) {
	rsp, err := c.PostSendMessageWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostSendMessageResponse(rsp)
}

func (c *ClientWithResponses) PostSendMessageWithResponse(ctx context.Context, params *PostSendMessageParams, body PostSendMessageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostSendMessageResponse, error) {
	rsp, err := c.PostSendMessage(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostSendMessageResponse(rsp)
}

// ParsePostGetHistoryResponse parses an HTTP response from a PostGetHistoryWithResponse call
func ParsePostGetHistoryResponse(rsp *http.Response) (*PostGetHistoryResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostGetHistoryResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GetHistoryResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostSendMessageResponse parses an HTTP response from a PostSendMessageWithResponse call
func ParsePostSendMessageResponse(rsp *http.Response) (*PostSendMessageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostSendMessageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SendMessageResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
