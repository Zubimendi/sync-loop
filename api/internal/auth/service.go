package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Zubimendi/sync-loop/api/internal/model"
	"github.com/Zubimendi/sync-loop/api/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users *repo.UserRepo
}

func NewService(users *repo.UserRepo) *Service {
	return &Service{users: users}
}

func (s *Service) Register(ctx context.Context, email, password, workspaceName string) (*model.User, *model.Workspace, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, "", err
	}
	u, w, err := s.users.CreateUserAndWorkspace(ctx, email, string(hash), workspaceName)
	if err != nil {
		return nil, nil, "", err
	}
	token, err := s.makeToken(u.ID, w.ID)
	if err != nil {
		return nil, nil, "", err
	}
	return u, w, token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*model.User, *model.Workspace, string, error) {
u, err := s.users.GetByEmail(ctx, email)
if err != nil || u == nil {
return nil, nil, "", errors.New("invalid credentials")
}
if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
return nil, nil, "", errors.New("invalid credentials")
}
// fetch real workspace
w, err := s.users.GetWorkspaceByOwner(ctx, u.ID)
if err != nil || w == nil {
return nil, nil, "", errors.New("workspace not found")
}
token, err := s.makeToken(u.ID, w.ID)
if err != nil {
return nil, nil, "", err
}
return u, w, token, nil
}

func (s *Service) makeToken(userID, workspaceID string) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"wid": workspaceID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func(s * Service) ValidateToken(tokenString string)(userID, workspaceID string, err error) {
    token, err := jwt.Parse(tokenString, func(token * jwt.Token)(interface {}, error) {
        if _, ok := token.Method.( * jwt.SigningMethodHMAC);
        !ok {
            return nil, errors.New("unexpected signing method")
        }
        return [] byte(os.Getenv("JWT_SECRET")), nil
    })
    if err != nil || !token.Valid {
        return "", "", errors.New("invalid token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", "", errors.New("invalid claims")
    }
    uid, _:= claims["uid"].(string)
    wid, _:= claims["wid"].(string)
    if uid == "" || wid == "" {
        return "", "", errors.New("missing ids")
    }
    return uid, wid, nil
}