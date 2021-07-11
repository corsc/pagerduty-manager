package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/corsc/go-commons/testing/skip"
	"github.com/corsc/pagerduty-manager/internal/users"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestE2EUsers_Add(t *testing.T) {
	skip.IfNotSet(t, "E2E_TEST")

	// inputs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, _ := zap.NewDevelopment()

	cfg := &testConfig{
		baseURL: "https://api.pagerduty.com",
	}

	user := &testUser{
		name:     "Gee",
		email:    "corey.scott@sage42.com",
		timeZone: "Australia/Melbourne",
		role:     "user",
	}
	timeZone := "Australia/Melbourne"

	// call object under test
	manager := users.New(cfg, logger)
	resultUserID, resultErr := manager.Add(ctx, user, timeZone)

	// validation
	require.NoError(t, resultErr)
	require.NotEmpty(t, resultUserID)
}

type testUser struct {
	name     string
	email    string
	timeZone string
	role     string
}

func (t *testUser) GetName() string {
	return t.name
}

func (t *testUser) GetEmail() string {
	return t.email
}

func (t *testUser) GetTimeZone() string {
	return t.timeZone
}

func (t *testUser) GetRole() string {
	return t.role
}