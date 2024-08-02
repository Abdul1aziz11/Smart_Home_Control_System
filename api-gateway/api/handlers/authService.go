package handlers

import (
	"api-gateway/genproto"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body genproto.UserReq true "User registration info"
// @Success 200 {object} genproto.RespStatus
// @Failure 400 {string} string "Bad Request"
// @Router /users/register [post]
func (h *Handlers) Register(ctx *gin.Context) {
	var user genproto.UserReq
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	resp, err := h.Auth.SignUp(ctx, &user)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	ctx.JSON(200, resp)
}

// Verify godoc
// @Summary Verify user
// @Description Verify a user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body genproto.UserReq true "User verification info"
// @Success 200 {object} genproto.RespStatusUserID
// @Failure 400 {string} string "Bad Request"
// @Router /users/verify [post]
func (h *Handlers) Verify(ctx *gin.Context) {
	var user genproto.UserReq
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	resp, err := h.Auth.Verify(ctx, &user)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	ctx.JSON(200, resp)
}

// SignIn godoc
// @Summary Sign in user
// @Description Authenticate a user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body genproto.UserReq true "User credentials"
// @Success 200 {object} genproto.UserRespFull
// @Failure 400 {string} string "Bad Request"
// @Router /users/login [post]
func (h *Handlers) SignIn(ctx *gin.Context) {
	var user genproto.UserReq
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	resp, err := h.Auth.SignIn(ctx, &user)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	ctx.JSON(200, resp)
}

// Profile godoc
// @Summary Get user profile
// @Description Get the profile of an authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param user body genproto.UserReq true "User info"
// @Success 200 {object} genproto.UserProfileResp
// @Failure 400 {string} string "Bad Request"
// @Router /profile [get]
func (h *Handlers) Profile(ctx *gin.Context) {
	var user genproto.UserReq
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	resp, err := h.Auth.Profile(ctx, &user)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	ctx.JSON(200, resp)
}
