package main

import (
	"fmt"
	"maps"
	"net/url"
	"sync"
)

type factory struct {
	urls          []*url.URL
	injectValues  []string
	gMode         string
	vMode         string
	chunks        int
	threads       int
	goroutinesCap chan struct{}
}

func (f *factory) generateUrls(customParams []string) ([]string, error) {
	switch f.gMode {
	case "ignore":
		return f.ignoreMode(customParams), nil

	case "pitchfork":
		return f.pitchforkMode(customParams), nil

	case "all":
		return f.allMode(customParams), nil
	}

	return nil, fmt.Errorf("[!] Invalid generate strategy mode")
}

func (f *factory) ignoreMode(customParams []string) []string {
	if len(customParams) == 0 {
		return nil
	}
	var (
		finalUrls = make([]string, 0, len(f.urls)*3)
		wg        sync.WaitGroup
		params    url.Values
	)

	for _, urlObj := range f.urls {
		f.goroutinesCap <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, injectValue := range f.injectValues {
				params = urlObj.Query()

				addCustomParams(params, customParams, injectValue)
				addToFinals(urlObj, &finalUrls, params)
			}
			<-f.goroutinesCap
		}()
	}

	wg.Wait()
	return finalUrls
}

func (f *factory) pitchforkMode(customParams []string) []string {
	var (
		finalUrls = make([]string, 0, len(f.urls)*3)
		params    url.Values
		wg        sync.WaitGroup
	)

	for _, urlObj := range f.urls {
		f.goroutinesCap <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			for _, injectValue := range f.injectValues {
				params = urlObj.Query()

				if len(params) == 0 {
					if len(customParams) == 0 {
						continue
					}
					emptyParamsMap := make(url.Values, len(customParams))

					addCustomParams(emptyParamsMap, customParams, injectValue)
					addToFinals(urlObj, &finalUrls, emptyParamsMap)
					continue
				}

				for param, val := range params {
					cpParams := make(url.Values, len(params)+len(params))
					maps.Copy(cpParams, params)
					cpParams[param] = []string{generateValue(val[0], injectValue, f.vMode)}

					addCustomParams(cpParams, customParams, injectValue)
					addToFinals(urlObj, &finalUrls, cpParams)
				}
			}
			<-f.goroutinesCap
		}()
	}

	wg.Wait()
	return finalUrls
}

func (f *factory) allMode(customParams []string) []string {
	allUrlsCh := make(chan string, 2)
	allUrls := make([]string, 0, len(f.urls)*3)
	checkUnique := make(map[string]bool, len(f.urls)*3)

	var wg sync.WaitGroup
	wg.Add(2)
	f.goroutinesCap <- struct{}{}
	f.goroutinesCap <- struct{}{}
	go func() {
		defer wg.Done()

		ignoreMode := f.ignoreMode(customParams)
		for _, gUrl := range ignoreMode {
			allUrlsCh <- gUrl
		}
	}()
	go func() {
		defer wg.Done()

		pitchforkMode := f.pitchforkMode(customParams)
		for _, gUrl := range pitchforkMode {
			allUrlsCh <- gUrl
		}
	}()
	go func() {
		wg.Wait()
		<-f.goroutinesCap
		<-f.goroutinesCap
		close(allUrlsCh)
	}()

	for gUrl := range allUrlsCh {
		if _, alreadyExists := checkUnique[gUrl]; alreadyExists {
			continue
		}

		allUrls = append(allUrls, gUrl)
		checkUnique[gUrl] = true
	}

	return allUrls
}
