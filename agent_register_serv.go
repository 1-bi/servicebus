package servicebus

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"log"
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
	Info    AgentInfo
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
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				log.Printf("Recv reply from service: %s, ttl:%d", s.nodeId, ka.TTL)
			}
		}
	}
}

func (s *AgentServiceRegService) Stop() {
	s.stop <- nil
}

func (s *AgentServiceRegService) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	info := &s.Info

	key := "/agent/nodes/" + s.nodeId
	value, _ := json.Marshal(info)

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
	log.Printf("servide:%s stop\n", s.nodeId)
	return err
}
