package main

import (
	"github.com/1-bi/cron"
	"github.com/1-bi/log-api"
	"github.com/1-bi/log-zap"
	"github.com/1-bi/log-zap/appender"
	zaplayout "github.com/1-bi/log-zap/layout"
	"github.com/1-bi/servicebus"
	"github.com/coreos/etcd/clientv3"
	"log"
	"runtime"
	"time"
)

func main() {
	prepareLogSetting()
	// --- set the global logger ---

	var conf = servicebus.NewConfig()

	conf.SetRegisterServer(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})

	conf.SetNodeRoles([]string{
		"master", "minion",
	})

	var configErr = conf.CheckBeforeStart()

	if configErr != nil {
		logapi.GetLogger("agent").Fatal(configErr.Error(), nil)
		return
	}

	// detect the properties
	var agent = servicebus.NewAgent(conf)

	agent.Start()

	defer agent.Stop()

	// connect api --

	// ---- keep program running ----
	runtime.Goexit()

}

func prepareLogSetting() {

	// --- construct layout ---
	var jsonLayout = zaplayout.NewJsonLayout()
	//jsonLayout.SetTimeFormat("2006-01-02 15:04:05")
	jsonLayout.SetTimeFormat("2006-01-02 15:04:05 +0800 CST")
	//fmt.Println( time.Now().Location() )

	// --- set appender
	var consoleAppender = appender.NewConsoleAppender(jsonLayout)

	var mainOpt = logzap.NewLoggerOption()
	mainOpt.SetLevel("debug")
	mainOpt.AddAppender(consoleAppender)

	var agentOpt = logzap.NewLoggerOption()
	agentOpt.SetLoggerPattern("servicebus")
	agentOpt.SetLevel("warn")
	agentOpt.AddAppender(consoleAppender)

	var implReg = new(logzap.ZapFactoryRegister)

	_, err := logapi.RegisterLoggerFactory(implReg, mainOpt, agentOpt)

	if err != nil {
		log.Fatal(err)
	}
}

func regServer() {
	/*
		serviceName := "s-test"
		serviceInfo := servicebus.AgentInfo{IP: "vicenteyou"}

		s, err := servicebus.NewAgentRegisterService(serviceName, serviceInfo, []string{
			"http://localhost:2379",
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("name:%s, ip:%s\n", s.nodeId, s.Info.IP)

		go func() {
			time.Sleep(time.Second * 20)
			s.Stop()
		}()

		s.Start()
	*/
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
