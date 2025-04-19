package config

import (
	"database/sql"
	"log/slog"
)

func Initialise() (*Configuration, *slog.Logger, *sql.DB, error) {
	conf := loadConf()
	logger := setupLog(conf.Log)
	db, err := setupDb(conf.Db)

	return conf, logger, db, err
}
