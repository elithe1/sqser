package sqsercore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func getSession() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), "")
	options := session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		Config: aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: creds,
		},
	}
	return session.NewSessionWithOptions(options)

}
