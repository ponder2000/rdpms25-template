package pkg

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ponder2000/rdpms25-template/pkg/config"
	"github.com/ponder2000/rdpms25-template/pkg/util"
)

func Start(version, buildTime string) {
	conf, logger, dbConn, err := config.Initialise()
	logger.Info("starting application", "version", version, "build", buildTime)
	if err != nil {
		logger.Error("unable to start application", "err", err)
		os.Exit(1)
	}

	logger.Info("configuration loaded ")

	util.Periodic(
		time.Minute,
		util.LogRuntimeStatsBasic,
		func() { slog.Info("db status", "stats", dbConn.Stats()) },
	)

	app := gin.New()

	if e := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), app); e != nil {
		slog.Error("unable to start server", "error", e.Error(), "port", conf.Port)
	}
}
