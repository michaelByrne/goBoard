package membersvc

import (
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type MemberService struct {
	memberRepo ports.MemberRepo
	logger     *zap.SugaredLogger
}

func NewMemberService(memberRepo ports.MemberRepo, logger *zap.SugaredLogger) MemberService {
	return MemberService{memberRepo, logger}
}

func (s MemberService) Save(member domain.Member) (int, error) {
	id, err := s.memberRepo.SaveMember(member)
	if err != nil {
		s.logger.Errorf("error saving member: %v", err)
		return 0, err
	}

	return id, nil
}

func (s MemberService) GetMemberByID(id int) (*domain.Member, error) {
	return s.memberRepo.GetMemberByID(id)
}

func (s MemberService) GetMemberIDByUsername(username string) (int, error) {
	return s.memberRepo.GetMemberIDByUsername(username)
}

func (s MemberService) GetMemberByUsername(username string) (*domain.Member, error) {
	return s.memberRepo.GetMemberByUsername(username)
}
