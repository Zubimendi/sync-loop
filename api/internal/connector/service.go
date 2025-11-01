package connector

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Zubimendi/sync-loop/api/internal/encrypt"
	"github.com/Zubimendi/sync-loop/api/internal/model"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service { return &Service{repo: repo} }

var hardCodedTypes = map[string]bool{
	"pg": true, "mysql": true, "s3": true, "excel": true,
	"gsheets": true, "sf": true, "rest": true,
}

func (s *Service) CreateSource(ctx context.Context, name, ctype string, config map[string]interface{}, userID, workspaceID string) (*model.Connector, error) {
	if !hardCodedTypes[ctype] {
		return nil, errors.New("unsupported connector type")
	}
	plain, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}
	cipher, err := encrypt.Encrypt(string(plain))
	if err != nil {
		return nil, fmt.Errorf("encrypt config: %w", err)
	}
	c := &model.Connector{
		ID:        hex.EncodeToString(randBytes(16)),
		Name:      name,
		Type:      ctype,
		Config:    cipher,
		CreatedBy: userID,
		Workspace: workspaceID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create connector: %w", err)
	}
	return c, nil
}
func (s *Service) ListSources(ctx context.Context, workspaceID string) ([]model.Connector, error) {
	return s.repo.ListByWorkspace(ctx, workspaceID)
}

func randBytes(n int) []byte { b := make([]byte, n); rand.Read(b); return b }