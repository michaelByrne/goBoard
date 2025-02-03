package imagerepo

import (
	"bytes"
	"context"
	"fmt"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"goBoard/internal/gberrors"
	"image"
	"image/jpeg"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

type ImageRepo struct {
	awsS3      ports.AWSS3
	awsPresign ports.AWSPresign

	logger zap.SugaredLogger

	bucket string
}

func NewImageRepo(awsS3 ports.AWSS3, awsPresign ports.AWSPresign, logger zap.SugaredLogger, bucket string) *ImageRepo {
	return &ImageRepo{
		awsS3:      awsS3,
		awsPresign: awsPresign,
		logger:     logger,
		bucket:     bucket,
	}
}

func (r ImageRepo) PresignURL(ctx context.Context, key string) (string, error) {
	presignURL, err := r.awsPresign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    &key,
	}, s3.WithPresignExpires(time.Minute*1))
	if err != nil {
		r.logger.Errorw("failed to presign url", "error", err)
		return "", err
	}

	return presignURL.URL, nil
}

func (r ImageRepo) UploadImage(ctx context.Context, imageBytes []byte) (*domain.Image, error) {
	imageReader := bytes.NewReader(imageBytes)

	key := uuid.New().String() + ".jpg"

	_, err := r.awsS3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &r.bucket,
		Key:    &key,
		Body:   imageReader,
	})
	if err != nil {
		r.logger.Errorw("failed to upload image", "error", err)
		return nil, err
	}

	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", r.bucket, key)

	return &domain.Image{
		URL:  imageURL,
		Name: key,
	}, nil

}

func (r ImageRepo) ResizeImage(imageBytes []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		if format != "jpeg" {
			return nil, gberrors.UnsupportedImageFormat{
				Format: format,
			}
		}

		r.logger.Errorw("failed to decode image", "error", err)
		return nil, err
	}

	if format != "jpeg" {
		r.logger.Errorw("image format is not jpeg", "format", format)
		return nil, gberrors.UnsupportedImageFormat{
			Format: format,
		}
	}

	// newWidth and newHeight are now adjusted to be within the range of 500px and 700px
	resizedImg := resize.Resize(550, 0, img, resize.Lanczos3)

	resizedImgBytes, err := imgToBytes(resizedImg)
	if err != nil {
		r.logger.Errorw("failed to convert image to bytes", "error", err)
		return nil, err
	}

	return resizedImgBytes, nil
}

func imgToBytes(img image.Image) ([]byte, error) {
	var opt jpeg.Options
	opt.Quality = 80

	buff := bytes.NewBuffer(nil)
	err := jpeg.Encode(buff, img, &opt)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
