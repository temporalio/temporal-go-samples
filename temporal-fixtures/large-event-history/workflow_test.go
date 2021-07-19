package largeeventhistory

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/worker"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_LargeEventHistoryWorkflow() {
	env := s.NewTestWorkflowEnvironment()
	env.SetWorkerOptions(worker.Options{})

	env.OnActivity(Activity, mock.Anything).Return(nil)

	env.RegisterActivity(Activity)

	env.ExecuteWorkflow(LargeEventHistoryWorkflow, 1024, false)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}
