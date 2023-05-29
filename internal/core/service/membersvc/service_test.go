package membersvc

import (
	"fmt"
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
