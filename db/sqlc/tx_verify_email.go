package db

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		log.Info().Int64("email_id", arg.EmailId).Str("secret_code", arg.SecretCode).Msg("verify email")

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			log.Err(err).Msg("cannot verify email")
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		if err != nil {
			log.Err(err).Msg("cannot update user")
			return err
		} else {
			log.Info().Str("username", result.User.Username).Msg("successfully verified")
		}
		return err
	})

	return result, err
}
