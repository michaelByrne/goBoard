package membersvc

import (
	"goBoard/internal/core/domain"
	"goBoard/internal/core/ports"
)

type MemberService struct {
	memberRepo ports.MemberRepo
}

func NewMemberService(memberRepo ports.MemberRepo) MemberService {
	return MemberService{memberRepo}
}

func (s MemberService) Save(member domain.Member) error {
	_, err := s.memberRepo.SaveMember(member)
	if err != nil {
		return err
	}

	return nil
}

func (s MemberService) GetMemberByID(id int) (*domain.Member, error) {
	return s.memberRepo.GetMemberByID(id)
}
