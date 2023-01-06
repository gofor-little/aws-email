## A package for sending emails via AWS SES

![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/gofor-little/aws-email?include_prereleases)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gofor-little/aws-email)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://raw.githubusercontent.com/gofor-little/aws-email/main/LICENSE)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/gofor-little/aws-email/ci.yaml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofor-little/aws-email)](https://goreportcard.com/report/github.com/gofor-little/aws-email)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gofor-little/aws-email)](https://pkg.go.dev/github.com/gofor-little/aws-email)

### Introduction
* Send emails via AWS SES
* Fetch attachments locally, from an S3 Bucket or via HTTP
* Supports both text/plain and text/html emails

### Example
```go
package main

import (
	"context"

	email "github.com/gofor-little/aws-email"
)

func main() {
	// Initialize the email package.
	if err := email.Initialize("AWS_PROFILE", "AWS_REGION"); err != nil {
		panic(err)
	}

	// Build the email data.
	data := email.Data{
		To:          []string{"joe@example.com"},
		From:        "john@example.com",
		Subject:     "Example Email",
		Body:        "Example Body",
		ContentType: email.ContentTypeTextPlain,
		Attachments: []email.Attachment{
			{
				Path: "/home/ubuntu/app/logo.png",
				Type: email.AttachmentTypeLocal,
			},
			{
				Path: "bucket-name/images/logo.png",
				Type: email.AttachmentTypeS3,
			},
			{
				Path: "https://example.com/logo.png",
				Type: email.AttachmentTypeHTTP,
			},
		},
	}

	// Send the email.
	if _, err := email.Send(context.Background(), data); err != nil {
		panic(err)
	}
}
```

### Testing
Ensure the following environment variables are set, usually with a .env file.
* ```AWS_PROFILE``` (an AWS CLI profile name)
* ```AWS_REGION``` (a valid AWS region)
* ```TEST_EMAIL_FROM``` (a valid email address verified in SES)
* ```TEST_EMAIL_RECIPIENTS``` (a comma separated list of valid email addresses)
