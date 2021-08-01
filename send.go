package email

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/go-mail/mail"
	"github.com/gofor-little/env"
	"github.com/gofor-little/xerror"
)

// Send loads the attachments, builds the email and sends it.
func Send(ctx context.Context, data Data) (string, error) {
	// Check that the package clients have been initialized.
	if SESClient == nil {
		return "", xerror.Newf("SESClient is nil")
	}

	if S3Client == nil {
		return "", xerror.Newf("S3Client is nil")
	}

	// Build the destination emails.
	var destinations []string
	destinations = append(destinations, data.To...)

	if data.CC != nil {
		destinations = append(destinations, *data.CC...)
	}
	if data.BCC != nil {
		destinations = append(destinations, *data.BCC...)
	}

	// Set the headers and body.
	message := mail.NewMessage()
	message.SetHeader("To", destinations...)

	if data.ReplyTo != nil {
		message.SetHeader("Reply-To", *data.ReplyTo...)
	}

	message.SetHeader("From", data.From)
	message.SetHeader("Subject", data.Subject)
	message.SetBody(string(data.ContentType), data.Body)

	// Load and attach the Attachments.
	if data.Attachments != nil {
		for _, a := range data.Attachments {
			data, err := a.Load(ctx)
			if err != nil {
				return "", xerror.Wrap("failed to load Attachment", err)
			}
			message.AttachReader(filepath.Base(a.Path), bytes.NewReader(data))
		}

		// Remove any attachments on disc once we're done with them.
		defer func() {
			// Don't delete attachments when we're developing locally.
			if env.Get("ENVIRONMENT", "development") == "development" {
				return
			}

			for _, a := range data.Attachments {
				_ = os.Remove(a.Path)
			}
		}()
	}

	// Write the email to a buffer.
	var buf bytes.Buffer
	_, err := message.WriteTo(&buf)
	if err != nil {
		return "", xerror.Wrap("failed to write to buffer", err)
	}

	// Send email.
	output, err := SESClient.SendRawEmail(ctx, &ses.SendRawEmailInput{
		Destinations: destinations,
		FromArn:      aws.String(env.Get("AWS_FROM_ARN", "")),
		RawMessage: &types.RawMessage{
			Data: buf.Bytes(),
		},
	})
	if err != nil {
		return "", xerror.Wrap("failed to send raw email", err)
	}

	return *output.MessageId, nil
}
