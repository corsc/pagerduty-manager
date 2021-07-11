package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/corsc/pagerduty-manager/internal/pd"

	"go.uber.org/zap"
)

const (
	getURI  = "/services/%s"
	listURI = "/services"
	addURI  = "/services"
)

var ErrNoSuchService = errors.New("no such service")

func New(cfg Config, logger *zap.Logger) *Manager {
	return &Manager{
		cfg:    cfg,
		logger: logger,
		api:    pd.New(cfg, logger),
	}
}

// Manager allows for loading and creating services
type Manager struct {
	cfg    Config
	logger *zap.Logger
	api    *pd.API
}

func (u *Manager) Get(ctx context.Context, serviceID string) (*Service, error) {
	uri := fmt.Sprintf(getURI, serviceID)

	services := &getServiceResponse{}

	err := u.api.Get(ctx, uri, nil, services)
	if err != nil {
		return nil, fmt.Errorf("failed to get service '%s' with err: %s", serviceID, err)
	}

	if services.Service == nil {
		return nil, ErrNoSuchService
	}

	return services.Service, nil
}

func (u *Manager) GetByName(ctx context.Context, name string) (*Service, error) {
	params := url.Values{}
	params.Set("query", name)
	params.Set("total", "false")
	params.Set("limit", "1")

	services := &getServicesResponse{}

	err := u.api.Get(ctx, listURI, nil, services)
	if err != nil {
		return nil, fmt.Errorf("failed to get services '%s' with err: %s", name, err)
	}

	if len(services.Service) == 0 {
		return nil, ErrNoSuchService
	}

	return services.Service[0], nil
}

func (u *Manager) Add(ctx context.Context, service NewService) (string, error) {
	reqDTO := &addRequest{
		Service: &Service{
			Name:        service.GetName(),
			Description: service.GetDescription(),
			Status:      "active",
			EscalationPolicy: &escalationPolicy{
				ID: service.GetEscalationPolicyID(),
			},
			Teams: []*team{
				{
					ID: service.GetTeamID(),
				},
			},
			IncidentUrgencyRule: &incidentUrgency{
				Type: "constant",
			},
			AlertCreation: "create_alerts_and_incidents",
			AlertGroupingParameters: &alertGroupParameters{
				Type: "intelligent",
			},
		},
	}

	respDTO := &addResponse{}

	err := u.api.Post(ctx, addURI, reqDTO, respDTO)
	if err != nil {
		return "", fmt.Errorf("failed to add service '%#v' with err: %s", reqDTO, err)
	}

	return respDTO.Service.ID, nil
}

type NewService interface {
	GetName() string
	GetDescription() string
	GetEscalationPolicyID() string
	GetTeamID() string
}

type getServiceResponse struct {
	Service *Service `json:"service"`
}

type getServicesResponse struct {
	Service []*Service `json:"services"`
}

type Service struct {
	ID                      string                `json:"id"`
	Type                    string                `json:"type"`
	Name                    string                `json:"name"`
	Description             string                `json:"description"`
	Status                  string                `json:"status"`
	EscalationPolicy        *escalationPolicy     `json:"escalation_policy"`
	Teams                   []*team               `json:"teams"`
	IncidentUrgencyRule     *incidentUrgency      `json:"incident_urgency_rule"`
	AlertCreation           string                `json:"alert_creation"`
	AlertGroupingParameters *alertGroupParameters `json:"alert_grouping_parameters"`
}

type escalationPolicy struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type team struct {
	ID string `json:"id"`
}

type incidentUrgency struct {
	Type string `json:"type"`
}

type alertGroupParameters struct {
	Type string `json:"type"`
}

type addRequest struct {
	Service *Service `json:"service"`
}

type addResponse struct {
	Service *Service `json:"service"`
}

type Config interface {
	Debug() bool
	BaseURL() string
	AuthToken() string
}
