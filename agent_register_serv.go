package servicebus

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/1-bi/log-api"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

//the detail of service
type AgentInfo struct {
	LastUpdatedTime int64
}

func (myself *AgentInfo) SetLastUpdatedTime(t int64) {
	myself.LastUpdatedTime = t
}

func NewAgentInfo() *AgentInfo {
	var agentInfo = new(AgentInfo)

	return agentInfo
}

type AgentServiceRegService struct {
	nodeId  string
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
	_prefix string
}

func NewAgentRegisterService(nodeId string, cli *clientv3.Client) *AgentServiceRegService {

	// --- create  AgentService ---
	var agentRegServ = new(AgentServiceRegService)
	agentRegServ.client = cli
	agentRegServ.nodeId = nodeId
	agentRegServ.stop = make(chan error)

	agentRegServ._prefix = "/agent/nodes/"

	return agentRegServ
}

func (s *AgentServiceRegService) Start() error {

	ch, err := s.keepAlive()
	if err != nil {
		log.Fatal(err)
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				logapi.GetLogger("servicebus.AgentServiceRegService.start").Info("keep alive channel closed.", nil)
				//log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {

				// ---  update status ---

				structBean := logapi.NewStructBean()
				structBean.LogString("nodeId", s.nodeId)
				structBean.LogInt64("ttl time ", ka.TTL)

				logapi.GetLogger("servicebus.AgentServiceRegService.start").Debug("Recv reply from service: %s, ttl:%d", structBean)
				//log.Printf("Recv reply from service: %s, ttl:%d", s.nodeId, ka.TTL)
			}
		}
	}
}

func (s *AgentServiceRegService) Stop() {
	s.stop <- nil
}

func (s *AgentServiceRegService) updateAgentInfo() {

	// --- get the lastupdate register service ---
	info := NewAgentInfo()
	info.SetLastUpdatedTime(time.Now().UnixNano())

	//info := &s.Info
	value, _ := json.Marshal(info)

	// --- get properties key --
	key := s._prefix + s.nodeId

	_, err := s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	s.leaseid = resp.ID

}

func (s *AgentServiceRegService) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	// --- get the lastupdate register service ---
	info := NewAgentInfo()
	info.SetLastUpdatedTime(time.Now().UnixNano())

	//info := &s.Info
	value, _ := json.Marshal(info)

	// --- get properties key --
	key := s._prefix + s.nodeId

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *AgentServiceRegService) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}

	structBean := logapi.NewStructBean()
	structBean.LogString("nodeId", s.nodeId)
	logapi.GetLogger("servicebus.AgentServiceRegService.revoke").Info("servide:%s stop\n", structBean)

	//log.Printf("servide:%s stop\n", s.nodeId)
	return err
}
