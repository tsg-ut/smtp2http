package main

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/zaccone/spf"
	"go.uber.org/zap"
)

const (
	MessageSpammer = "Your host isn't configured correctly or you are a spammer -_-"
)

type Handler struct {
	Validation bool
	Webhook    string
}

func StrictValidation(c Context) error {
	ipstr, _, _ := net.SplitHostPort(c.RemoteAddr().String())
	ip := net.ParseIP(ipstr)

	// trust emails from gmail.com
	// https://datatracker.ietf.org/doc/html/rfc7208#section-4
	spfgmail, _, _ := spf.CheckHost(ip, "gmail.com", "gmail.com")
	nonspam := (spfgmail == spf.Pass) || ip.IsLoopback() || ip.IsPrivate()

	if !nonspam {
		// should say "you are spammer" to spammer.
		return errors.New(MessageSpammer)
	}
	return nil
}

func (h *Handler) handle(c Context) error {
	log := zap.S().With(
		"from", c.From(),
		"to", c.To(),
		"remote_addr", c.RemoteAddr(),
	)
	log.Infow("receive message")

	// validate the from data
	if h.Validation {
		if err := StrictValidation(c); err != nil {
			log.Infow("spam detected")
			return err
		}
	}

	msg, err := c.Parse()
	if err != nil {
		log.Errorw("fail to read", "error", err)
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
			"mailfrom": c.From().Address,
			"from":     strings.Join(extractEmails(msg.From), ","),
			"to":       strings.Join(extractEmails(msg.To), ","),
			"cc":       strings.Join(extractEmails(msg.Cc), ","),
			"bcc":      strings.Join(extractEmails(msg.Bcc), ","),
		},
		"file_count": strconv.Itoa(len(msg.Attachments)),
	}

	body, err := json.Marshal(params)
	if err != nil {
		log.Errorw("internal error while encoding params", "error", err)
		return nil
	}

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(h.Webhook)
	if err != nil {
		log.Errorw("internal error", "error", err)
		return nil
	}
	if resp.StatusCode() != 200 {
		log.Warnw("backend returned error", "status", resp.Status())
		return nil
	}

	log.Infow("successfully processed")
	return nil
}
