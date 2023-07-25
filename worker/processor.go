package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/mail"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/hibiken/asynq"
	log "github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	config util.Config
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, rr error) {
			log.Error().Msg("process task failed")
		}),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		config: config,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessSendVerifyEmail)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshall payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == db.ErrRecordNotFound {
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
	subject := processor.config.VerifyEmailSubject
	content := fmt.Sprintf(processor.config.VerifyEmailBody, user.Name, user.Surname, verifyURL)

	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}
