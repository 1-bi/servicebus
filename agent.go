package servicebus

import (
	"fmt"
	"github.com/1-bi/log-api"
	"github.com/1-bi/servicebus/schema"
	"github.com/coreos/etcd/clientv3"
)

// Agent define service bus agent proxy
type Agent struct {
	conf *Config
}

func (myself *Agent) Start() {

	// --- open thread
	go myself.startRegisterServer()

	// --- start watch server
	go myself.startWatchServer()

}

func (myself *Agent) Stop() {

}

// On implement event name
func (myself *Agent) On(eventName string, fn func(ServiceEventHandler)) error {

	return nil
}

// Fire call by event name and define callback
func (myself *Agent) Fire(eventName string, msgBody []byte, callback ...Callback) error {

	// serialization runtimeArgs
	reqEvent := new(schema.ReqEvent)
	reqEvent.Name = eventName
	reqEvent.ParamsBody = msgBody

	// --- sent msg body ---

	return nil
}

// ---------------------  private method ---
func (myself *Agent) startRegisterServer() {

	cli, err := clientv3.New(myself.conf._etcdConfig)

	if err != nil {
		structBean := logapi.NewStructBean()
		structBean.LogStringArray("etcd.server", myself.conf._etcdConfig.Endpoints)
		logapi.GetLogger("serviebus.agent").Fatal("Connect etcd server fail.", structBean)
		return
	}

	var nodeRoles = []string{"master", "minion"}
	if len(myself.conf.nodeRoles) == 0 {
		nodeRoles = myself.conf.nodeRoles
	}

	var serv = NewAgentRegisterService(myself.conf._agentNodeId, cli, nodeRoles)

	err = serv.Start()
	if err != nil {
		fmt.Println(err)
	}

}

func (myself *Agent) startWatchServer() {

	cli, err := clientv3.New(myself.conf._etcdConfig)

	if err != nil {
		structBean := logapi.NewStructBean()
		structBean.LogStringArray("etcd.server", myself.conf._etcdConfig.Endpoints)
		logapi.GetLogger("serviebus.agent").Fatal("Connect etcd server fail.", structBean)
		return
	}

	var serv = NewAgentWatchService(myself.conf._agentNodeId, cli)

	err = serv.Start()
	if err != nil {
		fmt.Println(err)
	}

}

func (myself *Agent) checkRegCenterConnect() {

}

func (myself *Agent) getNodeAgentList() {

}

// NewAgent check agent
func NewAgent(conf *Config) *Agent {

	var agentLogger = logapi.GetLogger("servicebus.Agent")

	agentLogger.Info("base", nil)

	var agent = new(Agent)
	agent.conf = conf

	//  start scheduler

	// check health

	// start describe queue to nats
	return agent

}
