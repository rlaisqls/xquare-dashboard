//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"

	"github.com/xquare-dashboard/pkg/build"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	os.Exit(build.RunCmd())
}
