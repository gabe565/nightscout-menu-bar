package util

import "runtime/debug"

func GetCommit() string {
	var commit string
	var modified bool
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					modified = true
				}
			}
		}
	}

	if commit != "" {
		if len(commit) > 8 {
			commit = commit[:8]
		}
		if modified {
			commit = "*" + commit
		}
	}
	return commit
}
