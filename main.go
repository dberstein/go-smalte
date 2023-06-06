package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	isatty "github.com/mattn/go-isatty"
)

func Usage(exitCode int) {
	out := os.Stdout
	if exitCode > 0 {
		out = os.Stderr
	}
	fmt.Fprintf(out, "Usage:\n%s <template file>\n... or ...\ncat <template file> | %s [-]\n", path.Base(os.Args[0]), path.Base(os.Args[0]))
	os.Exit(exitCode)
}

// getSortedEnvironment returns environment sorted by variable name length, longest first
func getSortedEnvironment() []string {
	environment := os.Environ()
	// ... longest strings first ...
	sort.Slice(environment, func(i, j int) bool {
		envI := strings.SplitN(environment[i], "=", 2)[0]
		envJ := strings.SplitN(environment[j], "=", 2)[0]
		return len(envI) >= len(envJ)
	})
	return environment
}

// readTemplateBody returns template body read from STDIN or file
func readTemplateBody() (string, error) {
	var (
		err error
		bs  []byte
	)

	if !isatty.IsTerminal(os.Stdin.Fd()) {
		if len(os.Args) == 1 || os.Args[1] == "-" {
			bs, err = io.ReadAll(os.Stdin)
		} else {
			return string(bs), fmt.Errorf("cannot read from STDIN")
		}
	} else {
		if len(os.Args) == 2 {
			bs, err = os.ReadFile(os.Args[1])
		} else {
			return string(bs), fmt.Errorf("cannot read from file: " + strings.Join(os.Args[1:], ", "))
		}
	}

	return string(bs), err
}

func main() {
	// help?
	if len(os.Args) > 1 {
		for _, value := range os.Args {
			if value == "-h" || value == "--help" {
				Usage(0)
			}
		}
	}

	// read template body ...
	body, err := readTemplateBody()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		Usage(1)
	}

	// replace each available environmental variable ...
	for _, e := range getSortedEnvironment() {
		env := strings.SplitN(e, "=", 2)
		quoted := regexp.QuoteMeta(env[0])
		// ... replace "$VAR" and "${VAR}", but not "\$VAR" or "\${VAR}" ...
		for _, from := range []string{
			fmt.Sprintf("(^|[^\\\\])\\$%s\\b", quoted),
			fmt.Sprintf("(^|[^\\\\])\\${%s}", quoted),
		} {
			r := regexp.MustCompile(from)
			body = r.ReplaceAllString(body, "$1"+env[1])
		}
	}

	// ... output replaced text
	fmt.Print(body)
}
