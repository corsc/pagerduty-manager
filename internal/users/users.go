package users

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/corsc/pagerduty-manager/internal/pd"

	"go.uber.org/zap"
)

const uri = "/users"

var ErrNoSuchUser = errors.New("no such user")

func New(cfg Config, logger *zap.Logger) *Manager {
	return &Manager{
		cfg:    cfg,
		logger: logger,
		api:    pd.New(cfg, logger),
	}
}

// Manager allows for loading and creating users
type Manager struct {
	cfg    Config
	logger *zap.Logger
	api    *pd.API
}

func (u *Manager) Get(ctx context.Context, email string) (*User, error) {
	params := url.Values{}
	params.Set("query", url.QueryEscape(email))
	params.Set("total", "false")
	params.Set("limit", "1")

	users := &getResponse{}

	err := u.api.Get(ctx, uri, params, users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user '%s` with err: %s", email, err)
	}

	if len(users.Users) == 0 {
		return nil, ErrNoSuchUser
	}

	return users.Users[0], nil
}

func (u *Manager) Add(ctx context.Context, user NewUser, defaultTimeZone string) (string, error) {
	reqDTO := newNewUserRequest(user, defaultTimeZone)

	respDTO := &newUserResponse{}

	err := u.api.Post(ctx, uri, reqDTO, respDTO)
	if err != nil {
		return "", fmt.Errorf("failed to add user '%#v` with err: %s", user, err)
	}

	return respDTO.User.ID, nil
}

type NewUser interface {
	GetName() string
	GetEmail() string
	GetTimeZone() string
	GetRole() string
}

func newNewUserRequest(user NewUser, defaultTimeZone string) *newUserRequest {
	out := &newUserRequest{
		User: userFormat{
			Type:     "user",
			Name:     user.GetName(),
			Email:    user.GetEmail(),
			TimeZone: user.GetTimeZone(),
			Role:     user.GetRole(),
		},
	}

	if out.User.TimeZone == "" {
		out.User.TimeZone = defaultTimeZone
	}

	return out
}

type newUserRequest struct {
	User userFormat `json:"user"`
}

type newUserResponse struct {
	User userFormat `json:"user"`
}

type userFormat struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	TimeZone string `json:"time_zone"`
	Role     string `json:"role"`
}

type getResponse struct {
	Users []*User `json:"users"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Teams []Team `json:"teams"`
}

type Team struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type Config interface {
	Debug() bool
	BaseURL() string
	AuthToken() string
}
