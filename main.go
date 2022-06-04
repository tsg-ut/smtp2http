package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/alash3al/go-smtpsrv/v3"
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Printf("ver:%s rev:%s\n", Version, Revision)
		return
	}

	cfg := smtpsrv.ServerConfig{
		ListenAddr:      *flagListenAddr,
		BannerDomain:    *flagServerName,
		ReadTimeout:     2 * time.Second,
		WriteTimeout:    2 * time.Second,
		MaxMessageBytes: int(*flagMaxMessageSize),
		Handler:         func(c *smtpsrv.Context) error { return handler(c) },
	}

	fmt.Println("start the smtp server on address: ", *flagListenAddr)
	fmt.Println("specified maximum body size: ", *flagMaxMessageSize, " bytes")
	fmt.Println("specified server name: ", *flagServerName)
	fmt.Println("specified webhook: ", *flagWebhook)
	fmt.Println("validating the incoming FROM header: ", *flagStrictValidation)

	fmt.Println(smtpsrv.ListenAndServe(&cfg))
}
