package api

import (
	docs "github.com/eugeniopolito/gobetemplate/docs"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// The App routes
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
	authRoutes.Use(AuthMiddleware(server.tokenMaker, false, 0))
	authRoutes.POST("/users/logout", server.logoutUser)
	authRoutes.GET("/user/:username", server.getUser)

	// admin, authenticated
	authAdminRoutes := router.Group("/v1/admin")
	authAdminRoutes.Use(DefaultStructuredLogger())
	authAdminRoutes.Use(AuthMiddleware(server.tokenMaker, true, util.ADMIN))
	authAdminRoutes.GET("/users", server.listUsers)
	authAdminRoutes.GET("/users/count", server.countUsers)

	server.router = router
}
