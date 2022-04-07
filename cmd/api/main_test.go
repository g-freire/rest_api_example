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
	"strings"
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

	// uncommment to re-create the db state after the tests
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
	sql := "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	if _, err := db.Exec(context.Background(), sql); err != nil {
		log.Fatal(err)
	}
}

// seeds the database with test data - the creation_data is hardcoded to make sure the test is deterministic
// we could also have an dedicated QA db instance for this purpose
func seedDB(db *pgxpool.Pool) {
	sql := `
			-- CLASS
			INSERT INTO class(name, start_date, end_date,capacity, creation_time) VALUES  ('pilates','2021-12-01','2021-12-20', 20 , '2023-01-01 00:00:00.000000');
			INSERT INTO class(name, start_date, end_date,capacity, creation_time) VALUES  ('crossfit','2021-12-01','2021-12-30', 40, '2023-01-01 00:00:00.000000');
			INSERT INTO class(name, start_date, end_date,capacity, creation_time) VALUES  ('jiu-jitsu','2022-04-04','2022-04-20', 30,'2023-01-01 00:00:00.000000');
			
			-- MEMBER
			INSERT INTO member(name, creation_time) VALUES  ('Alice',  '2023-01-01 00:00:00.000000');
			INSERT INTO member(name, creation_time) VALUES  ('Bob',    '2023-01-01 00:00:00.000000');
			INSERT INTO member(name, creation_time) VALUES  ('Charlie','2023-01-01 00:00:00.000000');
			
			-- BOOKING
			INSERT INTO booking(class_id, member_id, date, creation_time) VALUES (1,1,'2021-12-01','2023-01-01 00:00:00.000000');
			INSERT INTO booking(class_id, member_id, date, creation_time) VALUES (2,1,'2021-12-02','2023-01-01 00:00:00.000000');
			INSERT INTO booking(class_id, member_id, date, creation_time) VALUES (1,3,'2021-12-01','2023-01-01 00:00:00.000000');
			INSERT INTO booking(class_id, member_id, date, creation_time) VALUES (1,3,'2021-12-02','2023-01-01 00:00:00.000000');
			INSERT INTO booking(class_id, member_id, date, creation_time) VALUES (1,3,'2021-12-03','2023-01-01 00:00:00.000000');
		;`
	if _, err := db.Exec(context.Background(), sql); err != nil {
		log.Fatal(err)
	}
}

func TestHttpEndpoints(t *testing.T) {
	log.Print(constants.Yellow + "------------------------------------------------------------" + constants.Reset)
	log.Print(constants.Yellow + "STARTING API ENDPOINTS TESTS" + constants.Reset)
	log.Print(constants.Yellow + "------------------------------------------------------------" + constants.Reset)

	r := setup()

	// GET /version
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	////////////////////////////////////////////////////////////////////////////////////////
	// v1/classes
	////////////////////////////////////////////////////////////////////////////////////////
	// GET /v1/classes
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes", nil)
	r.ServeHTTP(w, req)
	responseString := "[{\"id\":1,\"creation_time\":\"2023-01-01T00:00:00Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}," +
		"{\"id\":2,\"creation_time\":\"2023-01-01T00:00:00Z\",\"name\":\"crossfit\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-30T00:00:00Z\",\"capacity\":40}," +
		"{\"id\":3,\"creation_time\":\"2023-01-01T00:00:00Z\",\"name\":\"jiu-jitsu\",\"start_date\":\"2022-04-04T00:00:00Z\",\"end_date\":\"2022-04-20T00:00:00Z\",\"capacity\":30}]"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// Get /v1/classes?name={name}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes?name=pilates", nil)
	r.ServeHTTP(w, req)
	responseString = "[{\"id\":1,\"creation_time\":\"2023-01-01T00:00:00Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}]"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// GET /v1/classes/{:id}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/2", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"id\":2,\"creation_time\":\"2023-01-01T00:00:00Z\",\"name\":\"crossfit\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-30T00:00:00Z\",\"capacity\":40}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())
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

	// GET /v1/classes?date?star={timestamp}&end={timestamp}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/date?start=2021-12-22&end=2021-12-20", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":404,\"Type\":\"Resource Not Found\",\"Message\":[\"Invalid timestamps: start \\u003c end\"]}"
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/date?start=2021-12-01&end=2021-12-20", nil)
	r.ServeHTTP(w, req)
	responseString = "[{\"id\":1,\"creation_time\":\"2023-01-01T00:00:00Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}]"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// GET /v1/classes/count
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/classes/count", nil)
	r.ServeHTTP(w, req)
	responseString = "3"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// POST /v1/classes/
	w = httptest.NewRecorder()
	body := strings.NewReader(`{
				"name": "pilates2",
				"start_date": "2022-12-19T00:00:00.000000Z",
				"end_date": "2022-12-20T00:00:00.000000Z",
				"capacity": 30
	}`)
	req, _ = http.NewRequest("POST", "/v1/classes", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":4,\"Message\":\"Created Class successfully\",\"Status\":201}"
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	w = httptest.NewRecorder()
	body = strings.NewReader(`{
				"name": "pilates2",
				"start_date": "2020-12-19T00:00:00.000000Z",
				"end_date": "2022-12-20T00:00:00.000000Z",
				"capacity": 30
	}`)
	req, _ = http.NewRequest("POST", "/v1/classes", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":400,\"Type\":\"Database Operation Error\",\"Message\":[\"Invalid timestamps: start should be \\u003e current time\"]}"
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// UPDATE /v1/classes/{:id}
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
        "name": "soccer",
        "start_date": "2022-12-18T00:00:00.000000Z",
        "end_date": "2022-12-20T00:00:00.000000Z",
        "capacity": 102
	}`)
	req, _ = http.NewRequest("PUT", "/v1/classes/4", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":\"4\",\"Message\":\"Updated Class with successfully\",\"Status\":200}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// DELETE /v1/classes/{:id}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/v1/classes/4", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":\"4\",\"Message\":\"Deleted Class successfully\",\"Status\":200}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	////////////////////////////////////////////////////////////////////////////////////////
	// v1/members
	////////////////////////////////////////////////////////////////////////////////////////
	// GET /v1/members/
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/members", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String())

	// GET /v1/members/{:id}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/members/1", nil)
	r.ServeHTTP(w, req)
	//responseString := "{\"id\":1,\"creation_time\":\"2022-04-05T09:08:31.95438Z\",\"name\":\"pilates\",\"start_date\":\"2021-12-01T00:00:00Z\",\"end_date\":\"2021-12-20T00:00:00Z\",\"capacity\":20}"
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, responseString, w.Body.String())

	// GET /v1/members/count
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/members/count", nil)
	r.ServeHTTP(w, req)
	responseString = "3"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// POST /v1/members/
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
        "name":"Alice"
	}`)
	req, _ = http.NewRequest("POST", "/v1/members", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":4,\"Message\":\"Created Member successfully\",\"Status\":201}"
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// UPDATE /v1/members/{:id}
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
	   "name":"Bob"
	}`)
	req, _ = http.NewRequest("PUT", "/v1/members/1", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":\"1\",\"Message\":\"Updated Member successfully\",\"Status\":200}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// DELETE /v1/bookings/{:id}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/v1/members/3", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":\"3\",\"Message\":\"Deleted Member successfully\",\"Status\":200}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	////////////////////////////////////////////////////////////////////////////////////////
	// v1/booking
	////////////////////////////////////////////////////////////////////////////////////////
	// GET /v1/bookings/
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/bookings", nil)
	r.ServeHTTP(w, req)
	responseString = "[{\"id\":1,\"class_id\":1,\"member_id\":1,\"date\":\"2021-12-01T00:00:00Z\",\"creation_time\":\"2023-01-01T00:00:00Z\"}," +
		"{\"id\":2,\"class_id\":2,\"member_id\":1,\"date\":\"2021-12-02T00:00:00Z\",\"creation_time\":\"2023-01-01T00:00:00Z\"}," +
		"{\"id\":3,\"class_id\":1,\"member_id\":3,\"date\":\"2021-12-01T00:00:00Z\",\"creation_time\":\"2023-01-01T00:00:00Z\"}," +
		"{\"id\":4,\"class_id\":1,\"member_id\":3,\"date\":\"2021-12-02T00:00:00Z\",\"creation_time\":\"2023-01-01T00:00:00Z\"}," +
		"{\"id\":5,\"class_id\":1,\"member_id\":3,\"date\":\"2021-12-03T00:00:00Z\",\"creation_time\":\"2023-01-01T00:00:00Z\"}]"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// GET /v1/bookings/{:id}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/bookings/1", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"id\":1,\"class_id\":1,\"member_id\":1,\"date\":\"2021-12-01T00:00:00Z\",\"creation_time\":\"2023-01-01T00:00:00Z\"}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// GET /v1/bookings/count
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/bookings/count", nil)
	r.ServeHTTP(w, req)
	responseString = "5"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// POST /v1/booking/
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
			"class_id": 1,
			"member_id": 1,
			"date": "2022-12-01T00:00:00.000000Z"
	}`)
	req, _ = http.NewRequest("POST", "/v1/bookings", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":6,\"Message\":\"Created Booking successfully\",\"Status\":201}"
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	w = httptest.NewRecorder()
	body = strings.NewReader(`{
			"class_id": 1,
			"member_id": 1,
			"date": "1980-12-01T00:00:00.000000Z"
	}`)
	req, _ = http.NewRequest("POST", "/v1/bookings", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":400,\"Type\":\"Invalid Request Body\",\"Message\":[\"Invalid timestamps: start \\u003c end\"]}"
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// UPDATE /v1/bookings/{:id}
	w = httptest.NewRecorder()
	body = strings.NewReader(`{
	   "class_id": 1,
	   "member_id": 2,
	   "date": "2050-12-01T00:00:00.000000Z"
	}`)
	req, _ = http.NewRequest("PUT", "/v1/bookings/1", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":\"1\",\"Message\":\"Updated Booking with successfully\",\"Status\":200}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	w = httptest.NewRecorder()
	body = strings.NewReader(`{
	   "class_id": 1,
	   "member_id": 22,
	   "date": "2050-12-01T00:00:00.000000Z"
	}`)
	req, _ = http.NewRequest("PUT", "/v1/bookings/1", body)
	r.ServeHTTP(w, req)
	responseString = "{\"Status\":400,\"Type\":\"Database Operation Error\",\"Message\":[\"ERROR: insert or update on table \\\"booking\\\" violates foreign key constraint \\\"booking_member_id_fkey\\\" (SQLSTATE 23503)\"]}"
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, responseString, w.Body.String())

	// DELETE /v1/bookings/{:id}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/v1/bookings/1", nil)
	r.ServeHTTP(w, req)
	responseString = "{\"Id\":\"1\",\"Message\":\"Deleted Booking successfully\",\"Status\":200}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, responseString, w.Body.String())
}
