package authenticationsvc

import (
	"context"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"

	"go.uber.org/zap"
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

func (s AuthenticationService) Authenticate(ctx context.Context, username, password string) (*domain.Member, *domain.Token, error) {
	token, err := s.authRepo.Authenticate(ctx, username, password)
	if err != nil {
		s.logger.Errorw("failed to authenticate", "error", err)
		return nil, nil, err
	}

	member, err := s.memberRepo.GetMemberByUsername(username)
	if err != nil {
		s.logger.Errorw("failed to get member by username", "error", err)
		return nil, nil, err
	}

	prefs, err := s.memberRepo.GetMemberPrefs(member.ID)
	if err != nil {
		s.logger.Errorw("failed to get member prefs", "error", err)
		return nil, nil, err
	}

	member.Prefs = *prefs

	return member, token, nil
}
