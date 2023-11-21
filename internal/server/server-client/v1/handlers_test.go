package clientv1_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	clientv12 "github.com/keepcalmist/chat-service/internal/server/server-client/v1"
	"github.com/keepcalmist/chat-service/internal/server/server-client/v1/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/keepcalmist/chat-service/internal/middlewares"
	"github.com/keepcalmist/chat-service/internal/testingh"
	"github.com/keepcalmist/chat-service/internal/types"
)

type HandlersSuite struct {
	testingh.ContextSuite

	ctrl              *gomock.Controller
	getHistoryUseCase *clientv1mocks.MockgetHistoryUseCase
	sendMsgUseCase    *clientv1mocks.MocksendMessageUseCase
	handlers          clientv12.Handlers

	clientID types.UserID
}

func TestHandlersSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlersSuite))
}

func (s *HandlersSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.getHistoryUseCase = clientv1mocks.NewMockgetHistoryUseCase(s.ctrl)
	s.sendMsgUseCase = clientv1mocks.NewMocksendMessageUseCase(s.ctrl)
	{
		var err error
		s.handlers, err = clientv12.NewHandlers(clientv12.NewOptions(s.getHistoryUseCase, s.sendMsgUseCase))
		s.Require().NoError(err)
	}
	s.clientID = types.NewUserID()

	s.ContextSuite.SetupTest()
}

func (s *HandlersSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *HandlersSuite) newEchoCtx(
	requestID types.RequestID,
	path string,
	body string,
) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderXRequestID, requestID.String())

	resp := httptest.NewRecorder()

	ctx := echo.New().NewContext(req, resp)
	middlewares.SetToken(ctx, s.clientID)

	return resp, ctx
}
