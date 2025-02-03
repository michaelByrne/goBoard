package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"time"

	"log/slog"
)

type AWSS3 interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type AWSPresign interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type Handler struct {
	awsS3      AWSS3
	awsPresign AWSPresign

	logger slog.Logger

	bucket string
}

func NewHandler(awsS3 AWSS3, awsPresign AWSPresign, logger slog.Logger, bucket string) *Handler {
	return &Handler{
		awsS3:      awsS3,
		awsPresign: awsPresign,
		logger:     logger,
		bucket:     bucket,
	}
}

func (h Handler) Handle(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	key, ok := event.QueryStringParameters["key"]
	if !ok {
		h.logger.Error("key query param is required")
		return events.APIGatewayProxyResponse{}, nil
	}

	input := &s3.GetObjectInput{
		Bucket: &h.bucket,
		Key:    &key,
	}

	presignOutput, err := h.awsPresign.PresignGetObject(ctx, input, s3.WithPresignExpires(time.Minute*15))
	if err != nil {
		h.logger.Error("error getting presigned url ", err)
		return events.APIGatewayProxyResponse{}, err
	}

	presignBytes, err := json.Marshal(presignOutput)
	if err != nil {
		h.logger.Error("error marshalling presigned body ", err)
		return events.APIGatewayProxyResponse{}, err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            string(presignBytes),
		IsBase64Encoded: false,
		Headers:         headers,
	}, nil
}

func main() {
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("BUCKET_NAME env var is required")
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	awsS3 := s3.NewFromConfig(awsConfig)
	awsPresign := s3.NewPresignClient(awsS3)

	logger := slog.Default()

	handler := NewHandler(awsS3, awsPresign, *logger, bucketName)

	lambda.Start(handler.Handle)
}
