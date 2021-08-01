package email

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofor-little/xerror"
)

// AttachmentType is a string type alias for the different supported attachment types.
type AttachmentType string

const (
	// AttachmentTypeLocal is used for local attachments. As this is
	// deployed as a Lambda function it's primary use is testing.
	AttachmentTypeLocal AttachmentType = "local"
	// AttachmentTypeS3 is used for attachments stored in S3.
	AttachmentTypeS3 AttachmentType = "s3"
	// AttachmentTypeHTTP is used for attachments publicly accessible via HTTP.
	AttachmentTypeHTTP AttachmentType = "http"
)

// Attachment stores a path and type which are used to load it locally.
type Attachment struct {
	Path string         `json:"path"`
	Type AttachmentType `json:"type"`
}

// Load loads an Attachment's data and returns it as a byte slice.
func (a *Attachment) Load(ctx context.Context) ([]byte, error) {
	switch a.Type {
	case AttachmentTypeLocal:
		data, err := os.ReadFile(a.Path)
		if err != nil {
			return nil, xerror.Wrap("failed to read local file", err)
		}

		return data, nil
	case AttachmentTypeS3:
		// Split the path up to get the bucket and key.
		parts := strings.Split(a.Path, "/")
		if len(parts) < 2 {
			return nil, xerror.Newf("invalid path for AttachmentTypeS3: %s, path must be in the following format <bucket>/<key>", a.Path)
		}

		// Get the object from S3 and read it into a bytes.Buffer.
		output, err := S3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(parts[0]),
			Key:    aws.String(a.Path[len(parts[0])+1:]),
		})
		if err != nil {
			return nil, xerror.Wrap("failed to get object from S3", err)
		}

		buffer := bytes.Buffer{}
		_, err = buffer.ReadFrom(output.Body)
		if err != nil {
			return nil, xerror.Wrap("failed to read s3.GetObjectOutput's body", err)
		}

		return buffer.Bytes(), nil
	case AttachmentTypeHTTP:
		if HTTPClient == nil {
			HTTPClient = http.DefaultClient
		}

		// Create new http.Request.
		request, err := http.NewRequest(http.MethodGet, a.Path, nil)
		if err != nil {
			return nil, xerror.Wrap("failed to create http.Request", err)
		}

		// Get the Attachment data.
		response, err := HTTPClient.Do(request)
		if err != nil {
			return nil, xerror.Wrap("failed to do HTTP request on attachment", err)
		}
		defer func() {
			_ = response.Body.Close()
		}()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, xerror.Wrap("failed to read response body", err)
		}

		return data, nil
	default:
		// Fallthrough to the error case as a.Type is not a supported AttachmentType.
	}

	return nil, fmt.Errorf("unable to load Attachment with unexpected type: %s", a.Type)
}
