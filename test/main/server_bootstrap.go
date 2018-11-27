package main

import (
	"github.com/1-bi/servicebus"
	"log"
	"runtime"
)

func main() {

	// ---- create service bus manager ----

	serviceManager, err := servicebus.NewServiceManager()

	if err != nil {
		log.Fatal(err)
	}

	serviceManager.On()

	serviceManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}
