package authenticationsvc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"goBoard/internal/core/domain"
	"goBoard/internal/core/service/mocks"
	"testing"
)

func TestAuthenticationService_Authenticate(t *testing.T) {
	l := zap.NewNop().Sugar()

	t.Run("when everything goes right and a user is authenticated", func(t *testing.T) {
		var actualUsername string
		var actualMemberID int
		memberPrefs := domain.MemberPrefs{
			"hometown": domain.MemberPref{
				Value: "Colorado Springs",
			},
			"favorite_food": domain.MemberPref{
				Value: "curry",
			},
		}

		authRepoMock := &mocks.AuthenticationRepoMock{
			AuthenticateFunc: func(username, password string) (int, error) {
				return 666, nil
			},
		}

		memberRepoMock := &mocks.MemberRepoMock{
			GetMemberByUsernameFunc: func(username string) (*domain.Member, error) {
				actualUsername = username
				return &domain.Member{
					ID:      666,
					Name:    "test-name",
					Pass:    "test-pass",
					IsAdmin: true,
					Banned:  false,
				}, nil
			},
			GetMemberPrefsFunc: func(memberID int) (*domain.MemberPrefs, error) {
				actualMemberID = memberID
				return &memberPrefs, nil
			},
		}

		svc := NewAuthenticationService(authRepoMock, memberRepoMock, l)

		member, err := svc.Authenticate(nil, "test-name", "test-pass")
		require.NoError(t, err)

		assert.Equal(t, 666, member.ID)
		assert.True(t, member.IsAdmin)
		assert.Equal(t, memberPrefs, member.Prefs)
		assert.Equal(t, "test-name", actualUsername)
		assert.Equal(t, 666, actualMemberID)
	})

	t.Run("when everything goes right and a user is not authenticated", func(t *testing.T) {
		authRepoMock := &mocks.AuthenticationRepoMock{
			AuthenticateFunc: func(username, password string) (int, error) {
				return 0, nil
			},
		}

		memberRepoMock := &mocks.MemberRepoMock{}

		svc := NewAuthenticationService(authRepoMock, memberRepoMock, l)

		member, err := svc.Authenticate(context.Background(), "test-name", "test-pass")
		require.NoError(t, err)

		assert.Nil(t, member)
	})
}
