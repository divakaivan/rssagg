package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/divakaivan/rssagg/internal/database"
	_ "github.com/lib/pq"
)

var testQueries *database.Queries
var testDB *sql.DB

// entry point for the test suite; called only once, before all tests; does setup and teardown
func TestMain(m *testing.M) {
	var err error
	fmt.Println("Initializing test suite")

	pgContainer, err := NewPGContainer()
	if err != nil {
		log.Fatal("cannot create a new container:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal("cannot terminate the container:", err)
		}
	}()

	testDB = pgContainer.DB()
	testQueries = database.New(testDB)

	_, err = testDB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
