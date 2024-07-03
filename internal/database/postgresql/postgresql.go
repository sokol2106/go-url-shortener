package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"strings"
	"time"
)

type Postgresql struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string
	db       *sql.DB
}

func New(cnf string) *Postgresql {
	var pstg = Postgresql{}
	pstg.host = "disable"
	pstg.port = "5432" // standart port postgresql
	params := strings.Fields(cnf)
	for _, value := range params {
		res := strings.Split(value, "=")
		switch res[0] {
		case "host":
			pstg.host = res[1]
		case "port":
			pstg.port = res[1]
		case "user":
			pstg.user = res[1]
		case "password":
			pstg.password = res[1]
		case "dbname":
			pstg.dbname = res[1]
		case "sslmode":
			pstg.sslmode = res[1]
		}
	}

	return &pstg
}

func (pstg *Postgresql) Connect() error {
	var err error
	ps := fmt.Sprintf("host=%s posrt=%S user=%s password=%s dbname=%s sslmode=%s",
		pstg.host, pstg.port, pstg.user, pstg.password, pstg.dbname, pstg.sslmode)
	pstg.db, err = sql.Open("pgx", ps)
	return err
}

func (pstg *Postgresql) Disconnect() error {
	return pstg.db.Close()
}

func (pstg *Postgresql) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return pstg.db.PingContext(ctx)
}
