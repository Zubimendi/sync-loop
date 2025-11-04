package scheduler

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"github.com/google/uuid"
)

type Service struct {
	temporalClient client.Client
}

type ScheduleConfig struct {
	ID          string
	ConnectorID string
	CronExpr    string
	Table       string
	IsActive    bool
}

func NewService(temporalClient client.Client) *Service {
	return &Service{temporalClient: temporalClient}
}

// CreateSchedule creates a new Temporal schedule for automatic syncs
func (s *Service) CreateSchedule(ctx context.Context, config ScheduleConfig) error {
	scheduleID := fmt.Sprintf("schedule-%s-%s", config.ConnectorID, config.ID)
	
	// Create schedule with cron expression (every minute by default)
	_, err := s.temporalClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{config.CronExpr}, // "* * * * *" for every minute
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        fmt.Sprintf("sync-%s-%d", config.ConnectorID, time.Now().Unix()),
			Workflow:  "CopyTableWorkflow",
			TaskQueue: "sync-loop-task-queue",
			Args:      []interface{}{config.Table, config.ConnectorID, true}, // true = incremental
		},
		State: &client.ScheduleState{
			Paused: !config.IsActive,
			Notes:  fmt.Sprintf("Auto-sync for connector %s", config.ConnectorID),
		},
	})
	
	return err
}

// UpdateSchedule updates an existing schedule
func (s *Service) UpdateSchedule(ctx context.Context, scheduleID string, cronExpr string, isActive bool) error {
	handle := s.temporalClient.ScheduleClient().GetHandle(ctx, scheduleID)
	
	_, err := handle.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(schedule client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			return &client.ScheduleUpdate{
				Schedule: &schedule.Schedule,
				Spec: &client.ScheduleSpec{
					CronExpressions: []string{cronExpr},
				},
				State: &client.ScheduleState{
					Paused: !isActive,
					Notes:  "Updated schedule",
				},
			}, nil
		},
	})
	
	return err
}

// PauseSchedule pauses a schedule
func (s *Service) PauseSchedule(ctx context.Context, scheduleID string, reason string) error {
	handle := s.temporalClient.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Pause(ctx, client.SchedulePauseOptions{Note: reason})
}

// UnpauseSchedule unpauses a schedule
func (s *Service) UnpauseSchedule(ctx context.Context, scheduleID string, reason string) error {
	handle := s.temporalClient.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Unpause(ctx, client.ScheduleUnpauseOptions{Note: reason})
}

// ListSchedules lists all schedules for a workspace
func (s *Service) ListSchedules(ctx context.Context, workspaceID string) ([]ScheduleInfo, error) {
	var schedules []ScheduleInfo
	
	iter := s.temporalClient.ScheduleClient().List(ctx, client.ScheduleListOptions{
		PageSize: 50,
	})
	
	for iter.HasNext() {
		entry, err := iter.Next()
		if err != nil {
			break
		}
		
		// Filter by workspace ID (extract from schedule ID)
		if entry.ID != "" && len(entry.ID) > len("schedule-") {
			schedules = append(schedules, ScheduleInfo{
				ID:          entry.ID,
				WorkflowID:  entry.Schedule.Action.GetWorkflow().ID,
				CronExpr:    entry.Schedule.Spec.CronExpressions[0],
				IsActive:    !entry.Schedule.State.Paused,
				LastRunTime: entry.Schedule.State.LastProcessedTime,
				NextRunTime: entry.Schedule.State.NextActionTime,
			})
		}
	}
	
	return schedules, nil
}

type ScheduleInfo struct {
	ID          string    `json:"id"`
	WorkflowID  string    `json:"workflow_id"`
	CronExpr    string    `json:"cron_expr"`
	IsActive    bool      `json:"is_active"`
	LastRunTime time.Time `json:"last_run_time"`
	NextRunTime time.Time `json:"next_run_time"`
}