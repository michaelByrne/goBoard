package membersvc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/service/mocks"
	"testing"
)

func TestMemberService_ValidateMembers(t *testing.T) {
	l := zap.NewNop()
	sugar := l.Sugar()

	t.Run("successfully returns only valid members", func(t *testing.T) {
		byUsernameCalls := 0
		memberRepoMock := &mocks.MemberRepoMock{
			GetMemberIDByUsernameFunc: func(username string) (int, error) {
				if byUsernameCalls == 0 {
					byUsernameCalls++
					return 1, nil
				}

				return 0, fmt.Errorf("error getting member id")
			},
			GetMemberByIDFunc: func(id int) (*domain.Member, error) {
				return &domain.Member{
					ID:     1,
					Name:   "tester",
					Banned: false,
				}, nil
			},
		}

		memberService := NewMemberService(memberRepoMock, sugar)

		members, err := memberService.ValidateMembers([]string{"tester", "tester2"})
		require.NoError(t, err)

		require.Equal(t, 1, len(members))
		require.Equal(t, "tester", members[0].Name)
		require.Equal(t, 1, members[0].ID)
	})

	t.Run("rejects a banned member", func(t *testing.T) {
		byUsernameCalls := 0
		calledWithID := 0
		memberRepoMock := &mocks.MemberRepoMock{
			GetMemberIDByUsernameFunc: func(username string) (int, error) {
				if byUsernameCalls == 0 {
					byUsernameCalls++
					return 1, nil
				}

				return 0, fmt.Errorf("error getting member id")
			},
			GetMemberByIDFunc: func(id int) (*domain.Member, error) {
				calledWithID = id
				return &domain.Member{
					ID:     1,
					Name:   "tester",
					Banned: true,
				}, nil
			},
		}

		memberService := NewMemberService(memberRepoMock, sugar)

		members, err := memberService.ValidateMembers([]string{"tester", "tester2"})
		require.NoError(t, err)

		require.Len(t, members, 0)
		require.Equal(t, 1, calledWithID)
	})

	t.Run("ignores members that don't exist", func(t *testing.T) {
		memberRepoMock := &mocks.MemberRepoMock{
			GetMemberIDByUsernameFunc: func(username string) (int, error) {
				if username == "tester" {
					return 1, nil
				}

				return 0, fmt.Errorf("error getting member id")
			},

			GetMemberByIDFunc: func(id int) (*domain.Member, error) {
				if id == 1 {
					return &domain.Member{
						ID:     1,
						Name:   "tester1",
						Banned: false,
					}, nil
				}
				return nil, fmt.Errorf("error getting member")
			},
		}

		memberService := NewMemberService(memberRepoMock, sugar)

		members, err := memberService.ValidateMembers([]string{"tester", "tester2"})
		require.NoError(t, err)

		require.Len(t, members, 1)
		require.Equal(t, "tester1", members[0].Name)
	})
}

func TestMemberService_GetMergedPrefs(t *testing.T) {
	l := zap.NewNop().Sugar()

	t.Run("when everything goes right and merged prefs are returned", func(t *testing.T) {
		fakeMemberPrefs := &domain.MemberPrefs{
			"theme": domain.MemberPref{
				Value: "dark",
				Type:  "checkbox",
			},
			"location": domain.MemberPref{
				Value: "US",
				Type:  "text",
			},
		}

		fakeAllPrefs := []domain.Pref{
			{
				Name:    "theme",
				Display: "visual theme",
				Type:    "checkbox",
			},
			{
				Name:    "location",
				Display: "location",
				Type:    "text",
			},
			{
				Name:    "timezone",
				Display: "timezone",
				Type:    "text",
			},
		}

		mockMemberRepo := &mocks.MemberRepoMock{
			GetMemberPrefsFunc: func(id int) (*domain.MemberPrefs, error) {
				return fakeMemberPrefs, nil
			},
			GetAllPrefsFunc: func(ctx context.Context) ([]domain.Pref, error) {
				return fakeAllPrefs, nil
			},
		}

		expectedPrefs := []domain.Pref{
			{
				Name:    "theme",
				Display: "visual theme",
				Type:    "checkbox",
				Value:   "dark",
			},
			{
				Name:    "location",
				Display: "location",
				Type:    "text",
				Value:   "US",
			},
			{
				Name:    "timezone",
				Display: "timezone",
				Type:    "text",
				Value:   "",
			},
		}

		memberService := NewMemberService(mockMemberRepo, l)

		mergedPrefs, err := memberService.GetMergedPrefs(context.Background(), 1)
		require.NoError(t, err)

		require.Len(t, mergedPrefs, 3)
		assert.Equal(t, expectedPrefs, mergedPrefs)
	})

	t.Run("when a user doesn't have any prefs yet", func(t *testing.T) {
		fakeAllPrefs := []domain.Pref{
			{
				Name:    "theme",
				Display: "visual theme",
				Type:    "checkbox",
			},
			{
				Name:    "location",
				Display: "location",
				Type:    "text",
			},
			{
				Name:    "timezone",
				Display: "timezone",
				Type:    "text",
			},
		}

		emptyPrefs := make(domain.MemberPrefs)

		mockMemberRepo := &mocks.MemberRepoMock{
			GetMemberPrefsFunc: func(id int) (*domain.MemberPrefs, error) {
				return &emptyPrefs, nil
			},
			GetAllPrefsFunc: func(ctx context.Context) ([]domain.Pref, error) {
				return fakeAllPrefs, nil
			},
		}

		memberService := NewMemberService(mockMemberRepo, l)

		mergedPrefs, err := memberService.GetMergedPrefs(context.Background(), 1)
		require.NoError(t, err)

		require.Len(t, mergedPrefs, 3)
		assert.Equal(t, fakeAllPrefs, mergedPrefs)
	})
}
