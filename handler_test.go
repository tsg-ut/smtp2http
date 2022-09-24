package main_test

import (
	"testing"

	main "github.com/tsg-ut/smtp2http"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		req   *main.TestContext
		valid bool
	}{
		{
			name: "hoge",
			req: &main.TestContext{
				FromEmail:        "hoge@example.com",
				RemoteAddrString: "1.1.1.1:1111",
			},
			valid: false,
		},
		{
			name: "fake-localhost",
			req: &main.TestContext{
				FromEmail:        "hoge@localhost.localdomain",
				RemoteAddrString: "1.1.1.1:1111",
			},
			valid: false,
		},
		{
			name: "fake-gmail",
			req: &main.TestContext{
				FromEmail:        "hoge@gmail.com",
				RemoteAddrString: "1.1.1.1:1111",
			},
			valid: false,
		},
		{
			name: "valid-gmail",
			req: &main.TestContext{
				FromEmail:        "hoge@gmail.com",
				RemoteAddrString: "35.190.247.1:1111",
			},
			valid: true,
		},
		{
			name: "lo",
			req: &main.TestContext{
				FromEmail:        "hoge@localhost.localdomain",
				RemoteAddrString: "127.0.0.1:1111",
			},
			valid: true,
		},
		{
			name: "private",
			req: &main.TestContext{
				FromEmail:        "hoge@anotherhost.localdomain",
				RemoteAddrString: "172.16.0.1:1111",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := main.Validate(tt.req)
			if valid != tt.valid {
				t.Fatal("fail")
			}
		})
	}
}
