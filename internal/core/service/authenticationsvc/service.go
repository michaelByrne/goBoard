package authenticationsvc

import (
	"context"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type AuthenticationService struct {
	logger     *zap.SugaredLogger
	authRepo   ports.AuthenticationRepo
	memberRepo ports.MemberRepo
}

func NewAuthenticationService(authRepo ports.AuthenticationRepo, memberRepo ports.MemberRepo, logger *zap.SugaredLogger) AuthenticationService {
	return AuthenticationService{
		logger:     logger,
		authRepo:   authRepo,
		memberRepo: memberRepo,
	}
}

func (s AuthenticationService) Authenticate(ctx context.Context, username, password string) (*domain.Member, error) {
	id, err := s.authRepo.Authenticate(username, password)
	if err != nil {
		s.logger.Errorw("failed to authenticate", "error", err)
		return nil, err
	}

	if id == 0 {
		s.logger.Infow("authentication failed", "username", username)
		return nil, nil
	}

	member, err := s.memberRepo.GetMemberByUsername(username)
	if err != nil {
		s.logger.Errorw("failed to get member by username", "error", err)
		return nil, err
	}

	prefs, err := s.memberRepo.GetMemberPrefs(member.ID)
	if err != nil {
		s.logger.Errorw("failed to get member prefs", "error", err)
		return nil, err
	}

	member.Prefs = *prefs

	return member, nil
}
