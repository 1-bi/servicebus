package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"strings"
)

type EtcdServiceOperations struct {
	client *clientv3.Client

	_nodePrefix string
}

func NewEtcdServiceOperations(client *clientv3.Client, opts map[string]string) *EtcdServiceOperations {

	var etcdBiz = new(EtcdServiceOperations)

	if opts["nodeprefix"] != "" {
		etcdBiz._nodePrefix = opts["nodeprefix"]
	} else {
		etcdBiz._nodePrefix = "nodes/"
	}

	etcdBiz.client = client

	return etcdBiz

}

func (myself EtcdServiceOperations) GetAllNodeIds(role string) ([]string, error) {

	var prefix = myself._nodePrefix + role + "="

	// --- watch message of node changed

	resp, err := myself.client.Get(context.Background(), prefix, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	nodeIds := make([]string, 0)
	var nodeId string

	for _, value := range resp.Kvs {

		nodeId = strings.Replace(string(value.Key), prefix, "", 0)

		nodeIds = append(nodeIds, nodeId)
	}

	return nodeIds, nil
}
