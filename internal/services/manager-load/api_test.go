package managerload_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	managerload "github.com/keepcalmist/chat-service/internal/services/manager-load"
	managerloadmocks "github.com/keepcalmist/chat-service/internal/services/manager-load/mocks"
	"github.com/keepcalmist/chat-service/internal/testingh"
	"github.com/keepcalmist/chat-service/internal/types"
)

type ServiceSuite struct {
	testingh.ContextSuite

	ctrl *gomock.Controller

	problemsRepo *managerloadmocks.MockproblemsRepository
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.problemsRepo = managerloadmocks.NewMockproblemsRepository(s.ctrl)
	s.ContextSuite.SetupTest()
}

func (s *ServiceSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *ServiceSuite) TestCanManagerTakeProblem_Validation() {
	tCase := []struct {
		name              string
		maxProblemsAtTime int
		repo              *managerloadmocks.MockproblemsRepository
		isError           bool
	}{
		{
			name:              "maxProblemsAtTime is incorrect #1",
			maxProblemsAtTime: 0,
			repo:              s.problemsRepo,
			isError:           true,
		},
		{
			name:              "maxProblemsAtTime is incorrect #2",
			maxProblemsAtTime: 35,
			repo:              s.problemsRepo,
			isError:           true,
		},
		{
			name:              "repository is nil",
			maxProblemsAtTime: 20,
			repo:              nil,
			isError:           true,
		},
		{
			name:              "correct validation",
			maxProblemsAtTime: 20,
			repo:              s.problemsRepo,
			isError:           false,
		},
	}

	for _, c := range tCase {
		s.Run(c.name, func() {
			managerLoad, err := managerload.New(managerload.NewOptions(c.maxProblemsAtTime, c.repo))
			if c.isError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(managerLoad)
			}
		})
	}
}

func (s *ServiceSuite) TestCanManagerTakeProblem_Successful() {
	managerLoadService, err := managerload.New(managerload.NewOptions(20, s.problemsRepo))
	s.Require().NoError(err)
	managerID := types.NewUserID()

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).
		Return(10, nil)

	ok, err := managerLoadService.CanManagerTakeProblem(s.Ctx, managerID)
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s *ServiceSuite) TestCanManagerTakeProblem_ManagerIsBusy() {
	managerLoadService, err := managerload.New(managerload.NewOptions(5, s.problemsRepo))
	s.Require().NoError(err)
	managerID := types.NewUserID()

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).
		Return(10, nil)

	ok, err := managerLoadService.CanManagerTakeProblem(s.Ctx, managerID)
	s.Require().NoError(err)
	s.Require().False(ok)
}

func (s *ServiceSuite) TestCanManagerTakeProblem_BorderCases() {
	maxProblems := 10
	managerLoadService, err := managerload.New(managerload.NewOptions(maxProblems, s.problemsRepo))
	s.Require().NoError(err)
	managerID := types.NewUserID()

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).
		Return(maxProblems, nil)

	ok, err := managerLoadService.CanManagerTakeProblem(s.Ctx, managerID)
	s.Require().NoError(err)
	s.Require().False(ok)
}

func (s *ServiceSuite) TestCanManagerTakeProblem_ErrorFromRepo() {
	maxProblems := 10
	managerLoadService, err := managerload.New(managerload.NewOptions(maxProblems, s.problemsRepo))
	s.Require().NoError(err)

	managerID := types.NewUserID()
	errFromRepo := errors.New("unknown error")

	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).
		Return(0, errFromRepo)

	ok, err := managerLoadService.CanManagerTakeProblem(s.Ctx, managerID)
	s.Require().ErrorIs(err, errFromRepo)
	s.Require().False(ok)
}
