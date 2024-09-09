package util

import (
	"net/http"
	"strings"
)

func NewUserAgentTransport(name, version string) *UserAgentTransport {
	ua := name + "/"
	commit := strings.TrimPrefix(GetCommit(), "*")
	if version != "" {
		ua += version
		if commit != "" {
			ua += "-" + commit
		}
	} else if commit != "" {
		ua += commit
	}

	return &UserAgentTransport{
		transport: http.DefaultTransport,
		userAgent: ua,
	}
}

type UserAgentTransport struct {
	transport http.RoundTripper
	userAgent string
}

func (u *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", u.userAgent)
	return u.transport.RoundTrip(r)
}
