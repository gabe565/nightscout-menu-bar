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

var (
	//go:embed PKGBUILD.tmpl
	pkgbuild string
	//go:embed SRCINFO.tmpl
	srcinfo string
)

const binaryName = "nightscout-menu-bar"

type SpecVars struct {
	Tag         string
	Version     string
	Description string
	Sums        map[string]string
}

func (s SpecVars) URL(arch string) string {
	return "https://github.com/gabe565/nightscout-menu-bar/releases/download/" + s.Tag + "/" + archive(arch)
}

func (s SpecVars) SHA256(arch string) string {
	return s.Sums[arch]
}

func archive(arch string) string {
	return binaryName + "_linux_" + arch + ".tar.gz"
}

func main() {
	var dist, out string
	flag.StringVar(&dist, "dist", "dist", "Directory containing the release archives")
	flag.StringVar(&out, "out", ".", "Directory to write PKGBUILD and .SRCINFO to")
	flag.Parse()

	if flag.NArg() < 1 {
		slog.Error("Usage: aur [--dist DIR] [--out DIR] VERSION")
		os.Exit(1)
	}

	tag := flag.Arg(0)
	values := SpecVars{
		Tag: tag,
		// pkgver can't contain hyphens
		Version:     strings.ReplaceAll(strings.TrimPrefix(tag, "v"), "-", "_"),
		Description: "Display live blood sugar data from Nightscout in your system tray",
		Sums:        make(map[string]string, 2),
	}

	for _, arch := range []string{"amd64", "arm64"} {
		sum, err := checksum(filepath.Join(dist, archive(arch)))
		if err != nil {
			slog.Error("computing checksum", "err", err)
			os.Exit(1)
		}
		values.Sums[arch] = sum
	}

	for name, spec := range map[string]string{"PKGBUILD": pkgbuild, ".SRCINFO": srcinfo} {
		if err := render(filepath.Join(out, name), spec, values); err != nil {
			slog.Error("rendering "+name, "err", err)
			os.Exit(1)
		}
	}
}

func render(path, spec string, values SpecVars) error {
	tmpl, err := template.New("").Parse(spec)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(f, values); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
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
