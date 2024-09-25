package settings

import (
	"github.com/fynelabs/selfupdate"
)

func restart() error {
	updater, err := selfupdate.Manage(&selfupdate.Config{})
	if err != nil {
		return err
	}
	return updater.Restart()
}
