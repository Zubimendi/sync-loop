package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type CopyTableParams struct {
	Table        string
	ConnectorID  string
	Incremental  bool
	LastSyncTime time.Time
}

func CopyTableWorkflow(ctx workflow.Context, params CopyTableParams) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting CopyTableWorkflow", 
		"table", params.Table, 
		"incremental", params.Incremental)

	// Set workflow options for retries
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
	})

	// Query handler for current state
	var currentState = "starting"
	workflow.SetQueryHandler(ctx, "state", func() (string, error) {
		return currentState, nil
	})

	// Get last sync time if incremental
	var lastSyncTime = params.LastSyncTime
	if params.Incremental && lastSyncTime.IsZero() {
		currentState = "fetching_last_sync_time"
		var lastSync LastSyncInfo
		err := workflow.ExecuteActivity(ctx, "GetLastSyncTimeActivity", GetLastSyncTimeParams{
			ConnectorID: params.ConnectorID,
			Table:       params.Table,
		}).Get(ctx, &lastSync)
		
		if err != nil {
			logger.Error("Failed to get last sync time", "error", err)
			// Continue with full sync if we can't get last sync time
			params.Incremental = false
		} else if !lastSync.LastSyncTime.IsZero() {
			lastSyncTime = lastSync.LastSyncTime
		}
	}

	// Extract data
	currentState = "extracting"
	var extractResult ExtractResult
	err := workflow.ExecuteActivity(ctx, "ExtractActivity", ExtractParams{
		Table:        params.Table,
		ConnectorID:  params.ConnectorID,
		Incremental:  params.Incremental,
		LastSyncTime: lastSyncTime,
	}).Get(ctx, &extractResult)
	
	if err != nil {
		currentState = "extract_failed"
		logger.Error("ExtractActivity failed", "error", err)
		return fmt.Errorf("extract failed: %w", err)
	}

	if extractResult.RowCount == 0 {
		currentState = "no_data_to_process"
		logger.Info("No new data to process")
		return nil
	}

	// Transform data
	currentState = "transforming"
	var transformResult TransformResult
	err = workflow.ExecuteActivity(ctx, "TransformActivity", TransformParams{
		Data:      extractResult.Data,
		Table:     params.Table,
	}).Get(ctx, &transformResult)
	
	if err != nil {
		currentState = "transform_failed"
		logger.Error("TransformActivity failed", "error", err)
		return fmt.Errorf("transform failed: %w", err)
	}

	// Load data
	currentState = "loading"
	var loadResult LoadResult
	err = workflow.ExecuteActivity(ctx, "LoadActivity", LoadParams{
		Data:        transformResult.Data,
		Table:       params.Table,
		ConnectorID: params.ConnectorID,
	}).Get(ctx, &loadResult)
	
	if err != nil {
		currentState = "load_failed"
		logger.Error("LoadActivity failed", "error", err)
		return fmt.Errorf("load failed: %w", err)
	}

	// Update last sync time if incremental
	if params.Incremental {
		currentState = "updating_sync_time"
		err = workflow.ExecuteActivity(ctx, "UpdateLastSyncTimeActivity", UpdateLastSyncTimeParams{
			ConnectorID: params.ConnectorID,
			Table:       params.Table,
			SyncTime:    extractResult.MaxTimestamp,
		}).Get(ctx, nil)
		
		if err != nil {
			currentState = "sync_time_update_failed"
			logger.Error("UpdateLastSyncTimeActivity failed", "error", err)
			// Don't fail workflow for this - just log
		}
	}

	currentState = "completed"
	logger.Info("CopyTableWorkflow completed successfully", 
		"rows_processed", loadResult.RowsProcessed,
		"incremental", params.Incremental)
	
	return nil
}

// Activity parameter types
type ExtractParams struct {
	Table        string
	ConnectorID  string
	Incremental  bool
	LastSyncTime time.Time
}

type ExtractResult struct {
	Data         []map[string]interface{}
	RowCount     int64
	MaxTimestamp time.Time
	Checksum     string
}

type TransformParams struct {
	Data  []map[string]interface{}
	Table string
}

type TransformResult struct {
	Data     []map[string]interface{}
	RowCount int64
}

type LoadParams struct {
	Data        []map[string]interface{}
	Table       string
	ConnectorID string
}

type LoadResult struct {
	RowsProcessed int64
	Success       bool
}

type GetLastSyncTimeParams struct {
	ConnectorID string
	Table       string
}

type LastSyncInfo struct {
	LastSyncTime time.Time
}

type UpdateLastSyncTimeParams struct {
	ConnectorID string
	Table       string
	SyncTime    time.Time
}