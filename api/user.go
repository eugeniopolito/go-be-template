package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	db "github.com/eugeniopolito/gobetemplate/db/sqlc"
	"github.com/eugeniopolito/gobetemplate/token"
	"github.com/eugeniopolito/gobetemplate/util"
	"github.com/eugeniopolito/gobetemplate/worker"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

// @BasePath /api/v1

// Create User godoc
// @Summary create a new user
// @Schemes
// @Description creates a new user who receives a verification email on his/her email address to confirm the registration.
// @Tags users
// @Accept json
// @Produce json
// @Param req body CreateUserRequest true "CreateUserRequest"
// @Success 200 {object} UserResponse
// @Router /users [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Internal Server Error"))
		return
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username: req.Username,
			Email:    req.Email,
			Name:     req.Name,
			Role:     pgtype.Int4{Int32: int32(req.Role), Valid: true},
			Surname:  req.Surname,
			Password: hashedPassword,
			Enabled:  true,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if strings.Contains(err.Error(), "users_pkey") {
			log.Error().Str("user", req.Username).Msg("User already exists")
			ctx.JSON(http.StatusBadRequest, errorResponse(("User already exists")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	rsp := CreateUserResponse(txResult.User)

	ctx.JSON(http.StatusCreated, rsp)
}

// @BasePath /api/v1

// Verify the user godoc
// @Summary perform the user verification with email check
// @Schemes
// @Description check the code received in the email during registration is correct
// @Tags users
// @Accept json
// @Produce json
// @Param req body VerifyEmailRequest true "VerifyEmailRequest"
// @Success 200 {object} VerifyEmailResponse
// @Failure 500 "failed to verify email"
// @Router /verify_email [get]
func (server *Server) verifyEmail(ctx *gin.Context) {
	var req VerifyEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	arg := db.VerifyEmailTxParams{
		EmailId:    int64(req.EmailId),
		SecretCode: req.SecretCode,
	}

	txResult, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to verify email"))
		return
	}

	rsp := &VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	ctx.JSON(http.StatusCreated, rsp)
}

// @BasePath /api/v1

// Login User godoc
// @Summary perform a new user login
// @Schemes
// @Description returns a new PASETO token and the logged user info
// @Tags users
// @Accept json
// @Produce json
// @Param req body LoginUserRequest true "LoginUserRequest"
// @Success 200 {object} LoginUserResponse
// @Failure 404 "no rows in resultset"
// @Failure 400 "user not verified"
// @Failure 401 "invalid credentials"
// @Router /users/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	session := sessions.Default(ctx)

	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse("user not found"))
		return
	}

	if !user.IsEmailVerified {
		ctx.JSON(http.StatusBadRequest, errorResponse("user not verified"))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse("invalid credentials"))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(req.Username, int(user.Role.Int32), server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	session.Set(userkey, req.Username)
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	log.Info().Str("user", user.Username).Msg("successfully logged in")

	loginResponse := LoginUserResponse{
		AccessToken: accessToken,
		User:        CreateUserResponse(user),
	}
	ctx.JSON(http.StatusOK, loginResponse)
}

// @BasePath /api/v1

// Get user info godoc
// @Summary get the user info
// @Schemes
// @Description returns the user info
// @Tags users
// @Param username path string  true  "Username"
// @Param authorization header string  true  "Authorization"
// @Produce json
// @Success 200 {object} UserResponse
// @Router /user/{username} [get]
func (server *Server) getUser(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != user.Username {
		err := errors.New("user doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, CreateUserResponse(user))
}

// @BasePath /api/v1

// Logout User godoc
// @Summary perform a user logout
// @Schemes
// @Description delete the user session
// @Tags users
// @Produce json
// @Success 200 "Successfully logged out"
// @Router /users/logout [post]
func (server *Server) logoutUser(c *gin.Context) {
	u, r := LoggedUsernameAndRole(c, server.tokenMaker)
	log.Info().Int("role", r).Str("user", u).Msg("successfully logged out")
	session := sessions.Default(c)
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// @BasePath /api/v1

// List the users godoc
// @Summary get the user list paginated
// @Schemes
// @Description get paginated user list
// @Tags users
// @Param authorization header string  true  "Authorization"
// @Accept json
// @Produce json
// @Param req query PaginationRequest true "PaginationRequest"
// @Success 200 {array} UserResponse
// @Router /admin/users [get]
func (server *Server) listUsers(ctx *gin.Context) {
	var req PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}
	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}
	var lUsers []UserResponse
	for _, u := range users {
		lUsers = append(lUsers, CreateUserResponse(u))
	}
	ctx.JSON(http.StatusOK, lUsers)
}

// @BasePath /api/v1

// Get the users count godoc
// @Summary get the user count for pagination
// @Schemes
// @Description get user count
// @Tags users
// @Param authorization header string  true  "Authorization"
// @Produce json
// @Success 200 {object} CountUsersResponse
// @Router /admin/users/count [get]
func (server *Server) countUsers(ctx *gin.Context) {
	countUsers, err := server.store.CountUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}
	countResponse := CountUsersResponse{
		Count: int(countUsers),
	}
	ctx.JSON(http.StatusOK, countResponse)
}
