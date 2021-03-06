package sql

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gopkg.in/DataDog/dd-trace-go.v1/contrib/internal/sqltest"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

// tableName holds the SQL table that these tests will be run against. It must be unique cross-repo.
const tableName = "testsql"

func TestMain(m *testing.M) {
	_, ok := os.LookupEnv("INTEGRATION")
	if !ok {
		fmt.Println("--- SKIP: to enable integration test, set the INTEGRATION environment variable")
		os.Exit(0)
	}
	defer sqltest.Prepare(tableName)()
	os.Exit(m.Run())
}

func TestMySQL(t *testing.T) {
	Register("mysql", &mysql.MySQLDriver{})
	db, err := Open("mysql", "test:test@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	testConfig := &sqltest.Config{
		DB:         db,
		DriverName: "mysql",
		TableName:  tableName,
		ExpectName: "mysql.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName:     "mysql.db",
			ext.SpanType:        ext.SpanTypeSQL,
			ext.TargetHost:      "127.0.0.1",
			ext.TargetPort:      "3306",
			ext.DBUser:          "test",
			ext.DBName:          "test",
			ext.EventSampleRate: nil,
		},
	}
	sqltest.RunAll(t, testConfig)
}

func TestPostgres(t *testing.T) {
	Register("postgres", &pq.Driver{}, WithServiceName("postgres-test"), WithAnalyticsRate(0.2))
	db, err := Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	testConfig := &sqltest.Config{
		DB:         db,
		DriverName: "postgres",
		TableName:  tableName,
		ExpectName: "postgres.query",
		ExpectTags: map[string]interface{}{
			ext.ServiceName:     "postgres-test",
			ext.SpanType:        ext.SpanTypeSQL,
			ext.TargetHost:      "127.0.0.1",
			ext.TargetPort:      "5432",
			ext.DBUser:          "postgres",
			ext.DBName:          "postgres",
			ext.EventSampleRate: 0.2,
		},
	}
	sqltest.RunAll(t, testConfig)
}
