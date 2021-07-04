package teams

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/corsc/pagerduty-manager/internal/pd"

	"go.uber.org/zap"
)

const (
	getURI         = "/teams/%s"
	listMembersURI = "/teams/%s/members"
	addMemberURI   = "/teams/%s/users/%s"
)

var (
	ErrNoSuchTeam = errors.New("no such user")
	ErrNoMembers  = errors.New("no members")
)

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

func (u *Manager) Get(ctx context.Context, teamID string) (*Team, error) {
	uri := fmt.Sprintf(getURI, teamID)

	teams := &getTeamResponse{}

	err := u.api.Get(ctx, uri, nil, teams)
	if err != nil {
		return nil, fmt.Errorf("failed to get team '%s' with err: %s", teamID, err)
	}

	if teams.Team == nil {
		return nil, ErrNoSuchTeam
	}

	return teams.Team, nil
}

func (u *Manager) GetMembers(ctx context.Context, teamID string) ([]*Member, error) {
	uri := fmt.Sprintf(listMembersURI, teamID)

	params := url.Values{}
	params.Set("total", "true")

	team := &getTeamMembersResponse{}

	err := u.api.Get(ctx, uri, params, team)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members for team '%s' with err: %s", teamID, err)
	}

	if team.Members == nil {
		return nil, ErrNoMembers
	}

	members := make([]*Member, len(team.Members))

	for index, member := range team.Members {
		members[index] = &Member{
			ID:   member.User.ID,
			Role: member.Role,
		}
	}

	return members, nil
}

func (u *Manager) AddMember(ctx context.Context, teamID string, user User) error {
	uri := fmt.Sprintf(addMemberURI, teamID, user.GetUserID())

	payload := &addMemberRequest{
		Role: user.GetRole(),
	}

	err := u.api.Put(ctx, uri, payload)
	if err != nil {
		return fmt.Errorf("failed to add user '%#v' to team '%s' with err: %s", user, teamID, err)
	}

	return nil
}

type User interface {
	GetUserID() string
	GetRole() string
}

type getTeamResponse struct {
	Team *Team `json:"team"`
}

type Team struct {
	Name string `json:"name"`
}

type Member struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type getTeamMembersResponse struct {
	Members []*member `json:"members"`
}

type member struct {
	User *user  `json:"user"`
	Role string `json:"role"`
}

type user struct {
	ID string `json:"id"`
}

type addMemberRequest struct {
	Role string `json:"role"`
}

type Config interface {
	Debug() bool
	BaseURL() string
	AuthToken() string
}
