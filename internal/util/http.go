package util

import (
	"runtime"
	"strings"

	"gabe565.com/utils/httpx"
)

func NewUserAgentTransport(name, version string) *httpx.UserAgentTransport {
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
	ua += " (" + runtime.GOOS + "/" + runtime.GOARCH + ")"

	return httpx.NewUserAgentTransport(nil, ua)
}
