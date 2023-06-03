package membersvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"

	"go.uber.org/zap"
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

func (s MemberService) ValidateMembers(names []string) ([]domain.Member, error) {
	var validMembers []domain.Member
	for _, name := range names {
		id, err := s.memberRepo.GetMemberIDByUsername(name)
		if err != nil {
			s.logger.Errorf("error getting member id: %v", err)
			continue
		}

		member, err := s.memberRepo.GetMemberByID(id)
		if err != nil {
			s.logger.Errorf("error getting member: %v", err)
			return nil, err
		}

		if !member.Banned {
			validMembers = append(validMembers, *member)
		}
	}

	return validMembers, nil
}
