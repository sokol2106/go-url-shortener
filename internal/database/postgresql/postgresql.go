package postgresql

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Postgresql struct {
	db     *sql.DB
	config string
}

func New(cnf string) *Postgresql {
	var pstg = Postgresql{}
	pstg.config = cnf

	return &pstg
}

func (pstg *Postgresql) Connect() error {
	var err error
	pstg.db, err = sql.Open("pgx", pstg.config)
	return err
}

func (pstg *Postgresql) Disconnect() error {
	if pstg.db != nil {
		return pstg.db.Close()
	}
	return nil
}

func (pstg *Postgresql) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return pstg.db.PingContext(ctx)
}
