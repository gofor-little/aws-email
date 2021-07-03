package email

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/gofor-little/xerror"
)

var (
	// SESClient is used to send emails via SES.
	SESClient sesiface.SESAPI
	// S3Client is used to fetch attachments from S3.
	S3Client s3iface.S3API
	// HTTPClient is used to fetch attachments from the web.
	HTTPClient httpClient
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Initialize will initialize the sms package. Both the profile
// and region parameters are optional if authentication can be achieved
// via another method. For example, environment variables or IAM roles.
func Initialize(profile string, region string) error {
	var sess *session.Session
	var err error

	if profile != "" && region != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
			Profile: profile,
		})
	} else {
		sess, err = session.NewSession()
	}
	if err != nil {
		return xerror.Wrap("failed to create session.Session", err)
	}

	SESClient = ses.New(sess)
	S3Client = s3.New(sess)

	return nil
}
