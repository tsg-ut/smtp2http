package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/alash3al/go-smtpsrv/v3"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("ver:%s rev:%s\n", Version, Revision)
		return
	}

	{
		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync()
		zap.ReplaceGlobals(logger)
	}

	log := zap.S()

	handler := Handler{
		Validation: *flagStrictValidation,
		Webhook:    *flagWebhook,
	}

	cfg := smtpsrv.ServerConfig{
		ListenAddr:      *flagListenAddr,
		BannerDomain:    *flagServerName,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		MaxMessageBytes: int(*flagMaxMessageSize),
		Handler:         func(c *smtpsrv.Context) error { return handler.handle(c) },
	}

	log.Infow("start",
		"addr", cfg.ListenAddr,
		"maxmessage", cfg.MaxMessageBytes,
		"banner", cfg.BannerDomain,
		"webhook", handler.Webhook,
		"validation", handler.Validation,
	)
	log.Sync()

	log.Info(smtpsrv.ListenAndServe(&cfg))
}
