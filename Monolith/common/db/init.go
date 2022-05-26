package db

import (
	"log"
	"os"
	"runtime"
	"strings"
)

var configPath string

func init() {
	var sbuilder strings.Builder
	var separator string

	sysType := runtime.GOOS
	if sysType == "linux" {
		separator = "/"
	} else if sysType == "windows" {
		separator = "\\"
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	strs := []string{cwd, "common", "config"}
	for i := 0; i < len(strs); i++ {
		sbuilder.WriteString(strs[i])
		sbuilder.WriteString(separator)
	}
	sbuilder.WriteString("application.yaml")
	configPath = sbuilder.String()
}
