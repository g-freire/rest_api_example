package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/go-playground/validator.v9"
	"gym/internal/cache"
	"gym/internal/config"
	"gym/internal/constants"
	pg "gym/internal/db/postgres"
	"gym/pkg/booking"
	"gym/pkg/class"
	"gym/pkg/member"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var migrationsRootFolder = "file://migration"

func handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "GYM API v1 - 2022-08-03")
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, time.Now().UTC())
}

func setup() (*gin.Engine, *pgxpool.Pool, *config.Config) {
	// CONFIGURATION
	conf := config.GetConfig()
	log.Print(constants.Green + "LOAD CONFIG" + constants.Reset)

	// CACHE
	copts := &redis.UniversalOptions{
		MasterName: "",
		Addrs:      []string{conf.RedisAddress},
		Password:   conf.RedisPassword,
	}
	ch := cache.NewClient(*copts)

	// MIGRATIONS
	postgresConn := pg.NewPostgresConnectionPool(conf.PostgresHost)
	err := pg.Migrate(conf.PostgresHost, migrationsRootFolder, "up", 0)
	if err != nil {
		log.Fatal(err)
	}
	// SQL REPOSITORIES
	classRepository := class.NewRepository(postgresConn, ch)
	memberRepository := member.NewRepository(postgresConn)
	bookingRepository := booking.NewRepository(postgresConn)

	// SERVICES
	classService := class.NewService(classRepository)
	bookingService := booking.NewService(bookingRepository)

	// WEB SERVER
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// HTTP HANDLERS
	validator := validator.New()
	r.GET("/", handleVersion)
	r.GET("/health", handleHealthCheck)
	class.NewHandler(r, "classes", validator, classService, classRepository)
	member.NewHandler(r, "members", validator, memberRepository)
	booking.NewHandler(r, "bookings", validator, bookingService, bookingRepository)
	return r, postgresConn, conf
}

func main() {
	r, postgresConn, conf := setup()

	// SERVER SETUP
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + conf.Port,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print(constants.Blue, "WEB SERVER PORT: ", conf.Port, constants.Reset)

	// GRACEFULL SHUTDOWNS
	//DB
	defer func() {
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := postgresConn.Close
		if err != nil {
			log.Fatalf("DB SHUTDOWN ERROR")
		}
	}()
	//SERVER
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
