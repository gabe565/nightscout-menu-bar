package config

var reloaders []func()

func AddReloader(fn func()) {
	reloaders = append(reloaders, fn)
}
