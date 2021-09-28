package api

import (
	"database/sql"
	db "github.com/diantanjung/wecom/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

type createUserDirRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID int64  `json:"user_id" binding:"required,min=1"`
}

type userDirResponse struct {
	Name      string    `json:"name"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserDirResponse(dir db.Directory) userDirResponse {
	return userDirResponse{
		Name:      dir.Name,
		UserID:    dir.UserID,
		CreatedAt: dir.CreatedAt,
	}
}

func (server *Server) CreateUserDir(ctx *gin.Context) {
	var req createUserDirRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.CreateUserDirParams{
		Name:   req.Name,
		UserID: req.UserID,
	}
	dir, err := server.querier.CreateUserDir(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserDirResponse(dir)
	ctx.JSON(http.StatusOK, rsp)
}

type getUserDirsRequest struct {
	UserID int64 `json:"user_id" binding:"required,min=1"`
}

func (server *Server) GetUserDirs(ctx *gin.Context) {
	var req getUserDirsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dirs, err := server.querier.GetUserDirs(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dirs)
}

type deleteUserDirsRequest struct {
	UserID int64  `json:"user_id" binding:"required,min=1"`
	Name   string `json:"name" binding:"required"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func (server *Server) DeleteUserDirs(ctx *gin.Context) {
	var req deleteUserDirsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args := db.DeleteUserDirParams{
		UserID: req.UserID,
		Name:   req.Name,
	}
	err := server.querier.DeleteUserDir(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	message := messageResponse{Message: "Data Delete successfully"}
	ctx.JSON(http.StatusOK, message)
}
