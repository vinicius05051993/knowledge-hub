package integration

import (
	"database/sql"
	"testing"

	_ "github.com/microsoft/go-mssqldb/azuread"
)

func TestAzureSQLDriverRegistered(t *testing.T) {

	found := false

	for _, driver := range sql.Drivers() {
		if driver == "azuresql" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("azuresql driver not registered")
	}
}