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
	)

	flag.StringVar(&f.gMode, "generate strategy mode", "all", "strategy mode for generating urls")
	flag.StringVar(&f.vMode, "value strategy mode", "append", "strategy mode for generating urls")
	flag.IntVar(&f.chunks, "chunks ", 40, "total number of parameter in each url")

	flag.StringVar(&customParams, "params", "", "path to file containing parameters separated by \\n")
	flag.StringVar(&urlsPath, "urls", "", "path to file containing urls separated by \\n")

	flag.Parse()
	// if f.gMode == "ignore" && customParams == "" {
	// 	log.Fatalln("[!] ignore mode requires at least 1 custom parameter.")
	// }

	if urlsPath == "" {
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			os.Exit(0)
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

	urls, err := f.generateUrls([]string{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, u := range urls {
		fmt.Println(u)
	}
}
