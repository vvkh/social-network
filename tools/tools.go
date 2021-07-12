// +build tools

package tools

//go:generate go build -o ../bin github.com/golang/mock/mockgen
//go:generate go build -o ../bin github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go build -o ../bin golang.org/x/tools/cmd/goimports

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/tools/cmd/goimports"
)
