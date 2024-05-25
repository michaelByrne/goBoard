package membersvc

import (
	"context"
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

func (s MemberService) GetAllPrefs(ctx context.Context) ([]domain.Pref, error) {
	return s.memberRepo.GetAllPrefs(ctx)
}

func (s MemberService) UpdatePrefs(ctx context.Context, memberID int, updatedPrefs domain.MemberPrefs) error {
	return s.memberRepo.UpdatePrefs(ctx, memberID, updatedPrefs)
}

func (s MemberService) GetMergedPrefs(ctx context.Context, memberID int) ([]domain.Pref, error) {
	memberPrefs, err := s.memberRepo.GetMemberPrefs(memberID)
	if err != nil {
		s.logger.Errorf("error getting member prefs: %v", err)
		return nil, err
	}

	allPrefs, err := s.memberRepo.GetAllPrefs(ctx)
	if err != nil {
		s.logger.Errorf("error getting all prefs: %v", err)
		return nil, err
	}

	mergedPrefs := mergePrefs(*memberPrefs, allPrefs)

	return mergedPrefs, nil
}

func (s MemberService) UpdateMember(ctx context.Context, member domain.Member) error {
	return s.memberRepo.UpdateMember(ctx, member)
}

func (s MemberService) UpdatePostalCode(ctx context.Context, memberID int, postalCode string) error {
	err := s.memberRepo.UpdatePostalCode(ctx, memberID, postalCode)
	if err != nil {
		s.logger.Errorf("error updating postal code: %v", err)
		return err
	}

	return nil
}

func mergePrefs(memberPrefs domain.MemberPrefs, prefs []domain.Pref) []domain.Pref {
	var prefsOut []domain.Pref

	for _, pref := range prefs {
		if value, ok := memberPrefs[pref.Name]; ok {
			prefsOut = append(prefsOut, domain.Pref{
				Name:    pref.Name,
				Value:   value.Value,
				Width:   pref.Width,
				Display: pref.Display,
				Type:    pref.Type,
			})
		} else {
			prefsOut = append(prefsOut, domain.Pref{
				Name:    pref.Name,
				Value:   pref.Value,
				Width:   pref.Width,
				Display: pref.Display,
				Type:    pref.Type,
			})
		}
	}

	return prefsOut
}
