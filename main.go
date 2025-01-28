package main

import (
	"rms_proxy/v2/src/clientserver"
)

func main() {
	cs := &clientserver.ClientServer{}
	cs.StartServer()
}
