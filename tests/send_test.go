package email_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	email "github.com/gofor-little/aws-email"
)

func TestSend(t *testing.T) {
	recipients, from := setup(t)

	testCases := []struct {
		name string
		data email.Data
		want error
	}{
		{"TestSend_Integration", email.Data{To: recipients, From: from, Subject: "Test Subject", Body: "Test Body", ContentType: email.ContentTypeTextPlain}, nil},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%s_%d", tc.name, i)

		t.Run(name, func(t *testing.T) {
			_, err := email.Send(context.Background(), tc.data)
			require.Equal(t, tc.want, err)
		})
	}
}
