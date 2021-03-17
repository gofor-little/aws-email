package email_test

import (
	"strings"
	"testing"

	"github.com/gofor-little/env"
	"github.com/stretchr/testify/require"

	email "github.com/gofor-little/aws-email"
)

func setup(t *testing.T) ([]string, string) {
	if err := env.Load("../../.env"); err != nil {
		t.Logf("failed to load .env file, ignore if running in CI/CD: %v", err)
	}

	require.NoError(t, email.Initialize(env.Get("AWS_PROFILE", ""), env.Get("AWS_REGION", "")))

	from, err := env.MustGet("TEST_EMAIL_FROM")
	require.NoError(t, err)

	recipients, err := env.MustGet("TEST_EMAIL_RECIPIENTS")
	require.NoError(t, err)

	return strings.Split(recipients, ","), from
}
