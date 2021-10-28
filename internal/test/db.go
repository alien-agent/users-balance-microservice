package test

import (
	"path"
	"runtime"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq" // initialize posgresql for test
	"users-balance-microservice/internal/config"
	"users-balance-microservice/pkg/dbcontext"
	"users-balance-microservice/pkg/log"
)

var db *dbcontext.DB

// DB returns the database connection for testing purpose.
func DB(t *testing.T) *dbcontext.DB {
	if db != nil {
		return db
	}
	logger, _ := log.NewForTest()
	dir := getSourcePath()
	cfg, err := config.Load(dir+"/../../config/test.yml", logger)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dbc, err := dbx.MustOpen("postgres", cfg.StorageUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dbc.LogFunc = logger.Infof
	db = dbcontext.New(dbc)
	return db
}

// ResetTables truncates all data in the specified tables.
func ResetTables(t *testing.T, db *dbcontext.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.DB().TruncateTable(table).Execute()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

// getSourcePath returns the directory containing the source code that is calling this function.
func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}