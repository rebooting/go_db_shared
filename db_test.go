package go_db_shared

import (
	"context"
	"fmt"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/cockroachdb"
)

func TestMigrate(t *testing.T) {
	p := cockroachdb.Preset(
		cockroachdb.WithDatabase("gnomock"),
		cockroachdb.WithVersion("v21.2.8"),
	)
	container, err := gnomock.Start(p)
	if err != nil {
		t.Error(err)
	}

	connStr := fmt.Sprintf("cockroach://root@%s:%d/gnomock?sslmode=disable", container.Host, container.DefaultPort())

	err = Migrate("file://test_migration", connStr)
	if err != nil {
		t.Error(err)
	}
	connStr = fmt.Sprintf("postgres://root@%s:%d/gnomock?sslmode=disable", container.Host, container.DefaultPort())
	dbo, err := New(connStr)
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()

	tag, err := dbo.DB.Exec(ctx, `INSERT INTO example.sample_table VALUES (1, 'whatever')`)
	if err != nil {
		t.Error(err)
	}
	if tag.RowsAffected() != 1 {
		t.Error(tag)
	}
	var resInt int
	var resString string

	dbo.DB.QueryRow(ctx, `SELECT * FROM example.sample_table WHERE test_integer=1`).Scan(&resInt, &resString)

	if resInt != 1 || resString != "whatever" {
		t.Errorf("got %v, %v", resInt, resString)
	}

}
