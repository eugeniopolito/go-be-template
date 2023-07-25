package test

import (
	"context"
	"log"
	"os"
	"testing"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot open connection to db:", err)
	}
	testQueries = db.New(connPool)

	os.Exit(m.Run())
}
