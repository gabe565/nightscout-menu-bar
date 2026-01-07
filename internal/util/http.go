package util

import (
	"runtime"
	"strings"

	"gabe565.com/utils/httpx"
)

func NewUserAgentTransport(name, version string) *httpx.UserAgentTransport {
	var ua strings.Builder
	ua.WriteString(name)
	ua.WriteRune('/')
	commit := strings.TrimPrefix(GetCommit(), "*")
	if version != "" {
		ua.WriteString(version)
		if commit != "" {
			ua.WriteRune('-')
			ua.WriteString(commit)
		}
	} else if commit != "" {
		ua.WriteString(commit)
	}
	ua.WriteString(" (" + runtime.GOOS + "/" + runtime.GOARCH + ")")

	return httpx.NewUserAgentTransport(nil, ua.String())
}
