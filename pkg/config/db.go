package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func setupDb(dbConf *Db) (*sql.DB, error) {
	if myDb, e := sql.Open("postgres", fmt.Sprintf(`dbname=%s host=%s port=%d user=%s password=%s sslmode=disable`, dbConf.Name, dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password)); e != nil {
		slog.Error("Error in opening connection", "err", e)
		return nil, e
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if e = myDb.PingContext(ctx); e != nil {
			return nil, e
		}

		if dbConf.IdleConnection != 0 {
			myDb.SetMaxIdleConns(dbConf.IdleConnection)
		}

		if dbConf.OpenConnection != 0 {
			myDb.SetMaxOpenConns(dbConf.OpenConnection)
		}

		boil.SetDB(myDb)
		if loc, e := time.LoadLocation("Asia/Kolkata"); e != nil {
			slog.Error("unable to set timezone as Asia/Kolkata", "err", e)
		} else {
			boil.SetLocation(loc)
		}
		slog.Info("DB setup completed")
		return myDb, nil
	}
}
