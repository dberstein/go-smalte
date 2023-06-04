package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

func Usage(exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "%s <template file>\n... or ...\ncat <template file> | %s [-]\n", path.Base(os.Args[0]), path.Base(os.Args[0]))
	os.Exit(exitCode)
}

func main() {
	var (
		err  error
		body string
		bs   []byte
	)

	// help?
	if len(os.Args) > 1 {
		for _, value := range os.Args {
			if value == "-h" || value == "--help" {
				Usage(0)
			}
		}
	}

	// read template body ...
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		if len(os.Args) == 1 || os.Args[1] == "-" {
			bs, err = io.ReadAll(os.Stdin)
		} else {
			Usage(1)
		}
	} else {
		if len(os.Args) == 2 {
			bs, err = os.ReadFile(os.Args[1])
		} else {
			Usage(1)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(2)
	}
	body = string(bs)

	// replace each available environmental variable ...
	for _, e := range os.Environ() {
		env := strings.SplitN(e, "=", 2)
		quoted := regexp.QuoteMeta(env[0])
		// ... replace "$VAR" and "${VAR}" ...
		for _, from := range []string{
			fmt.Sprintf("\\$%s\\b", quoted),
			fmt.Sprintf("\\${%s}", quoted),
		} {
			r := regexp.MustCompile(from)
			body = r.ReplaceAllString(body, env[1])
		}
	}

	// ... output replaced text
	fmt.Println(body)
}
