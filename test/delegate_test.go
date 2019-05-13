package test

import (
	"github.com/1-bi/log-api"
	"github.com/1-bi/log-zap"
	"github.com/1-bi/log-zap/appender"
	zaplayout "github.com/1-bi/log-zap/layout"
	"github.com/1-bi/servicebus"
	"github.com/coreos/etcd/clientv3"
	"log"
	"testing"
	"time"
)

func Test_Agent_BaseCase1(t *testing.T) {

	// --- create agent ---
	prepareLogSetting()

	conf, err := prepareConfig()

	if err != nil {
		logapi.GetLogger("agent_basecase1").Fatal(err.Error(), nil)
	}

	// dstart agent to complent
	var agent = servicebus.NewAgent(conf)

	agent.Start()

	defer agent.Stop()

	// --- check

	agent.Fire("agent.test.case", []byte("hello test"))

	//	runtime.Goexit()
	// stop main thread running
	select {}

}

func prepareConfig() (*servicebus.Config, error) {

	var conf = servicebus.NewConfig()

	conf.SetRegisterServer(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})

	conf.SetNodeRoles([]string{
		"master", "minion",
	})

	var configErr = conf.CheckBeforeStart()

	return conf, configErr

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
