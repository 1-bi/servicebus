package servicebus

import (
	"context"
	"fmt"
	"github.com/1-bi/log-api"
	"github.com/1-bi/servicebus/etcd"
	"github.com/1-bi/servicebus/schema"
	"github.com/bwmarrin/snowflake"
	"github.com/coreos/etcd/clientv3"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/stan.go"
	"strconv"
	"strings"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup

// Agent define service bus agent proxy
type Agent struct {
	conf *Config

	natsConn stan.Conn

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

	natsServer := strings.Join(myself.conf._natsHost, ",")
	myself.natsConn, err = stan.Connect("test-cluster", "clienttest", stan.NatsURL(natsServer))
	if err != nil {
		structBean := logapi.NewStructBean()
		structBean.LogStringArray("nats.server", myself.conf._natsHost)
		logapi.GetLogger("serviebus.Start").Fatal("Connect nats server fail.", structBean)
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

	// open and connect nats subscribe queue message

	go func() {
		myself.openNatsSubscribe(myself.natsConn)
	}()

	// --- start watch server
	waitgroup.Wait()
}

func (myself *Agent) Stop() {

}

// On implement event name
func (myself *Agent) On(eventName string, fn func(ReqMsgContext)) error {

	// --- send message to  nats ---
	natsServer := strings.Join(myself.conf._natsHost, ",")
	_, err := stan.Connect("serv-clusterId", "clienttest", stan.NatsURL(natsServer))
	if err != nil {
		return err
	}

	return nil
}

// FireByQueue call by event name and define callback
func (myself *Agent) FireByQueue(eventName string, msgBody []byte, callback ...Callback) error {

	// --- send message to  nats ---

	// serialization runtimeArgs
	reqEvent := new(schema.ReqEvent)

	reqEvent.ReqId = myself.nodeGenerater.Generate().Int64()
	reqEvent.Name = eventName
	reqEvent.MsgBody = msgBody

	// --- sent msg body ---
	var reqMsg []byte

	reqMsg, err := proto.Marshal(reqEvent)

	if err != nil {
		return err
	}

	reqQ := new(schema.ReqQ)
	reqQ.ReqId = reqEvent.ReqId
	reqQ.Name = reqEvent.Name
	reqQ.ComType = schema.ReqQ_QUE

	var req []byte
	req, err = proto.Marshal(reqQ)

	if err != nil {
		return err
	}

	// get minion runinng node

	// --- key ---
	var key = strings.Join([]string{"reqm", strconv.FormatInt(reqEvent.ReqId, 10)}, "/")

	// --- set the key value ---
	err = myself.etcdServOpt.SetMessage(key, reqMsg)

	myself.natsConn.Publish("reqm", req)

	return nil
}

func (myself *Agent) FireByPublish(eventName string, msgBody []byte, callback ...Callback) error {

	// --- send message to  nats ---

	// serialization runtimeArgs
	reqEvent := new(schema.ReqEvent)

	reqEvent.ReqId = myself.nodeGenerater.Generate().Int64()
	reqEvent.Name = eventName
	reqEvent.MsgBody = msgBody

	// --- sent msg body ---+
	var reqMsg []byte

	reqMsg, err := proto.Marshal(reqEvent)

	if err != nil {
		return err
	}

	reqQ := new(schema.ReqQ)
	reqQ.ReqId = reqEvent.ReqId
	reqQ.Name = reqEvent.Name
	reqQ.ComType = schema.ReqQ_SUB

	var req []byte
	req, err = proto.Marshal(reqQ)

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

	myself.natsConn.Publish("reqm", req)

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

func (myself *Agent) openNatsSubscribe(conn stan.Conn) {
	var cli *clientv3.Client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		fmt.Println("okkdf")
	}

	_, _ = conn.Subscribe("reqm", func(m *stan.Msg) {

		reqQ := new(schema.ReqQ)

		if err := proto.Unmarshal(m.Data, reqQ); err != nil {
			fmt.Println(err)
		}

		// --- get msg body from etcd cache --
		var key = strings.Join([]string{"reqm", strconv.FormatInt(reqQ.ReqId, 10)}, "/")

		fmt.Println("----------99- ")

		resp, err := cli.Get(context.Background(), key)

		if err != nil {
			fmt.Println("00")
		}

		fmt.Println("counter message")
		fmt.Println(resp.Count)

		// --- req message
		/*
			req, err := myself.etcdServOpt.GetMesssage(key)

			if err != nil {
				fmt.Println(err)
			}

			// 解码
			unmaReqEvent := new(schema.ReqEvent)
			if err := proto.Unmarshal(req, unmaReqEvent); err != nil {
				log.Fatal("failed to unmarshal: ", err)
			}

			fmt.Println(unmaReqEvent.ReqId)
			fmt.Println(unmaReqEvent.Name)
			fmt.Println(string(unmaReqEvent.MsgBody))
		*/

	})

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
