//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/princjef/mageutil/bintool"
	"github.com/princjef/mageutil/shellcmd"
)

var (
	gofumpt = bintool.Must(bintool.New(
		"gofumpt{{.BinExt}}",
		"v0.6.0",
		"https://github.com/mvdan/gofumpt/releases/download/{{.Version}}/gofumpt_{{.Version}}_{{.GOOS}}_{{.GOARCH}}{{.BinExt}}",
	))
	goimports = bintool.Must(bintool.NewGo(
		"golang.org/x/tools/cmd/goimports",
		"v0.21.0",
		bintool.WithVersionCmd(""),
	))
	golines = bintool.Must(bintool.NewGo(
		"github.com/segmentio/golines",
		"v0.11.0",
	))
	linter = bintool.Must(bintool.New(
		"golangci-lint{{.BinExt}}",
		"1.55.2",
		"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
	))
	documenter = bintool.Must(bintool.New(
		"gomarkdoc{{.BinExt}}",
		"0.4.1",
		"https://github.com/princjef/gomarkdoc/releases/download/v{{.Version}}/gomarkdoc_{{.Version}}_{{.GOOS}}_{{.GOARCH}}{{.ArchiveExt}}",
	))
	integrationHash = "integrationHash.out"
)

func EnsureFormatters() error {
	if err := gofumpt.Ensure(); err != nil {
		return err
	}
	if err := goimports.Ensure(); err != nil {
		return err
	}
	if err := golines.Ensure(); err != nil {
		return err
	}
	return nil
}

func EnsureLinter() error {
	return linter.Ensure()
}

func EnsureDocumenter() error {
	return documenter.Ensure()
}

func EnsureAllTools() error {
	if err := EnsureFormatters(); err != nil {
		return err
	}
	if err := EnsureLinter(); err != nil {
		return err
	}
	if err := EnsureDocumenter(); err != nil {
		return err
	}
	return nil
}

func Format() error {
	if err := EnsureFormatters(); err != nil {
		return err
	}
	if err := gofumpt.Command("-w .").Run(); err != nil {
		return err
	}
	if err := goimports.Command("-w .").Run(); err != nil {
		return err
	}
	if err := golines.Command("-m 120 --no-reformat-tags --base-formatter gofmt -w .").Run(); err != nil {
		return err
	}
	return nil
}

func Lint() error {
	if err := EnsureLinter(); err != nil {
		return err
	}
	return linter.Command("run").Run()
}

func Doc() error {
	if err := EnsureDocumenter(); err != nil {
		return err
	}
	return shellcmd.RunAll(
		documenter.Command("./lib/..."),
		documenter.Command("./cmd/..."),
	)
}

func DocVerify() error {
	if err := EnsureDocumenter(); err != nil {
		return err
	}
	return shellcmd.RunAll(
		documenter.Command("-c ./lib/..."),
		documenter.Command("-c ./cmd/..."),
	)
}

func Build() error {
	return shellcmd.RunAll(
		`go build -o bin/app ./cmd/app`,
	)
}

func BuildDebug() error {
	return shellcmd.RunAll(
		`go build -gcflags="all=-N -l" -o bin/app ./cmd/app`,
	)
}

func Test() error {
	return shellcmd.Command(`go test -timeout 30s -cover ./...`).Run()
}

func TestRace() error {
	return shellcmd.Command(`go test -race -timeout 30s -cover ./...`).Run()
}

func TestSingle(folder, filter string) error {
	return shellcmd.Command(fmt.Sprintf(`go test -timeout 5m -p=1 -count=1 -cover %s -run %s`, folder, filter)).Run()
}

func Cover() error {
	return shellcmd.RunAll(
		`go test -coverprofile=coverage.out ./...`,
		`go tool cover -html=coverage.out`,
	)
}

func CI() error {
	if err := Format(); err != nil {
		return err
	}
	if err := Lint(); err != nil {
		return err
	}
	if err := Doc(); err != nil {
		return err
	}
	if err := Test(); err != nil {
		return err
	}
	return nil
}

func CIVerify() error {
	if err := Format(); err != nil {
		return err
	}
	if err := Lint(); err != nil {
		return err
	}
	if err := DocVerify(); err != nil {
		return err
	}
	if err := Test(); err != nil {
		return err
	}
	return nil
}

func Start() error {
	if err := Build(); err != nil {
		return err
	}
	return shellcmd.Command(`./bin/app`).Run()
}
