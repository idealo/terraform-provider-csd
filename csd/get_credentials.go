package csd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func getCreds(p string, r string) (aws.Credentials, error) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(r),
		config.WithSharedConfigProfile(p),
	)
	if err != nil {
		var creds aws.Credentials
		return creds, err
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return creds, err
	}

	return creds, nil
}
