package main

import (
	"context"
	"log"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"github.com/Zubimendi/sync-loop/api/internal/temporal"
)

func main() {
	if err := temporal.Init(); err != nil {
		log.Fatal(err)
	}
	defer temporal.Close()

	resp, err := temporal.DefaultClient.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
		Query: `WorkflowType="CopyTableWorkflow" AND ExecutionStatus="Running"`,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range resp.Executions {
		if info.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			if err := temporal.DefaultClient.CancelWorkflow(context.Background(), info.Execution.WorkflowId, ""); err != nil {
				log.Printf("cancel %s: %v", info.Execution.WorkflowId, err)
			} else {
				log.Printf("canceled %s", info.Execution.WorkflowId)
			}
		}}
}		