package errhandler

import (
	clientv1 "github.com/keepcalmist/chat-service/internal/server/server-client/v1"
	"github.com/keepcalmist/chat-service/pkg/pointer"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	return Response{Error: clientv1.Error{
		Code:    clientv1.ErrorCode(code),
		Details: pointer.PtrWithZeroAsNil(details),
		Message: msg,
	}}
}
