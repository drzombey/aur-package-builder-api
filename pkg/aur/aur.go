package aur

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	aurBaseURL = "https://aur.archlinux.org/rpc"
	aurVersion = "5"
	useragent  = "aur-go-rpc-client"
)

var (
	client = new(http.Client)
)

func SetUserAgent(value string) {
	useragent = value
}

func Call(searchType SearchType, value url.Values, response interface{}) error {
	url := fmt.Sprintf("%s?v=%s&type=%s&%s",
		aurBaseURL,
		aurVersion,
		searchType.String(),
		value.Encode(),
	)

	payload, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	payload.Header.Add("User-Agent", useragent)

	resource, err := client.Do(payload)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resource.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}

	return nil
}
