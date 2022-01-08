package api

import (
	"fmt"

	db2 "github.com/diantanjung/weterm/backend/db/sqlc"
	"github.com/diantanjung/weterm/backend/token"
	"github.com/diantanjung/weterm/backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	querier    db2.Querier
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, querier db2.Querier) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		querier:    querier,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", server.config.FeUrl)
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(server.CORSMiddleware())

	router.POST("/users/login", server.loginUser)
	// router.GET("/ws", server.WebSocket)
	router.GET("/ws2/:username", server.WebSocket2)

	router.POST("/run", server.RunCommand)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/users", server.createUser)
	authRoutes.GET("/user", server.getUser)

	authRoutes.POST("/directory", server.CreateUserDir)
	authRoutes.GET("/directory", server.GetUserDirs)
	authRoutes.DELETE("/directory", server.DeleteUserDirs)

	authRoutes.GET("/commands", server.GetCommands)
	authRoutes.POST("/command/:dir/:cmd", server.CreateCommand)
	authRoutes.GET("/command/:dir/:cmd", server.GetSourceCode)
	authRoutes.PATCH("/command/:dir/:cmd", server.UpdateSourceCode)
	authRoutes.DELETE("/command/:dir/:cmd", server.DeleteCommand)
	authRoutes.GET("/exe/:dir/:cmd", server.ExeCommand)

	authRoutes.POST("/open", server.GetFileContent)
	authRoutes.PATCH("/open", server.UpdateFileContent)
	//authRoutes.GET("/terminal/:dir/:cmd/:exe", server.Terminal)

	authRoutes.POST("/opendir", server.GetDirContent)

	authRoutes.POST("/opendirfile", server.GetDirFileContent)

	authRoutes.POST("/terminal", server.Terminal)
	authRoutes.POST("/autocomplete", server.AutoComplete)

	//authRoutes.POST("/ws", server.WebSocket)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
