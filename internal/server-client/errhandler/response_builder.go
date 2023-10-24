package errhandler

import (
	clientv1 "github.com/keepcalmist/chat-service/internal/server-client/v1"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	// FIXME
	return Response{}
}
