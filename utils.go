package main

import (
	"fmt"
	"log"
	"maps"
	"net/url"
	"os"
	"strings"
)

func addCustomParams(currentParams url.Values, customParams []string, injectValue string) {
	for _, param := range customParams {
		if _, exists := currentParams[param]; !exists {
			currentParams[param] = []string{injectValue}
		}
	}
}

func generateValue(param, injectValue, vMode string) string {
	switch vMode {
	case "append":
		return fmt.Sprintf("%s%s", param, injectValue)
	case "replace":
		return injectValue
	}
	return ""
}

func addToFinals(urlObj *url.URL, finalUrls *[]string, params url.Values) {
	cpObj := *urlObj
	newParamsRaw := params.Encode()

	// happens when custom parameter was existed in url, and gs mode is ignore
	// in this case, url will not receive any change at all, so we skip it.
	if newParamsRaw == urlObj.RawQuery {
		return
	}
	cpObj.RawQuery = newParamsRaw

	*finalUrls = append(*finalUrls, cpObj.String())
}

func parseUrl(rawUrl string, urls *[]*url.URL) {
	urlObj, err := url.Parse(rawUrl)
	if urlObj.Scheme == "" {
		urlObj.Scheme = "https"
	}

	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasPrefix(urlObj.Path, "/") {
		urlObj = urlObj.JoinPath("/")
	}
	*urls = append(*urls, urlObj)
}

func readFile(fileName string) ([]string, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	sLines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	if len(sLines) == 1 && sLines[0] == "" {
		return nil, fmt.Errorf("[!] %s's file is empty", fileName)
	}

	return sLines, nil
}

func printAscii() {
	fmt.Println(`
  ____   _    ____      _    __  ____  __
 |  _ \ / \  |  _ \    / \  |  \/  \ \/ /
 | |_) / _ \ | |_) |  / _ \ | |\/| |\  / 
 |  __/ ___ \|  _ <  / ___ \| |  | |/  \ 
 |_| /_/   \_\_| \_\/_/   \_\_|  |_/_/\_\                                                            
	`)
}

// Separate custom parameters in different urls
// Useful when total parameters are higher than requested chunk.
// It slices the overflowed amount and creates a valid url with legal chunks
// slicing continues until overflow amount goes below 0.
func addParamsSeparately(urlObj *url.URL, customParams []string, injectValue string, finalUrls *[]string, params url.Values, chunksOverflow, chunks int) {
	for chunksOverflow > 0 {
		overflowParams := customParams[:len(customParams)-chunksOverflow]
		customParams = customParams[len(customParams)-chunksOverflow:]
		chunksOverflow -= len(overflowParams)

		cpParams := make(url.Values, chunks)
		maps.Copy(cpParams, params)

		addCustomParams(cpParams, overflowParams, injectValue)
		addToFinals(urlObj, finalUrls, cpParams)
	}
	addCustomParams(params, customParams, injectValue)
	addToFinals(urlObj, finalUrls, params)
}
