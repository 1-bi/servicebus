package main

import (
	"fmt"
	"github.com/1-bi/cron"
	dis "github.com/1-bi/servicebus"
	"log"
	"time"
)

func main() {

	serviceName := "s-test"
	serviceInfo := dis.ServiceInfo{IP: "vicenteyou"}

	s, err := dis.NewService(serviceName, serviceInfo, []string{
		"http://localhost:2379",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name:%s, ip:%s\n", s.Name, s.Info.IP)

	go func() {
		time.Sleep(time.Second * 20)
		s.Stop()
	}()

	s.Start()
}

func myFunc() {

	i := 0
	c := cron.New()
	spec := "@every 2s"
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}