package worker

import (
	"context"

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
