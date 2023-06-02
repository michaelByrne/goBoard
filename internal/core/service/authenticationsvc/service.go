package authenticationsvc

import (
	"context"
	"go.uber.org/zap"
	"goBoard/internal/core/ports"
)

type AuthenticationService struct {
	logger   *zap.SugaredLogger
	authRepo ports.AuthenticationRepo
}

func NewAuthenticationService(authRepo ports.AuthenticationRepo, logger *zap.SugaredLogger) AuthenticationService {
	return AuthenticationService{
		logger:   logger,
		authRepo: authRepo,
	}
}

func (s AuthenticationService) Authenticate(ctx context.Context, username, password string) (int, error) {
	id, err := s.authRepo.Authenticate(username, password)
	if err != nil {
		s.logger.Errorw("failed to authenticate", "error", err)
		return 0, err
	}

	if id == 0 {
		s.logger.Errorw("authentication failed", "username", username)
		return 0, nil
	}

	return id, nil
}
