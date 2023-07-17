package main

import (
	"context"

	"github.com/eugeniopolito/gobetemplate/api"
	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/mail"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/eugeniopolito/gobetemplate/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("Cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Error().Err(err).Msg("Cannot connect to db")
	}

	// run DB migration
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)
	runGinServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot run migration")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("Cannot run migrate up")
	}
	log.Info().Msg("DB migration done")
}

func runGinServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Err(err).Msg("Cannot start server")
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewEmailConfigSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword, config.SmtpAuthAddress, config.SmtpServerAddress)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, config, mailer)
	log.Info().Msg("Start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Error().Err(err).Msg("Cannot start Redis processor")
	}
}
