// +build integration

package main

import (
	"os"
	"testing"

	"github.com/prokosna/dementor/lib"
)

// Integration test entry point
func TestMain(m *testing.M) {
	dementor.InitConf()
	code := m.Run()
	defer os.Exit(code)
}
