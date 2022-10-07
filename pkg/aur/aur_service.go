package aur

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
)

func FindPackageByNameInAur(name string) (*ResponseInfo, error) {

	var response ResponseInfo

	query := fmt.Sprintf("by=name&arg=%s", name)

	value, err := url.ParseQuery(query)

	if err != nil {
		logrus.Errorf("Cannot parse query string [error: %s]", err)
		return nil, err
	}

	Call(Search, value, &response)
	return &response, nil
}
