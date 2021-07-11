package pdmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"
)

// map of our roles to PD roles
var validRoles = map[string]string{
	"member":    "user",
	"observer":  "observer",
	"lead":      "admin",
	"dept-head": "admin",
}

func New(cfg Config, logger *zap.Logger) *Manager {
	return &Manager{
		cfg:           cfg,
		logger:        logger,
		companyConfig: &companyConfig{},
	}
}

// Manager is the main entry point for this package/tool
type Manager struct {
	cfg    Config
	logger *zap.Logger

	companyConfig *companyConfig
}

// Parse attempts to parse the provide file into this manager
func (m *Manager) Parse(_ context.Context) error {
	m.logger.Debug("loading data from file", zap.String("file", m.cfg.Filename()))

	fileContents, err := ioutil.ReadFile(m.cfg.Filename())
	if err != nil {
		return fmt.Errorf("failed to read input file with err: %w", err)
	}

	err = json.Unmarshal(fileContents, m.companyConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config JSON with err: %w", err)
	}

	return m.validate()
}

func (m *Manager) validate() error {
	if len(m.companyConfig.Teams) == 0 {
		return errors.New("no teams found in the JSON")
	}

	for _, thisTeam := range m.companyConfig.Teams {
		for _, thisMember := range thisTeam.Members {
			_, ok := validRoles[thisMember.Role]
			if !ok {
				return fmt.Errorf("invalid role 'value in: %v", thisMember)
			}
		}
	}

	return nil
}

// SyncTeams attempts to download the existing teams and create any that do not yet exist.
// Note: existing data will not be modified in any way.
// Note: creating a team also creates a matching service so we can have an `@oncall-[team]` slack alias
func (m *Manager) SyncTeams(ctx context.Context) error {
	return errors.New("not implemented")
}

// SyncUsers attempts to download the existing users and create any that do not yet exist.
// Note: existing data will not be modified in any way.
func (m *Manager) SyncUsers(ctx context.Context) error {
	return errors.New("not implemented")
}

// SyncServices attempts to download the existing services and create any that do not yet exist.
// Note: existing data will not be modified in any way.
func (m *Manager) SyncServices(ctx context.Context) error {
	return errors.New("not implemented")
}

// SyncEscalation attempts to download the existing escalation policies and create any that do not yet exist.
// Note: existing data will not be modified in any way.
func (m *Manager) SyncEscalation(ctx context.Context) error {
	return errors.New("not implemented")
}

// SyncSchedules attempts to download the existing schedules and create any that do not yet exist.
// Note: existing data will not be modified in any way.
func (m *Manager) SyncSchedules(ctx context.Context) error {
	return errors.New("not implemented")
}

// Config is the config for this package
type Config interface {
	Debug() bool
	Filename() string
}

type companyConfig struct {
	Teams           []*Team `json:"teams"`
	DefaultTimezone string  `json:"default_timezone"`
}

type Team struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Slack       string      `json:"slack"`
	Escalation  *Escalation `json:"escalation"`
	Members     []*Member   `json:"members"`
	Services    []*Service  `json:"services"`
}

type Member struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Timezone string `json:"timezone"`
	Role     string `json:"role"`
}

type Service struct {
	Name      string `json:"name"`
	Dashboard string `json:"dashboard"`
}

type Escalation struct {
	Handover struct {
		Day  string `json:"day"`
		Time string `json:"time"`
	} `json:"handover"`
	After string `json:"after"`
}
