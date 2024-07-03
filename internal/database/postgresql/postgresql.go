package postgresql

import (
	"context"
	"database/sql"
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
	pstg.db, err = sql.Open("pgx", "")
	return err
}

func (pstg *Postgresql) Disconnect() error {
	return pstg.db.Close()
}

func (pstg *Postgresql) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pstg.db.PingContext(ctx)
}
