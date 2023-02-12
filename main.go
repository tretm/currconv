package main

import (
	"os"
	"os/signal"

	ccc "currconv/cc"

	pars "currconv/cc/parser"
	serv "currconv/cc/server"
	stor "currconv/cc/storage"
)

func main() {

	ccc.Debug = true
	stg := stor.New()

	exitch := make(chan struct{})

	pars.Run(stg, exitch)

	srv := &serv.HttpServer{St: stg}
	srv.HttpApi()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Stop(c)

	close(exitch)
	srv.HttpApiStop()
	stg.Stop()

}
