// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
	"sync"
)

// Ensure, that ImageRepoMock does implement ports.ImageRepo.
// If this is not the case, regenerate this file with moq.
var _ ports.ImageRepo = &ImageRepoMock{}

// ImageRepoMock is a mock implementation of ports.ImageRepo.
//
//	func TestSomethingThatUsesImageRepo(t *testing.T) {
//
//		// make and configure a mocked ports.ImageRepo
//		mockedImageRepo := &ImageRepoMock{
//			PresignURLFunc: func(ctx context.Context, key string) (string, error) {
//				panic("mock out the PresignURL method")
//			},
//			ResizeImageFunc: func(imageBytes []byte) ([]byte, error) {
//				panic("mock out the ResizeImage method")
//			},
//			UploadImageFunc: func(ctx context.Context, imageBytes []byte) (*domain.Image, error) {
//				panic("mock out the UploadImage method")
//			},
//		}
//
//		// use mockedImageRepo in code that requires ports.ImageRepo
//		// and then make assertions.
//
//	}
type ImageRepoMock struct {
	// PresignURLFunc mocks the PresignURL method.
	PresignURLFunc func(ctx context.Context, key string) (string, error)

	// ResizeImageFunc mocks the ResizeImage method.
	ResizeImageFunc func(imageBytes []byte) ([]byte, error)

	// UploadImageFunc mocks the UploadImage method.
	UploadImageFunc func(ctx context.Context, imageBytes []byte) (*domain.Image, error)

	// calls tracks calls to the methods.
	calls struct {
		// PresignURL holds details about calls to the PresignURL method.
		PresignURL []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
		}
		// ResizeImage holds details about calls to the ResizeImage method.
		ResizeImage []struct {
			// ImageBytes is the imageBytes argument value.
			ImageBytes []byte
		}
		// UploadImage holds details about calls to the UploadImage method.
		UploadImage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ImageBytes is the imageBytes argument value.
			ImageBytes []byte
		}
	}
	lockPresignURL  sync.RWMutex
	lockResizeImage sync.RWMutex
	lockUploadImage sync.RWMutex
}

// PresignURL calls PresignURLFunc.
func (mock *ImageRepoMock) PresignURL(ctx context.Context, key string) (string, error) {
	if mock.PresignURLFunc == nil {
		panic("ImageRepoMock.PresignURLFunc: method is nil but ImageRepo.PresignURL was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockPresignURL.Lock()
	mock.calls.PresignURL = append(mock.calls.PresignURL, callInfo)
	mock.lockPresignURL.Unlock()
	return mock.PresignURLFunc(ctx, key)
}

// PresignURLCalls gets all the calls that were made to PresignURL.
// Check the length with:
//
//	len(mockedImageRepo.PresignURLCalls())
func (mock *ImageRepoMock) PresignURLCalls() []struct {
	Ctx context.Context
	Key string
} {
	var calls []struct {
		Ctx context.Context
		Key string
	}
	mock.lockPresignURL.RLock()
	calls = mock.calls.PresignURL
	mock.lockPresignURL.RUnlock()
	return calls
}

// ResizeImage calls ResizeImageFunc.
func (mock *ImageRepoMock) ResizeImage(imageBytes []byte) ([]byte, error) {
	if mock.ResizeImageFunc == nil {
		panic("ImageRepoMock.ResizeImageFunc: method is nil but ImageRepo.ResizeImage was just called")
	}
	callInfo := struct {
		ImageBytes []byte
	}{
		ImageBytes: imageBytes,
	}
	mock.lockResizeImage.Lock()
	mock.calls.ResizeImage = append(mock.calls.ResizeImage, callInfo)
	mock.lockResizeImage.Unlock()
	return mock.ResizeImageFunc(imageBytes)
}

// ResizeImageCalls gets all the calls that were made to ResizeImage.
// Check the length with:
//
//	len(mockedImageRepo.ResizeImageCalls())
func (mock *ImageRepoMock) ResizeImageCalls() []struct {
	ImageBytes []byte
} {
	var calls []struct {
		ImageBytes []byte
	}
	mock.lockResizeImage.RLock()
	calls = mock.calls.ResizeImage
	mock.lockResizeImage.RUnlock()
	return calls
}

// UploadImage calls UploadImageFunc.
func (mock *ImageRepoMock) UploadImage(ctx context.Context, imageBytes []byte) (*domain.Image, error) {
	if mock.UploadImageFunc == nil {
		panic("ImageRepoMock.UploadImageFunc: method is nil but ImageRepo.UploadImage was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		ImageBytes []byte
	}{
		Ctx:        ctx,
		ImageBytes: imageBytes,
	}
	mock.lockUploadImage.Lock()
	mock.calls.UploadImage = append(mock.calls.UploadImage, callInfo)
	mock.lockUploadImage.Unlock()
	return mock.UploadImageFunc(ctx, imageBytes)
}

// UploadImageCalls gets all the calls that were made to UploadImage.
// Check the length with:
//
//	len(mockedImageRepo.UploadImageCalls())
func (mock *ImageRepoMock) UploadImageCalls() []struct {
	Ctx        context.Context
	ImageBytes []byte
} {
	var calls []struct {
		Ctx        context.Context
		ImageBytes []byte
	}
	mock.lockUploadImage.RLock()
	calls = mock.calls.UploadImage
	mock.lockUploadImage.RUnlock()
	return calls
}