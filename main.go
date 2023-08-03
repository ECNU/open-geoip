package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ECNU/open-geoip/controller"
	"github.com/ECNU/open-geoip/cron"
	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/models"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	csvFile := flag.String("csv", "", "internal file")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *csvFile != "" {
		fmt.Println("import csv file", *csvFile)
		g.InitInternalDB(*csvFile)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	// init redis pool only ratelimit enbaled
	if g.Config().RateLimit.Enabled {
		g.InitRedisConnPool()
	}

	if g.Config().SSO.Enabled {
		if g.Config().Oauth2.Enabled {
			g.InitOauth2()
		}
	}

	srv := controller.InitGin(g.Config().Http.Listen)
	g.InitLog(g.Config().Logger)

	err := models.InitReader()
	if err != nil {
		log.Fatalf("load geo db failed, %v", err)
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	if g.Config().AutoDownload.Enabled {
		go cron.SyncMaxmindDatabase()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: %v", err)
	}
	log.Println("Server exit")
}
