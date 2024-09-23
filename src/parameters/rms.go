package parameters

import "net/http"

type RMSConnectParameter struct {
	url      string
	login    string
	password string
}

func(rm *RMSConnectParameter) Handle( r *http.Request) (string,error) {

	return "", nil;
}
