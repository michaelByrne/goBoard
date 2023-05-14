package memberrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goBoard/db"
	"goBoard/internal/core/domain"
	"goBoard/internal/repos/seed"
	"testing"
)

func TestNewMemberRepo(t *testing.T) {
	dbContainer, connPool, err := db.SetupTestDatabase()
	require.NoError(t, err)

	defer dbContainer.Terminate(context.Background())

	require.NoError(t, seed.SeedData(t, connPool))

	repo := NewMemberRepo(connPool)

	t.Run("should successfully save a member", func(t *testing.T) {
		member := domain.Member{
			Name:       "roxy",
			Pass:       "test",
			Secret:     "test",
			Email:      "roxy@gmail.com",
			PostalCode: "48225",
			IP:         "127.0.0.1",
		}

		returnedID, err := repo.SaveMember(member)
		require.NoError(t, err)

		var id int
		var name string
		var email string
		err = connPool.QueryRow(context.Background(), "SELECT id, name, email_signup FROM member WHERE id = $1", returnedID).Scan(&id, &name, &email)
		require.NoError(t, err)

		assert.Equal(t, id, returnedID)
		assert.Equal(t, name, member.Name)
		assert.Equal(t, email, member.Email)
	})

	t.Run("should successfully get a member by id", func(t *testing.T) {
		member, err := repo.GetMemberByID(2)
		require.NoError(t, err)

		assert.Equal(t, member.ID, 2)
		assert.Equal(t, member.Name, "gofreescout")
		assert.Equal(t, member.Email, "gofreescout@gmail.com")
		assert.Equal(t, member.PostalCode, "48225")
	})
}
