package test

import (
	"testing"
	"time"

	"github.com/eugeniopolito/gobetemplate/token"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, util.RandomRole(), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJPasetoMaker(test *testing.T) {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(test, err)

	t, err := maker.CreateToken(util.RandomOwner(), util.RandomRole(), -time.Minute)
	require.NoError(test, err)
	require.NotEmpty(test, t)

	payload, err := maker.VerifyToken(t)
	require.Error(test, err)
	require.EqualError(test, err, token.ErrExpiredToken.Error())

	require.Nil(test, payload)
}
