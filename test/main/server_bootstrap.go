package main

import (
	"github.com/1-bi/servicebus"
	"log"
	"runtime"
)

func main() {

	// ---- create service bus manager ----

	servusManager, err := servicebus.NewServiceManager()

	if err != nil {
		log.Fatal(err)
	}

	servusManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}
