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
	filePtr := flag.String("f", "", "Template filename, default STDIN")
	flag.Parse()

	var body []byte
	var b string
	var err error

	if *filePtr != "" {
		body, err = os.ReadFile(*filePtr)
	} else {
		body, err = io.ReadAll(os.Stdin)
	}

	if err != nil {
		panic(err)
	}
	b = string(body)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		for _, from := range []string{
			fmt.Sprintf("\\$%s\\b", regexp.QuoteMeta(pair[0])),
			fmt.Sprintf("\\${%s}", regexp.QuoteMeta(pair[0])),
		} {
			r := regexp.MustCompile(from)
			b = r.ReplaceAllString(b, pair[1])
		}
	}
	fmt.Println(b)
}
