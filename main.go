package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	var (
		err  error
		body string
		bs   []byte
	)

	// flags ...
	filePtr := flag.String("f", "", "Template filename, default STDIN")
	flag.Parse()

	// read template body ...
	if *filePtr != "" {
		bs, err = os.ReadFile(*filePtr)
	} else {
		bs, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		panic(err)
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
