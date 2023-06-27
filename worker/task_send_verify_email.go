package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/hibiken/asynq"
	log "github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const TaskSendVerifyEmail = "task:send_verify_email"

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPaylod, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPaylod, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshall payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", asynq.SkipRetry)
	}

	arg := db.CreateVerifyEmailParams{
		Username:   payload.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, arg)

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", asynq.SkipRetry)
	}

	verifyURL := fmt.Sprintf(processor.config.VerifyEmailAddress+"?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	subject := "Welcome to GO BE Template!"
	content := fmt.Sprintf(`Hello %s %s,<br>

	Thank you for registering with us.<br>
	Please <a href="%s">click here</a> to verify your email address.<br>`, user.Name, user.Surname, verifyURL)

	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}
