package test

import (
	"context"
	"testing"
	"time"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5/pgtype"
)

func createRandomUser(t *testing.T) db.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Username: util.RandomString(15),
		Name:     util.RandomString(10),
		Surname:  util.RandomString(10),
		Enabled:  true,
		Role:     pgtype.Int4{Int32: int32(util.RandomRole()), Valid: true},
		Email:    util.RandomEmail(),
		Password: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Surname, user.Surname)
	require.Equal(t, arg.Enabled, user.Enabled)
	require.Equal(t, arg.Role, user.Role)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.PasswordChangeAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Surname, user2.Surname)
	require.Equal(t, user1.Enabled, user2.Enabled)
	require.Equal(t, user1.Role, user2.Role)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Username)

	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.Error(t, err)
	require.EqualError(t, err, "no rows in result set")
	require.Empty(t, user2)
}

func TestUpdateUserEmail(t *testing.T) {
	user1 := createRandomUser(t)

	nm := util.RandomEmail()

	arg := db.UpdateUserEmailParams{
		Email:    nm,
		Username: user1.Username,
	}

	user2, err := testQueries.UpdateUserEmail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.Name, user1.Name)
	require.Equal(t, user2.Surname, user1.Surname)
	require.Equal(t, user2.Enabled, user1.Enabled)
	require.Equal(t, user2.Role, user1.Role)
	require.Equal(t, user2.Email, nm)
	require.Equal(t, user2.Password, user1.Password)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, time.Second)
}
