package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	pg "gym/internal/db/postgres"
	class "gym/pkg/class"
	member "gym/pkg/member"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPostgresURILocal = "postgres://gym:gym@localhost:5432/gym?sslmode=disable"
	defaultPort             = "5000"
)

func handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "GYM API v1 - 2022-04-03")
}

func main() {

	// WEB SERVER SETUP
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// SQL REPOSITORIES
	postgresConn := pg.NewPostgresConnectionPool(defaultPostgresURILocal)
	classRepository := class.NewRepository(postgresConn)
	memberRepository := member.NewRepository(postgresConn)

	// HTTP HANDLERS
	validator := validator.New()
	r.GET("/", handleVersion)
	class.NewHandler(r, "classes", validator, classRepository)
	member.NewHandler(r, "members", validator, memberRepository)
	//class.NewHandler(r, "bookings", validator, classRepository)

	// SERVER SETUP
	port := os.Getenv("GIN_PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("WEB SERVER PORT: %s", port)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// GRACEFULL SHUTDOWNS
	//DB
	defer func() {
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := postgresConn.Close
		if err != nil {
			log.Fatal(err)
		}
	}()

	//SERVER
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
