package main_test

import (
	"testing"

	"github.com/alash3al/go-smtpsrv"

	main "github.com/tsg-ut/smtp2http"
)

func TestStrictValidation(t *testing.T) {
	tests := []struct {
		name string
		req  *smtpsrv.Request
		err  bool // err should be non-nil or not
	}{
		{
			name: "hoge",
			req: &smtpsrv.Request{
				From:       "hoge@example.com",
				RemoteAddr: "1.1.1.1:1111",
			},
			err: true,
		},
		{
			name: "valid",
			req: &smtpsrv.Request{
				From:       "hoge@gmail.com",
				RemoteAddr: "35.190.247.1:1111",
			},
			err: false,
		},
		{
			name: "lo",
			req: &smtpsrv.Request{
				From:       "hoge@localhost.localdomain",
				RemoteAddr: "127.0.0.1:1111",
			},
			err: false,
		},
		{
			name: "private",
			req: &smtpsrv.Request{
				From:       "hoge@anotherhost.localdomain",
				RemoteAddr: "172.16.0.1:1111",
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := main.StrictValidation(tt.req)
			if (err != nil) != (tt.err == true) {
				t.Fatal("fail")
			}
		})
	}
}
