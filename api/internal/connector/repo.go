package connector

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/Zubimendi/sync-loop/api/internal/model"
	"github.com/rs/zerolog/log"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo { return &Repo{db: db} }

func (r *Repo) Create(ctx context.Context, c *model.Connector) error {
	q := `INSERT INTO connector (id,name,type,config_json,created_by_user_id,workspace_id)
	      VALUES ($1,$2,$3,to_jsonb($4::text),$5,$6)`
	log.Printf("SQL: %s | args: %#v", q, []interface{}{c.ID, c.Name, c.Type, c.Config, c.CreatedBy, c.Workspace})
	_, err := r.db.ExecContext(ctx, q, c.ID, c.Name, c.Type, c.Config, c.CreatedBy, c.Workspace)
	return err
}

func (r *Repo) ListByWorkspace(ctx context.Context, workspaceID string) ([]model.Connector, error) {
	cc := make([]model.Connector, 0) // non-nil empty slice
	err := r.db.SelectContext(ctx, &cc,
		`SELECT id,name,type,created_at,updated_at FROM connector WHERE workspace_id=$1`,
		workspaceID)
	return cc, err
}