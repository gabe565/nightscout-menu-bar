package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	flag "github.com/spf13/pflag"
)

//go:embed nightscout-menu-bar.rb.tmpl
var spec string

const binaryName = "nightscout-menu-bar"

// Platforms lists the OS/arch pairs published in the cask, in output order.
//
//nolint:gochecknoglobals
var Platforms = []Platform{
	{OS: "darwin", Arch: "arm64"},
	{OS: "darwin", Arch: "amd64"},
	{OS: "linux", Arch: "arm64"},
	{OS: "linux", Arch: "amd64"},
}

type Platform struct {
	OS, Arch string
}

// Archive returns the release asset filename for the platform.
func (p Platform) Archive() string {
	return binaryName + "_" + p.OS + "_" + p.Arch + ".tar.gz"
}

type SpecVars struct {
	// Version without the leading "v"
	Version string
	// SHA256 maps "<os>_<arch>" to the archive's hex-encoded checksum
	SHA256 map[string]string
}

func main() {
	var dist string
	flag.StringVar(&dist, "dist", "dist", "Directory containing the release archives")
	flag.Parse()

	if flag.NArg() < 1 {
		slog.Error("Usage: cask [--dist DIR] VERSION")
		os.Exit(1)
	}

	version := flag.Arg(0)

	values := SpecVars{
		Version: strings.TrimPrefix(version, "v"),
		SHA256:  make(map[string]string, len(Platforms)),
	}

	for _, p := range Platforms {
		sum, err := checksum(filepath.Join(dist, p.Archive()))
		if err != nil {
			slog.Error("computing checksum", "err", err)
			os.Exit(1)
		}
		values.SHA256[p.OS+"_"+p.Arch] = sum
	}

	tmpl, err := template.New("").Parse(spec)
	if err != nil {
		slog.Error("parsing template", "err", err)
		os.Exit(1)
	}

	if err := tmpl.Execute(os.Stdout, values); err != nil {
		slog.Error("executing template", "err", err)
		os.Exit(1)
	}
}

func checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
