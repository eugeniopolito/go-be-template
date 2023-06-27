package api

import (
	"fmt"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	docs "github.com/eugeniopolito/gobetemplate/docs"
	"github.com/eugeniopolito/gobetemplate/token"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/eugeniopolito/gobetemplate/worker"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

type paginationRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=1,max=100"`
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

func (server *Server) setupRouter() {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/v1"

	// initialize sessions
	router.Use(sessions.Sessions("session", cookie.NewStore(secret)))
	router.Use(DefaultStructuredLogger())

	// user, public
	router.POST("/v1/users", server.createUser)
	router.GET("/v1/verify_email", server.verifyEmail)
	router.POST("/v1/users/login", server.loginUser)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// user, authenticated
	authRoutes := router.Group("/v1")
	authRoutes.Use(DefaultStructuredLogger())
	authRoutes.Use(authMiddleware(server.tokenMaker, false, 0))
	authRoutes.POST("/users/logout", server.logoutUser)
	authRoutes.GET("/user/:username", server.getUser)

	// admin, authenticated
	authAdminRoutes := router.Group("/v1/admin")
	authAdminRoutes.Use(DefaultStructuredLogger())
	authAdminRoutes.Use(authMiddleware(server.tokenMaker, true, util.ADMIN))
	authAdminRoutes.GET("/users", server.listUsers)
	authAdminRoutes.GET("/users/count", server.countUsers)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(str string) gin.H {
	return gin.H{"error": str}
}
