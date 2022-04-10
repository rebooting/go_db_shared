package go_db_shared

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v4"

	"github.com/rebooting/go_utils"
)

type DB struct {
	DB *pgx.Conn
}

func New(connStr string) (DB, error) {

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		go_utils.Log(go_utils.GetCaller(), err.Error())
		return DB{}, err
	}

	return DB{DB: conn}, nil

}

func Migrate(sqlpath string, connstr string) (resErr error) {

	m, err := migrate.New(sqlpath, connstr)
	if err != nil {
		go_utils.Log(go_utils.GetCaller(), err.Error())
		resErr = err
	}
	if err := m.Up(); err != nil {
		go_utils.Log(go_utils.GetCaller(), err.Error())
		resErr = err
	}
	return
}
