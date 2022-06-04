package main_test

import (
	"testing"

	main "github.com/tsg-ut/smtp2http"
)

func TestStrictValidation(t *testing.T) {
	tests := []struct {
		name string
		req  *main.TestContext
		err  bool // err should be non-nil or not
	}{
		{
			name: "hoge",
			req: &main.TestContext{
				FromEmail:        "hoge@example.com",
				RemoteAddrString: "1.1.1.1:1111",
			},
			err: true,
		},
		{
			name: "fake-localhost",
			req: &main.TestContext{
				FromEmail:        "hoge@localhost.localdomain",
				RemoteAddrString: "1.1.1.1:1111",
			},
			err: true,
		},
		{
			name: "fake-gmail",
			req: &main.TestContext{
				FromEmail:        "hoge@gmail.com",
				RemoteAddrString: "1.1.1.1:1111",
			},
			err: true,
		},
		{
			name: "valid-gmail",
			req: &main.TestContext{
				FromEmail:        "hoge@gmail.com",
				RemoteAddrString: "35.190.247.1:1111",
			},
			err: false,
		},
		{
			name: "lo",
			req: &main.TestContext{
				FromEmail:        "hoge@localhost.localdomain",
				RemoteAddrString: "127.0.0.1:1111",
			},
			err: false,
		},
		{
			name: "private",
			req: &main.TestContext{
				FromEmail:        "hoge@anotherhost.localdomain",
				RemoteAddrString: "172.16.0.1:1111",
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
