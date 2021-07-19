package openNclosed

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

func (s *UnitTestSuite) Test_LargePayloadWorkflow() {
	env := s.NewTestWorkflowEnvironment()
	env.SetWorkerOptions(worker.Options{})

	keepOpen := false
	env.OnActivity(Activity, mock.Anything, keepOpen).Return("hello", nil)

	env.RegisterActivity(Activity)

	env.ExecuteWorkflow(OpenNClosedWorkflow, keepOpen)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}
