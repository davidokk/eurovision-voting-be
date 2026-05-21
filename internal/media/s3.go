package media

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3 struct {
	client *s3.Client
	bucket string
	public string
}

func NewS3FromEnv(ctx context.Context) (*S3, error) {
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET not set")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "auto"
	}
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, err
	}
	endpoint := strings.TrimRight(os.Getenv("S3_ENDPOINT"), "/")
	var s3Opts []func(*s3.Options)
	if endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = os.Getenv("S3_USE_PATH_STYLE") == "true"
		})
	}
	public := os.Getenv("S3_PUBLIC_BASE_URL")
	if public == "" {
		if endpoint != "" {
			public = endpoint + "/" + bucket
		} else {
			public = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket, region)
		}
	}
	public = strings.TrimRight(public, "/")
	return &S3{
		client: s3.NewFromConfig(cfg, s3Opts...),
		bucket: bucket,
		public: public,
	}, nil
}

func (s *S3) Upload(ctx context.Context, folder string, data []byte, contentType string) (string, error) {
	ext := extForContentType(contentType)
	key := path.Join("media", folder, uuid.New().String()+ext)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return s.public + "/" + key, nil
}

func extForContentType(ct string) string {
	switch {
	case strings.Contains(ct, "jpeg"), strings.Contains(ct, "jpg"):
		return ".jpg"
	case strings.Contains(ct, "png"):
		return ".png"
	case strings.Contains(ct, "webm") && strings.Contains(ct, "video"):
		return ".webm"
	case strings.Contains(ct, "webm"):
		return ".webm"
	case strings.Contains(ct, "mp4"):
		return ".mp4"
	case strings.Contains(ct, "mpeg"), strings.Contains(ct, "mp3"):
		return ".mp3"
	default:
		return ".bin"
	}
}
