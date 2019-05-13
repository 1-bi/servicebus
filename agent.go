package servicebus

import (
	"fmt"
	"github.com/1-bi/log-api"
	"github.com/1-bi/servicebus/etcd"
	"github.com/1-bi/servicebus/schema"
	"github.com/bwmarrin/snowflake"
	"github.com/coreos/etcd/clientv3"
	"strconv"
	"strings"
	"sync"
)

var waitgroup sync.WaitGroup

// Agent define service bus agent proxy
type Agent struct {
	conf *Config

	nodeGenerater *snowflake.Node

	etcdServOpt *etcd.EtcdServiceOperations
}

func (myself *Agent) Start() {

	node, err := snowflake.NewNode(myself.conf.nodeNum)

	if err != nil {
		logapi.GetLogger("start").Fatal(err.Error(), nil)
	} else {
		myself.nodeGenerater = node
	}

	// --- connect client ---
	var cli *clientv3.Client
	cli, err = clientv3.New(myself.conf._etcdConfig)

	if err != nil {
		structBean := logapi.NewStructBean()
		structBean.LogStringArray("etcd.server", myself.conf._etcdConfig.Endpoints)
		logapi.GetLogger("serviebus.Start").Fatal("Connect etcd server fail.", structBean)
		return
	}

	servOptsMap := make(map[string]string, 0)
	myself.etcdServOpt = etcd.NewEtcdServiceOperations(cli, servOptsMap)

	waitgroup.Add(2)
	// --- open thread
	go func() {
		go myself.startRegisterServer(cli)
		waitgroup.Done()
	}()

	go func() {
		myself.startWatchServer(cli)
		waitgroup.Done()
	}()
	// --- start watch server
	waitgroup.Wait()
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

	reqEvent.ReqId = myself.nodeGenerater.Generate().Int64()
	reqEvent.Name = eventName
	reqEvent.MsgBody = msgBody

	// --- sent msg body ---
	var reqMsg []byte

	reqMsg, err := reqEvent.Marshal()

	if err != nil {
		return err
	}

	// get minion runinng node

	nodes, err := myself.etcdServOpt.GetAllNodeIds("minion")
	if err != nil {
		return err
	}

	// pub the message to content
	for _, node := range nodes {

		// --- key ---
		var key = strings.Join([]string{"reqm", strconv.FormatInt(reqEvent.ReqId, 10), "mi=" + node}, "/")

		// --- set the key value ---

		err = myself.etcdServOpt.SetMessage(key, reqMsg)

		if err != nil {
			break
		}

	}

	// --- set use nats ---

	fmt.Println(nodes)

	return nil
}

// ---------------------  private method ---
func (myself *Agent) startRegisterServer(cli *clientv3.Client) {

	var err error

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

func (myself *Agent) startWatchServer(cli *clientv3.Client) {

	var err error
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
