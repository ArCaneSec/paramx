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
	f := factory{urls: make([]*url.URL, 0)}

	var (
		customParamsPath string
		urlsPath         string
		injectValuesPath string
		totalThreads     int
	)

	flag.StringVar(&f.gMode, "gs", "ignore",
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

	flag.IntVar(&f.chunks, "c", 40, "total number of parameter in each url")
	flag.IntVar(&totalThreads, "t", 150, "maximum number of threads")

	flag.StringVar(&customParamsPath, "p", "", "path to parameters file separated by \\n (required if using ignore mode)")
	flag.StringVar(&injectValuesPath, "v", "", "path to values file separated by \\n")
	flag.StringVar(&urlsPath, "u", "", "path to urls file separated by \\n")
	flag.Parse()

	flag.Usage = func (){
		printAscii()
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(os.Args) == 1{
		flag.Usage()
	}

	// least amount of threads needed, when using all generate strategy, check factory.go
	if totalThreads < 3 {
		totalThreads = 3
	}
	f.goroutinesCap = make(chan struct{}, totalThreads)

	if f.gMode == "ignore" && customParamsPath == "" {
		log.Fatalln("[!] ignore mode requires at least 1 custom parameter.")
	}

	if injectValuesPath == "" {
		log.Fatalln("[!] you must provide at least 1 inject value.")
	}

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
			parseUrl(urlStr, &f.urls)
		}
	} else {
		urlLines, err := readFile(urlsPath)
		if err != nil {
			log.Fatalln(err)
		}

		for _, rawUrl := range urlLines {
			parseUrl(rawUrl, &f.urls)
		}
	}

	var params []string
	if customParamsPath != "" {
		paramLines, err := readFile(customParamsPath)
		if err != nil {
			log.Fatalln(err)
		}

		params = paramLines
	}

	parsedValues, err := readFile(injectValuesPath)
	if err != nil {
		log.Fatalln(err)
	}
	f.injectValues = parsedValues

	urls, err := f.generateUrls(params)
	if err != nil {
		log.Fatalln(err)
	}

	for _, u := range urls {
		fmt.Println(u)
	}
}
