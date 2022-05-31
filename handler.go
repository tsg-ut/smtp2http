package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/alash3al/go-smtpsrv"
	"github.com/go-resty/resty/v2"
	"github.com/zaccone/spf"
)

func StrictValidation(req *smtpsrv.Request) error {
	ipstr, _, _ := net.SplitHostPort(req.RemoteAddr)
	ip := net.ParseIP(ipstr)

	// trust mails from gmail.com
	// https://datatracker.ietf.org/doc/html/rfc7208#section-4
	spfgmail, _, _ := spf.CheckHost(ip, "gmail.com", "gmail.com")
	nonspam := (spfgmail == spf.Pass) || ip.IsLoopback() || ip.IsPrivate()

	if !nonspam {
		// should say "you are spammer" to spammer.
		return errors.New("Your host isn't configured correctly or you are a spammer -_-")
	}
	return nil
}

func handler(req *smtpsrv.Request) error {
	log.Printf(`[handler] from:"%s" to:"%s" remote_addr:"%s"`+"\n", req.From, req.To, req.RemoteAddr)

	// validate the from data
	if *flagStrictValidation {
		if err := StrictValidation(req); err != nil {
			log.Println(`[handler] spam detected`)
			return err
		}
	}

	msg, err := parsemail.Parse(req.Message)
	if err != nil {
		log.Println(`[handler] fail to read, maybe because of huge size`)
		return nil
	}

	params := map[string]interface{}{
		"id":      msg.Header.Get("Message-ID"),
		"subject": msg.Subject,
		"body": map[string]interface{}{
			"text": string(msg.TextBody),
			"html": string(msg.HTMLBody),
		},
		"addresses": map[string]interface{}{
			"mailfrom": req.From,
			"from":     strings.Join(extractEmails(msg.From), ","),
			"to":       strings.Join(extractEmails(msg.To), ","),
			"cc":       strings.Join(extractEmails(msg.Cc), ","),
			"bcc":      strings.Join(extractEmails(msg.Bcc), ","),
		},
		"file_count": strconv.Itoa(len(msg.Attachments)),
	}

	body, err := json.Marshal(params)
	if err != nil {
		log.Println(`[handler] internal error while encoding params: '` + err.Error() + `'`)
		return nil
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(*flagWebhook)
	if err != nil {
		log.Println(`[handler] internal error: '` + err.Error() + `'`)
		return nil
	}
	if resp.StatusCode() != 200 {
		log.Println(`[handler] backend returned error: status=` + resp.Status())
		return nil
	}

	log.Println(`[handler] successfully processed`)
	return nil
}
