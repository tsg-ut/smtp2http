package main

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/zaccone/spf"
	"go.uber.org/zap"
)

var (
	ErrSpammer = errors.New("Your host isn't configured correctly or you are a spammer -_-")
)

type Handler struct {
	Validation bool
	Webhook    string
}

// Validate returns true if the email is worth receiving
func Validate(c Context) bool {
	ipstr, _, _ := net.SplitHostPort(c.RemoteAddr().String())
	ip := net.ParseIP(ipstr)

	// trust emails from gmail.com
	// https://datatracker.ietf.org/doc/html/rfc7208#section-4
	spfgmail, _, _ := spf.CheckHost(ip, "gmail.com", "gmail.com")
	nonspam := (spfgmail == spf.Pass) || ip.IsLoopback() || ip.IsPrivate()

	return nonspam
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
		if !Validate(c) {
			log.Infow("spam detected")
			return ErrSpammer
		}
	}

	msg, err := c.Parse()
	if err != nil {
		log.Errorw("fail to read", "error", err)
		return nil
	}

	params := map[string]interface{}{
		"id":      msg.GetHeader("Message-ID"),
		"subject": msg.GetHeader("Subject"),
		"body": map[string]interface{}{
			"text": string(msg.Text),
			"html": string(msg.HTML),
		},
		"addresses": map[string]interface{}{
			"mailfrom": c.From().Address,
			"from":     msg.GetHeader("From"),
			"to":       msg.GetHeader("To"),
			"cc":       msg.GetHeader("Cc"),
			"bcc":      msg.GetHeader("Bcc"),
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
