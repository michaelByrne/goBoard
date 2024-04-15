package imagesvc

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"goBoard/internal/core/ports"
	"strings"
)

type ImageService struct {
	imageRepo ports.ImageRepo

	logger zap.SugaredLogger
}

func NewImageService(imageRepo ports.ImageRepo, logger zap.SugaredLogger) *ImageService {
	return &ImageService{
		imageRepo: imageRepo,
		logger:    logger,
	}
}

func (s ImageService) UploadImage(ctx context.Context, imageBytes []byte) (string, string, error) {
	resized, err := s.imageRepo.ResizeImage(imageBytes)
	if err != nil {
		s.logger.Errorw("failed to resize image", "error", err)
		return "", "", err
	}

	uploadedImage, err := s.imageRepo.UploadImage(ctx, resized)
	if err != nil {
		s.logger.Errorw("failed to upload image", "error", err)
		return "", "", err
	}

	presignedURL, err := s.imageRepo.PresignURL(ctx, uploadedImage.Name)
	if err != nil {
		s.logger.Errorw("failed to presign url", "error", err)
		return "", "", err
	}

	return uploadedImage.Name, presignedURL, nil
}

func (s ImageService) RefreshPresign(ctx context.Context, key string) (string, error) {
	presignedURL, err := s.imageRepo.PresignURL(ctx, key)
	if err != nil {
		s.logger.Errorw("failed to presign url", "error", err)
		return "", err
	}

	return presignedURL, nil
}

func (s ImageService) PresignPostImages(ctx context.Context, body string) (string, error) {
	bodyReader := strings.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		s.logger.Errorw("failed to parse post body", "error", err)
		return "", err
	}

	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		key, exists := selection.Attr("data-key")
		if !exists {
			return
		}

		presignedURL, err := s.imageRepo.PresignURL(ctx, key)
		if err != nil {
			s.logger.Errorw("failed to presign url", "error", err)
			return
		}

		selection.SetAttr("src", presignedURL)
	})

	html, err := doc.Html()
	if err != nil {
		s.logger.Errorw("failed to convert post body to html", "error", err)
		return "", err
	}

	return html, nil
}
