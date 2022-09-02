package aur

import (
	"fmt"
	"net/url"

	aurgo "github.com/drzombey/aur-rpc-client-go"
	"github.com/drzombey/aur-rpc-client-go/types"
	"github.com/sirupsen/logrus"
)

func FindPackageByNameInAur(name string) (*types.ResponseInfo, error) {

	var response types.ResponseInfo

	query := fmt.Sprintf("by=name&arg=%s", name)

	value, err := url.ParseQuery(query)

	if err != nil {
		logrus.Errorf("Cannot parse query string [error: %s]", err)
		return nil, err
	}

	aurgo.Call(types.Search, value, &response)
	return &response, nil
}
