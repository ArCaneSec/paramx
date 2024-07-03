package main

import (
	"fmt"
	"net/url"
)

func addCustomParams(currentParams url.Values, customParams []string, injectValue string) {
	for _, param := range customParams {
		currentParams[param] = []string{injectValue}
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
	cpObj.RawQuery = params.Encode()
	*finalUrls = append(*finalUrls, cpObj.String())
}
