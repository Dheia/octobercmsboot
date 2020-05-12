package exec

import (
	"reflect"
	"testing"
)

func TestCreateSchemaCommand(t *testing.T) {
	dbName := "example"
	userName := "root"
	password := "root"
	want := []string{"mysql", "--user=" + userName, "--password=" + password, "-e", "CREATE SCHEMA " + dbName + " DEFAULT CHARACTER SET utf8mb4"}
	got := CreateSchemaCommand(dbName, userName, password)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("not equal")
	}
}

