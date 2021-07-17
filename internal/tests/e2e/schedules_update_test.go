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

func TestE2ESchedules_Update(t *testing.T) {
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
		memberIDs: []string{"PXJHUO9", "P7X168R"},
	}

	scheduleID := "PJQM2NF"

	// call object under test
	manager := schedules.New(cfg, logger)
	resultErr := manager.Update(ctx, scheduleID, schedule)

	// validation
	require.NoError(t, resultErr)
}
