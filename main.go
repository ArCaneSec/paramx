package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

func main() {
	f := factory{injectValues: []string{"'arcane'", `"arcane"`}}
	f.urls = make([]*url.URL, 0)

	var (
		customParams string
		urlsPath     string
		injectValues string
	)

	flag.StringVar(&f.gMode, "gs", "all",
		"strategy mode for generating urls [ignore | pitchfork | all]\n"+
			"ignore: will not touch current parameters, only appending new ones to the url (custom parameters required)\n"+
			"pitchfork: use inject value in each parameter separately\n"+
			"all: using all methods",
	)
	flag.StringVar(&f.vMode, "vs", "append",
		"strategy mode for generating values\n"+
			"append: appending inject value to parameter's current value\n"+
			"replace: replacing parameter's value with inject value",
	)

	flag.IntVar(&f.chunks, "c ", 40, "total number of parameter in each url")
	flag.IntVar(&f.threads, "t", 150, "maximum number of threads")

	flag.StringVar(&customParams, "p", "", "path to file containing parameters separated by \\n (required if using ignore mode)")
	flag.StringVar(&injectValues, "v", "", "path to file containing inject values by \\n")
	flag.StringVar(&urlsPath, "u", "", "path to file containing urls separated by \\n")

	flag.Parse()

	// if f.gMode == "ignore" && customParams == "" {
	// 	log.Fatalln("[!] ignore mode requires at least 1 custom parameter.")
	// }

	// if injectValues == "" {
	// 	log.Fatalln("[!] you must provide at least 1 inject value.")
	// }

	if urlsPath == "" {
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			log.Fatalln("[!] you must provide at least 1 url, using flags or stdin.")
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			urlStr := scanner.Text()
			urlObj, err := url.Parse(urlStr)

			if err != nil {
				log.Fatal(err)
			}

			f.urls = append(f.urls, urlObj)
		}
	}

	urls, err := f.generateUrls([]string{"cmd"})
	if err != nil {
		log.Fatalln(err)
	}
	for _, u := range urls {
		fmt.Println(u)
	}
}
