package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"gym/internal/config"
	"gym/internal/constants"
	pg "gym/internal/db/postgres"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//TODO Add more tests

const (
	migrationsTestRootFolder = "file://../../migration"
)

func TestMain(m *testing.M) {
	// Prepare database and seed for the tests
	conf := config.GetConfig()
	log.Print(constants.Green + "LOAD CONFIG" + constants.Reset)
	postgresConn := pg.NewPostgresConnectionPool(conf.PostgresHost)
	resetDB(postgresConn)
	err := pg.Migrate(conf.PostgresHost, migrationsTestRootFolder, "up", 0)
	if err != nil {
		log.Fatal(err)
	}
	ensureTableExists(postgresConn)
	seedDB(postgresConn)

	// Start unit tests
	exitVal := m.Run()
	// Clean Seed after it
	//resetDB(postgresConn)
	os.Exit(exitVal)
}

func ensureTableExists(db *pgxpool.Pool) {
	sql := "SELECT COUNT(*) FROM class"
	if _, err := db.Exec(context.Background(), sql); err != nil {
		log.Fatal(err)
	}
}

func resetDB(db *pgxpool.Pool) {
	sql := "DROP SCHEMA public CASCADE; CREATE SCHEMA public;\n"
	if _, err := db.Exec(context.Background(), sql); err != nil {
		log.Fatal(err)
	}
}

func seedDB(db *pgxpool.Pool) {
	sql := `
			-- CLASS
			INSERT INTO class(name, start_date, end_date,capacity) VALUES  ('pilates','2021-12-01','2021-12-20', 20);
			INSERT INTO class(name, start_date, end_date,capacity) VALUES  ('crossfit','2021-12-01','2021-12-30', 40);
			INSERT INTO class(name, start_date, end_date,capacity) VALUES  ('jiu-jitsu','2022-04-04','2022-04-20', 30);
			
			-- MEMBER
			INSERT INTO member(name) VALUES  ('Alice');
			INSERT INTO member(name) VALUES  ('Bob');
			INSERT INTO member(name) VALUES  ('Charlie');
			
			-- BOOKING
			INSERT INTO booking(class_id, member_id, date) VALUES (1,1,'2021-12-01');
			INSERT INTO booking(class_id, member_id, date) VALUES (2,1,'2021-12-02');
			INSERT INTO booking(class_id, member_id, date) VALUES (1,3,'2021-12-01');
			INSERT INTO booking(class_id, member_id, date) VALUES (1,3,'2021-12-02');
			INSERT INTO booking(class_id, member_id, date) VALUES (1,3,'2021-12-03');
		;`
	if _, err := db.Exec(context.Background(), sql); err != nil {
		log.Fatal(err)
	}
}

func TestGETRoutes(t *testing.T) {
	r := setup()
	// version
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	////////////////////////////////////////////////////////////////////////////////////////
	// v1/classes - GET
	////////////////////////////////////////////////////////////////////////////////////////
	// SUCCESS 200
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String()) // will break because of the dynamic creation time

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/1", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String()) // will break because of the dynamic creation time

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/date?start=2021-12-02&end=2021-12-20", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/count", nil)
	r.ServeHTTP(w, req)
	responseString := "3"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// ERRORS 404
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/1123", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":404,\"Type\":\"Resource Not Found\",\"Message\":[\"no rows in result set\"]}"
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/b", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":404,\"Type\":\"Resource Not Found\",\"Message\":[\"URL 'id' Parameter must be a number\"]}"
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/date?start=2021-12-22&end=2021-12-20", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":404,\"Type\":\"Resource Not Found\",\"Message\":[\"Invalid timestamps: start \\u003c end\"]}"
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	////////////////////////////////////////////////////////////////////////////////////////
	// v1/members - GET
	////////////////////////////////////////////////////////////////////////////////////////
	// SUCCESS 200
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/members", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String()) // will break because of the dynamic creation time

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/members/1", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String()) // will break because of the dynamic creation time

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/members/count", nil)
	r.ServeHTTP(w, req)
	responseString = "3"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	////////////////////////////////////////////////////////////////////////////////////////
	// v1/booking - GET
	////////////////////////////////////////////////////////////////////////////////////////
	// SUCCESS 200
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/bookings", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String()) // will break because of the dynamic creation time

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/bookings/1", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String()) // will break because of the dynamic creation time

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/bookings/count", nil)
	r.ServeHTTP(w, req)
	responseString = "5"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())
}
