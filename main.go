package main

import (
	"log/slog"
	"os"
	"rms_proxy/v2/src/clientserver"

	"github.com/MatusOllah/slogcolor"
)

func main() {
	 slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions)))
	cs := &clientserver.ClientServer{}
	cs.StartServer()
}
