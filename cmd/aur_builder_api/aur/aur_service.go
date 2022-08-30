package aur

import (
	"fmt"
	"net/url"

	aurgo "github.com/drzombey/aur-rpc-client-go"
	"github.com/sirupsen/logrus"
)

func FindPackageByNameInAur(name string) {

	var response aurgo.ResponseInfo

	query := fmt.Sprintf("by=name&arg=%s", name)

	value, err := url.ParseQuery(query)

	if err != nil {
		logrus.Errorf("Cannot parse query string [error: %s]", err)
		return
	}

	aurgo.Call(aurgo.Search, value, &response)
	fmt.Println(response)
}
