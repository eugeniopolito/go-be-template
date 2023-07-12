package api

import (
	"fmt"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/token"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/eugeniopolito/gobetemplate/worker"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

const userkey = "user"

var secret = []byte("S3CrEt!25#")

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(str string) gin.H {
	return gin.H{"error": str}
}
