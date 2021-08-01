package email

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

var (
	// SESClient is used to send emails via SES.
	SESClient *ses.Client
	// S3Client is used to fetch attachments from S3.
	S3Client *s3.Client
	// HTTPClient is used to fetch attachments from the web.
	HTTPClient httpClient
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Initialize will initialize the sms package. Both the profile
// and region parameters are optional if authentication can be achieved
// via another method. For example, environment variables or IAM roles.
func Initialize(ctx context.Context, profile string, region string) error {
	var cfg aws.Config
	var err error

	if profile != "" && region != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile), config.WithRegion(region))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}
	if err != nil {
		return fmt.Errorf("failed to load default config: %w", err)
	}

	SESClient = ses.NewFromConfig(cfg)
	S3Client = s3.NewFromConfig(cfg)

	return nil
}
