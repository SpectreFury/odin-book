package router

import (
	"os"

	"github.com/SpectreFury/odin-book/backend/src/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	return r
}

func RegisterRoutes(r *gin.Engine, conn *pgx.Conn) {
	h := handlers.Handlers{
		Conn: conn,
	}

	r.GET("/", h.IndexHandler)
	r.POST("/login", h.LoginHandler)
}

func Run(r *gin.Engine) error {
	port := os.Getenv("PORT")

	err := r.Run(":" + port)
	if err != nil {
		return err
	}

	return nil
}
