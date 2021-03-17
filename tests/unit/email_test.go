package email_test

import (
	"context"
	"fmt"
	"testing"

	mock "github.com/gofor-little/aws-sdk-mock"
	"github.com/stretchr/testify/require"

	email "github.com/gofor-little/aws-email"
)

func TestSend(t *testing.T) {
	email.SESClient = &mock.SESClient{}
	email.S3Client = &mock.S3Client{}

	testCases := []struct {
		name string
		data email.Data
		want error
	}{
		{"TestSend_Unit", email.Data{To: []string{"john@example.com"}, From: "jack@example.com", Subject: "Test Subject", Body: "Test Body", ContentType: email.ContentTypeTextPlain}, nil},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%s_%d", tc.name, i)

		t.Run(name, func(t *testing.T) {
			_, err := email.Send(context.Background(), tc.data)
			require.Equal(t, tc.want, err)
		})
	}
}
