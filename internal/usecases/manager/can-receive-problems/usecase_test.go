package canreceiveproblems_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/stretchr/testify/suite"

	"github.com/keepcalmist/chat-service/internal/testingh"
	canreceiveproblems "github.com/keepcalmist/chat-service/internal/usecases/manager/can-receive-problems"
	canreceiveproblemsmocks "github.com/keepcalmist/chat-service/internal/usecases/manager/can-receive-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *canreceiveproblemsmocks.MockmanagerLoadService
	mPoolMock *canreceiveproblemsmocks.MockmanagerPool
	uCase     canreceiveproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.mLoadMock = canreceiveproblemsmocks.NewMockmanagerLoadService(s.ctrl)
	s.mPoolMock = canreceiveproblemsmocks.NewMockmanagerPool(s.ctrl)

	uCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)
	s.uCase = uCase

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestUseCaseHandle_RequestValidation() {
	s.mPoolMock.EXPECT().Contains(s.Ctx, gomock.Any()).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, gomock.Any()).Return(true, nil)

	tCase := []struct {
		name    string
		request canreceiveproblems.Request
		wantErr bool
	}{
		{
			name: "valid",
			request: canreceiveproblems.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},
		{
			name: "empty request id",
			request: canreceiveproblems.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
			},
			wantErr: true,
		},
		{
			name: "empty manager id",
			request: canreceiveproblems.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
			},
			wantErr: true,
		},
		{
			name: "empty struct",
			request: canreceiveproblems.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.UserIDNil,
			},
			wantErr: true,
		},
	}

	for _, tc := range tCase {
		s.Run(tc.name, func() {
			resp, err := s.uCase.Handle(s.Ctx, tc.request)
			if tc.wantErr {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Equal(true, resp.Result)
		})
	}
}

func (s *UseCaseSuite) TestUseCaseHandle_ContainsError() {
	containsError := errors.New("contains error")

	s.mPoolMock.EXPECT().Contains(s.Ctx, gomock.Any()).Return(false, containsError)

	_, err := s.uCase.Handle(s.Ctx, canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	})

	s.ErrorIs(err, containsError)
}

func (s *UseCaseSuite) TestUseCaseHandle_ManagerAlreadyInThePool() {
	s.mPoolMock.EXPECT().Contains(s.Ctx, gomock.Any()).Return(true, nil)

	resp, err := s.uCase.Handle(s.Ctx, canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	})

	s.Require().NoError(err)
	s.False(resp.Result)
}

func (s *UseCaseSuite) TestUseCaseHandle_CanManagerTakeProblemError() {
	containsError := errors.New("contains error")

	s.mPoolMock.EXPECT().Contains(s.Ctx, gomock.Any()).Return(false, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, gomock.Any()).Return(false, containsError)

	_, err := s.uCase.Handle(s.Ctx, canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	})

	s.ErrorIs(err, containsError)
}

func (s *UseCaseSuite) TestUseCaseHandle_SuccessfulResults() {
	firstManager := types.NewUserID()
	secondManager := types.NewUserID()

	s.mPoolMock.EXPECT().Contains(s.Ctx, firstManager).Return(false, nil)
	s.mPoolMock.EXPECT().Contains(s.Ctx, secondManager).Return(false, nil)

	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, firstManager).Return(true, nil)
	s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, secondManager).Return(false, nil)

	resp, err := s.uCase.Handle(s.Ctx, canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: firstManager,
	})
	s.Require().NoError(err)
	s.True(resp.Result)

	resp, err = s.uCase.Handle(s.Ctx, canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: secondManager,
	})
	s.Require().NoError(err)
	s.False(resp.Result)
}
