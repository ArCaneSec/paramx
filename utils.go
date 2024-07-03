package main

import (
	"fmt"
	"net/url"
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
