package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"strings"
	"time"
)

type Postgresql struct {
	cnf    map[string]string
	db     *pgx.Conn
	config string
}

func New(cnf string) *Postgresql {
	var pstg = Postgresql{}
	pstg.cnf = make(map[string]string)
	pstg.config = cnf

	if cnf == "" {
		return &pstg
	}

	params := strings.Fields(cnf)
	for _, value := range params {
		res := strings.Split(value, "=")
		if len(res) != 2 {
			log.Printf("Erron new postgresql connection parameter: %s", cnf)
			return &pstg
		}
		pstg.cnf[res[0]] = res[1]
	}

	return &pstg
}

func (pstg *Postgresql) Connect() error {
	var err error
	//params := new(bytes.Buffer)
	//for key, value := range pstg.cnf {
	//	fmt.Fprintf(params, "%s=%s ", key, value)
	//}
	//pstg.db, err = sql.Open("pgx", pstg.config)
	pstg.db, err = pgx.Connect(context.Background(), pstg.config)
	return err
}

func (pstg *Postgresql) Disconnect() error {
	if pstg.db != nil {
		return pstg.db.Close(context.Background())
	}
	return nil
}

func (pstg *Postgresql) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return pstg.db.Ping(ctx)
}
