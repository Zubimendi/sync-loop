package main

import (
	"context"
	"log"

	"go.temporal.io/api/common/v1"        
	"go.temporal.io/api/workflowservice/v1"
	"github.com/Zubimendi/sync-loop/api/internal/temporal"
)

func main() {
	if err := temporal.Init(); err != nil {
		log.Fatal(err)
	}
	defer temporal.Close()

	// Handle both failed and cancelled workflows
	resp, err := temporal.DefaultClient.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
		Query: `WorkflowType="CopyTableWorkflow" AND (ExecutionStatus="Failed" OR ExecutionStatus="Canceled" OR ExecutionStatus="Terminated")`,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range resp.Executions {
		_, err := temporal.DefaultClient.WorkflowService().DeleteWorkflowExecution(context.Background(), &workflowservice.DeleteWorkflowExecutionRequest{
			Namespace: "default",
			WorkflowExecution: &common.WorkflowExecution{
				WorkflowId: info.Execution.WorkflowId,
				RunId:      info.Execution.RunId,
			},
		})
		if err != nil {
			log.Printf("failed to delete %s: %v", info.Execution.WorkflowId, err)
			continue
		}
		log.Printf("deleted %s (status: %s)", info.Execution.WorkflowId, info.Status.String())
	}
	log.Println("purge completed")
}