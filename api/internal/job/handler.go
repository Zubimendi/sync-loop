package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Zubimendi/sync-loop/api/internal/middleware"
	"github.com/Zubimendi/sync-loop/api/internal/workflow"
	"github.com/rs/zerolog/log"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	temporal client.Client
}

func NewHandler(temporal client.Client) *Handler { 
	return &Handler{temporal: temporal} 
}

type JobResp struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"start_time"`
	Table     string    `json:"table"`
	WorkflowType string `json:"workflow_type"` 
	Schedule  *ScheduleInfo `json:"schedule,omitempty"`
}

type ScheduleInfo struct {
	ID          string    `json:"id"`
	CronExpr    string    `json:"cron_expr"`
	IsActive    bool      `json:"is_active"`
	LastRunTime time.Time `json:"last_run_time"`
	NextRunTime time.Time `json:"next_run_time"`
}

// Update the List function to capture workflow type
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.CtxWorkspaceID).(string)
	
	resp, err := h.temporal.ListWorkflow(r.Context(), &workflowservice.ListWorkflowExecutionsRequest{
		Query: `WorkflowType="CopyTableWorkflow"`,
		PageSize: 20,
	})
	if err != nil {
		log.Error().Err(err).Msg("list workflows")
		json.NewEncoder(w).Encode(map[string]interface{}{"jobs": []JobResp{}})
		return
	}
	
	var out []JobResp
	for _, info := range resp.Executions {
		if info.Execution == nil || info.StartTime == nil {
			continue
		}
		
		// Extract workflow type and table name
		workflowType := info.Type.Name
		if workflowType == "" {
			workflowType = "CopyTableWorkflow" // fallback
		}
		
		out = append(out, JobResp{
			ID:        info.Execution.WorkflowId,
			Type:      info.Type.Name,
			WorkflowType: workflowType,
			Status:    info.Status.String(),
			StartTime: info.StartTime.AsTime(),
			Table:     extractTableFromWorkflowID(info.Execution.WorkflowId),
		})
	}
	
	schedules, err := h.listSchedules(r.Context())
	if err == nil && len(schedules) > 0 {
		for i := range out {
			if i < len(schedules) {
				out[i].Schedule = schedules[i]
			}
		}
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{"jobs": out})
}

// POST /api/v1/jobs/run-now – start workflow immediately (with incremental support)
func (h *Handler) RunNow(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Table       string `json:"table"`
		Incremental bool   `json:"incremental"`
		WorkflowType string `json:"workflow_type"` // Optional: specify workflow type
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	
	// Default to CopyTableWorkflow if not specified
	workflowType := req.WorkflowType
	if workflowType == "" {
		workflowType = "CopyTableWorkflow"
	}
	
	workflowID := fmt.Sprintf("run-now-%s-%s-%d", workflowType, req.Table, time.Now().Unix())
	
	// Use appropriate workflow based on type
	var workflowFunc interface{}
	var workflowArgs interface{}
	
	switch workflowType {
	case "CopyTableWorkflow":
		workflowFunc = workflow.CopyTableWorkflow
		workflowArgs = workflow.CopyTableParams{
			Table:       req.Table,
			Incremental: req.Incremental,
		}
	// Add more workflow types here as you create them
	default:
		http.Error(w, fmt.Sprintf("unknown workflow type: %s", workflowType), http.StatusBadRequest)
		return
	}
	
	we, err := h.temporal.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		TaskQueue: "sync-loop-task-queue",
		ID:        workflowID,
	}, workflowFunc, workflowArgs)
	
	if err != nil {
		log.Error().Err(err).Msg("execute workflow")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
		"workflow_type": workflowType,
	})
}

// POST /api/v1/jobs/cancel
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	var req struct{ WorkflowID string `json:"workflow_id"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	
	err := h.temporal.CancelWorkflow(r.Context(), req.WorkflowID, "")
	if err != nil {
		log.Error().Err(err).Msg("cancel workflow")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/v1/jobs/terminate-all
func (h *Handler) TerminateAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.temporal.ListWorkflow(r.Context(), &workflowservice.ListWorkflowExecutionsRequest{
		Query: `WorkflowType="CopyTableWorkflow" AND ExecutionStatus="Running"`,
	})
	if err != nil {
		log.Error().Err(err).Msg("list running workflows")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	for _, info := range resp.Executions {
		if info.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			_ = h.temporal.CancelWorkflow(r.Context(), info.Execution.WorkflowId, "")
		}
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// NEW: POST /api/v1/jobs/retry – retry a failed workflow
func (h *Handler) Retry(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	
	// Get original workflow info
	_, err := h.temporal.DescribeWorkflowExecution(r.Context(), req.WorkflowID, req.RunID)
	if err != nil {
		log.Error().Err(err).Msg("describe workflow")
		http.Error(w, "workflow not found", http.StatusNotFound)
		return
	}
	
	// Extract table name from workflow ID
	table := extractTableFromWorkflowID(req.WorkflowID)
	
	// Start new workflow with same parameters but new ID
	newWorkflowID := fmt.Sprintf("retry-%s-%d", req.WorkflowID, time.Now().Unix())
	
	we, err := h.temporal.ExecuteWorkflow(r.Context(), client.StartWorkflowOptions{
		TaskQueue: "sync-loop-task-queue",
		ID:        newWorkflowID,
	}, workflow.CopyTableWorkflow, workflow.CopyTableParams{
		Table:       table,
		Incremental: false, // Retry as full sync for safety
	})
	
	if err != nil {
		log.Error().Err(err).Msg("execute retry workflow")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
		"message":     "Retry started",
	})
}

// NEW: POST /api/v1/jobs/schedule – create/update schedule
func (h *Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConnectorID string `json:"connector_id"`
		Table       string `json:"table"`
		CronExpr    string `json:"cron_expr"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	
	// Default to every minute if no cron expression
	if req.CronExpr == "" {
		req.CronExpr = "* * * * *"
	}
	
	// Make schedule ID unique by adding timestamp
	scheduleID := fmt.Sprintf("schedule-%s-%s-%d", req.ConnectorID, req.Table, time.Now().Unix())
	
	// Check if schedule already exists
	existingHandle := h.temporal.ScheduleClient().GetHandle(r.Context(), scheduleID)
	_, err := existingHandle.Describe(r.Context())
	if err == nil {
		// Schedule already exists, return error
		http.Error(w, "schedule with this connector and table already exists", http.StatusConflict)
		return
	}
	
	// Create the schedule
	_, err = h.temporal.ScheduleClient().Create(r.Context(), client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{req.CronExpr},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        fmt.Sprintf("scheduled-%s-%d", req.Table, time.Now().Unix()),
			Workflow:  workflow.CopyTableWorkflow,
			TaskQueue: "sync-loop-task-queue",
			Args: []interface{}{workflow.CopyTableParams{
				Table:       req.Table,
				Incremental: true, // Schedules default to incremental
			}},
		},
		Paused: !req.IsActive,
	})
	
	if err != nil {
		log.Error().Err(err).Msg("create schedule")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"schedule_id": scheduleID,
		"message":     "Schedule created",
	})
}

// NEW: POST /api/v1/jobs/schedule/toggle – pause/unpause schedule
func (h *Handler) ToggleSchedule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ScheduleID string `json:"schedule_id"`
		Pause      bool   `json:"pause"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	
	handle := h.temporal.ScheduleClient().GetHandle(r.Context(), req.ScheduleID)
	
	if req.Pause {
		err := handle.Pause(r.Context(), client.SchedulePauseOptions{
			Note: "User requested pause",
		})
		if err != nil {
			log.Error().Err(err).Msg("pause schedule")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := handle.Unpause(r.Context(), client.ScheduleUnpauseOptions{
			Note: "User requested unpause",
		})
		if err != nil {
			log.Error().Err(err).Msg("unpause schedule")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"schedule_id": req.ScheduleID,
		"paused":      req.Pause,
	})
}

// Helper: list schedules for workspace
func (h *Handler) listSchedules(ctx context.Context) ([]*ScheduleInfo, error) {
	var schedules []*ScheduleInfo
	
	iter, err := h.temporal.ScheduleClient().List(ctx, client.ScheduleListOptions{
		PageSize: 50,
	})
	if err != nil {
		return nil, err
	}
	
	for iter.HasNext() {
		entry, err := iter.Next()
		if err != nil {
			break
		}
		
		handle := h.temporal.ScheduleClient().GetHandle(ctx, entry.ID)
		desc, err := handle.Describe(ctx)
		if err != nil {
			continue
		}
		
		// Safely extract information
		cronExpr := ""
		isActive := false
		
		// Check for cron expressions in the schedule spec
		if len(desc.Schedule.Spec.CronExpressions) > 0 {
			cronExpr = desc.Schedule.Spec.CronExpressions[0]
		}
		
		// Check if schedule is active (not paused)
		isActive = !desc.Schedule.State.Paused
		
		// For now, we'll skip the timing fields since they seem to have different names
		// You can add them back once we figure out the correct field names
		
		if cronExpr != "" {
			schedules = append(schedules, &ScheduleInfo{
				ID:       entry.ID,
				CronExpr: cronExpr,
				IsActive: isActive,
				// Leave timing fields empty for now
				LastRunTime: time.Time{},
				NextRunTime: time.Time{},
			})
		}
	}
	
	return schedules, nil
}

// Helper: extract table name from workflow ID
func extractTableFromWorkflowID(workflowID string) string {
	// Handle different workflow ID formats
	if len(workflowID) > len("run-now-") && workflowID[:len("run-now-")] == "run-now-" {
		return workflowID[len("run-now-"):]
	}
	if len(workflowID) > len("scheduled-") && workflowID[:len("scheduled-")] == "scheduled-" {
		parts := workflowID[len("scheduled-"):]
		return parts
	}
	if len(workflowID) > len("retry-") && workflowID[:len("retry-")] == "retry-" {
		parts := workflowID[len("retry-"):]
		if len(parts) > 0 {
			return parts
		}
	}
	return "unknown"
}

// GET /api/v1/jobs/:id/status - get detailed job status
func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	workflowID := chi.URLParam(r, "id")
	
	desc, err := h.temporal.DescribeWorkflowExecution(r.Context(), workflowID, "")
	if err != nil {
		http.Error(w, "workflow not found", http.StatusNotFound)
		return
	}
	
	// Get workflow history for error details
	historyIter := h.temporal.GetWorkflowHistory(r.Context(), workflowID, "", false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	if historyIter == nil {
		log.Error().Msg("get workflow history returned nil")
	}
	
	status := map[string]interface{}{
		"workflow_id": workflowID,
		"status":      desc.WorkflowExecutionInfo.Status.String(),
		"start_time":  desc.WorkflowExecutionInfo.StartTime.AsTime(),
	}
	
	if desc.WorkflowExecutionInfo.CloseTime != nil {
		status["close_time"] = desc.WorkflowExecutionInfo.CloseTime.AsTime()
	}
	
	if desc.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_FAILED {
		// Try to get the failure reason from history
		var failureReason string
		for historyIter.HasNext() {
			event, err := historyIter.Next()
			if err != nil {
				break
			}
			if event.GetWorkflowExecutionFailedEventAttributes() != nil {
				failure := event.GetWorkflowExecutionFailedEventAttributes().GetFailure()
				if failure != nil {
					failureReason = failure.GetMessage()
					break
				}
			}
		}
		if failureReason != "" {
			status["failure_reason"] = failureReason
		}
	}
	
	json.NewEncoder(w).Encode(status)
}