package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"text/template"
)

type Entry struct {
	Value  string
	Config map[string]Entry
}

func createMap(s []string, v string, m map[string]Entry) map[string]Entry {
	if m == nil {
		m = make(map[string]Entry)
	}
	var entry Entry
	if val, ok := m[s[0]]; ok {
		entry = val
	}

	if len(s) > 1 {
		entry.Config = createMap(s[1:], v, entry.Config)
	} else {
		entry.Value = v
	}
	m[s[0]] = entry

	return m
}

func createTemplate() {

	prefix := os.Getenv("TEMPLATER_PREFIX")
	outputFile := os.Getenv("TEMPLATER_" + prefix + "_OUTPUT")
	filename := path.Base(outputFile)

	tmpl, err := template.ParseGlob("*.tmpl")
	if err != nil {
		panic(err)
	}

	data := make(map[string]map[string]Entry)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], prefix) {
			options := strings.Split(strings.TrimPrefix(pair[0], prefix+"_"), "_")
			section := options[0]
			var m map[string]Entry
			if val, ok := data[section]; ok {
				m = val
			}
			data[section] = createMap(options[1:], pair[1], m)
		}
	}
	output, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(output, filename+".tmpl", data)
	output.Close()
}

func runCmd() {
	argsWithoutProg := os.Args[1:]
	cmd := exec.Command(argsWithoutProg[0], argsWithoutProg[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			fmt.Println("cmd.Wait: %v", err)
		}
	}
}

func main() {
	createTemplate()
	runCmd()
}
