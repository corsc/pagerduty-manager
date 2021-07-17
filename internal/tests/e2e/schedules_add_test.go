package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/corsc/pagerduty-manager/internal/schedules"

	"github.com/corsc/go-commons/testing/skip"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestE2ESchedules_Add(t *testing.T) {
	skip.IfNotSet(t, "E2E_TEST")

	// inputs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, _ := zap.NewDevelopment()

	cfg := &testConfig{
		baseURL: "https://api.pagerduty.com",
	}

	schedule := &testSchedule{
		teamName:  "Sage42",
		timeZone:  "Australia/Melbourne",
		teamID:    "PJVN6XK",
		memberIDs: []string{"PDPIGEC", "PXJHUO9"},
	}

	// call object under test
	manager := schedules.New(cfg, logger)
	resultID, resultErr := manager.Add(ctx, schedule)

	// validation
	require.NoError(t, resultErr)
	require.NotEmpty(t, resultID)
}

type testSchedule struct {
	teamName    string
	description string
	timeZone    string
	teamID      string
	memberIDs   []string
}

func (t *testSchedule) GetTimeZone() string {
	return t.timeZone
}

func (t *testSchedule) GetMemberIDs() []string {
	return t.memberIDs
}

func (t *testSchedule) GetTeamName() string {
	return t.teamName
}

func (t *testSchedule) GetDescription() string {
	return t.description
}

func (t *testSchedule) GetTeamID() string {
	return t.teamID
}
