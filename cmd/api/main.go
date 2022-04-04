package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gym/internal/config"
	pg "gym/internal/db/postgres"
	"gym/pkg/booking"
	"gym/pkg/class"
	"gym/pkg/constants"
	"gym/pkg/member"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "GYM API v1 - 2022-04-03")
}
func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, time.Now().UTC())
}

func main() {
	// Config loads info from env files
	conf := config.GetConfig()
	log.Print(constants.Green + "LOAD CONFIG" + constants.Reset)

	// WEB SERVER SETUP
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// SQL REPOSITORIES
	postgresConn := pg.NewPostgresConnectionPool(conf.PostgresHost)
	classRepository := class.NewRepository(postgresConn)
	memberRepository := member.NewRepository(postgresConn)
	bookingRepository := booking.NewRepository(postgresConn)

	// HTTP HANDLERS
	validator := validator.New()
	r.GET("/", handleVersion)
	r.GET("/health", handleHealthCheck)
	class.NewHandler(r, "classes", validator, classRepository)
	member.NewHandler(r, "members", validator, memberRepository)
	booking.NewHandler(r, "bookings", validator, bookingRepository)

	// SERVER SETUP
	port := conf.Port
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
