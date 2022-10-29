package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-smtp"
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

	s := smtp.NewServer(&backend{handler: handler.handle})
	s.Addr = *flagListenAddr
	s.Domain = *flagServerName
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = int(*flagMaxMessageSize)
	s.AllowInsecureAuth = true
	s.AuthDisabled = true
	s.EnableSMTPUTF8 = false

	log.Infow("start",
		"addr", s.Addr,
		"maxmessage", s.MaxMessageBytes,
		"banner", s.Domain,
		"webhook", handler.Webhook,
		"validation", handler.Validation,
	)

	log.Info(s.ListenAndServe())
}
