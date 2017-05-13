// +build integration

package dementor

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitConf()
	code := m.Run()
	defer os.Exit(code)
}
