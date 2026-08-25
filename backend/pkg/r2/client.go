package r2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client wraps an S3-compatible client pointed at Cloudflare R2.
// Used for "I made this" photos and thumbnails.
type Client struct {
	s3     *s3.Client
	bucket string
}

func New(ctx context.Context, accountID, accessKeyID, secretAccessKey, bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	endpoint := "https://" + accountID + ".r2.cloudflarestorage.com"

	return &Client{
		s3:     s3.NewFromConfig(cfg, func(o *s3.Options) { o.BaseEndpoint = &endpoint }),
		bucket: bucket,
	}, nil
}
