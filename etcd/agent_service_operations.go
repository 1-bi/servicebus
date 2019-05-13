package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
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

	fmt.Println(prefix)

	// --- watch message of node changed
	resp, err := myself.client.Get(context.Background(), prefix, clientv3.WithPrevKV())

	if err != nil {
		return nil, err
	}

	nodeIds := make([]string, 0)

	fmt.Println(resp.Count)

	for key, value := range resp.Kvs {

		fmt.Println(key)

		nodeIds = append(nodeIds, string(value.Key))
	}

	return nodeIds, nil
}
