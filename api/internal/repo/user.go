package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/Zubimendi/sync-loop/api/internal/model"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) CreateUserAndWorkspace(ctx context.Context, email, hash, workspaceName string) (*model.User, *model.Workspace, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	var u model.User
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'owner')
		RETURNING id, email, role, created_at`, email, hash).StructScan(&u)
	if err != nil {
		return nil, nil, err
	}

	var w model.Workspace
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO workspace (name, plan)
		VALUES ($1, 'free')
		RETURNING *`, workspaceName).StructScan(&w)
	if err != nil {
		return nil, nil, err
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO workspace_user (workspace_id, user_id, role)
		VALUES ($1, $2, 'owner')`, w.ID, u.ID); err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}
	return &u, &w, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE email = $1`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func(r * UserRepo) GetWorkspaceByOwner(ctx context.Context, userID string)( * model.Workspace, error) {
    var w model.Workspace
    err := r.db.GetContext(ctx, & w, `
		SELECT w.* 
		FROM workspace w
		JOIN workspace_user wu ON wu.workspace_id = w.id
		WHERE wu.user_id = $1 AND wu.role = 'owner'
		LIMIT 1`, userID)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &w, err
}