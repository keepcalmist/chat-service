package freehands_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/keepcalmist/chat-service/internal/testingh"
	"github.com/keepcalmist/chat-service/internal/types"
	freehands "github.com/keepcalmist/chat-service/internal/usecases/manager/free-hands"
	freehandsmocks "github.com/keepcalmist/chat-service/internal/usecases/manager/free-hands/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	loadService *freehandsmocks.MockmanagerLoadService
	managerPool *freehandsmocks.MockmanagerPool
	uCase       freehands.UseCase
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.loadService = freehandsmocks.NewMockmanagerLoadService(s.ctrl)
	s.managerPool = freehandsmocks.NewMockmanagerPool(s.ctrl)

	uCase, err := freehands.New(freehands.NewOptions(s.loadService, s.managerPool))
	s.Require().NoError(err)

	s.uCase = uCase

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) TestUseCaseHandle_RequestValidation() {
	s.loadService.EXPECT().CanManagerTakeProblem(s.Ctx, gomock.Any()).Return(true, nil)
	s.managerPool.EXPECT().Put(s.Ctx, gomock.Any()).Return(nil)

	tCase := []struct {
		name    string
		request freehands.Request
		wantErr bool
	}{
		{
			name: "valid",
			request: freehands.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},
		{
			name:    "invalid",
			request: freehands.Request{},
			wantErr: true,
		},
	}

	for _, tc := range tCase {
		s.Run(tc.name, func() {
			err := s.uCase.Handle(s.Ctx, tc.request)
			if tc.wantErr {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)
		})
	}
}

func (s *UseCaseSuite) TestUseCaseHandle() {
	managerID := types.NewUserID()

	s.Run("error from CanManagerTakeProblem", func() {
		unknownError := errors.New("unknown error")

		s.loadService.EXPECT().CanManagerTakeProblem(s.Ctx, managerID).Return(false, unknownError)

		err := s.uCase.Handle(s.Ctx, freehands.Request{
			ID:        types.NewRequestID(),
			ManagerID: managerID,
		})
		s.Require().ErrorIs(err, unknownError)
	})

	s.Run("CanManagerTakeProblem == false", func() {
		s.loadService.EXPECT().CanManagerTakeProblem(s.Ctx, managerID).Return(false, nil)

		err := s.uCase.Handle(s.Ctx, freehands.Request{
			ID:        types.NewRequestID(),
			ManagerID: managerID,
		})
		s.Require().ErrorIs(err, freehands.ErrManagerCannotTakeMoreProblems)
	})

	s.Run("successful case", func() {
		s.loadService.EXPECT().CanManagerTakeProblem(s.Ctx, managerID).Return(true, nil)
		s.managerPool.EXPECT().Put(s.Ctx, managerID).Return(nil)

		err := s.uCase.Handle(s.Ctx, freehands.Request{
			ID:        types.NewRequestID(),
			ManagerID: managerID,
		})
		s.Require().NoError(err)
	})
}
