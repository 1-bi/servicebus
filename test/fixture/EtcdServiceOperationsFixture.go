package fixture

import (
	"fmt"
	"github.com/1-bi/servicebus/etcd"
	"github.com/1-bi/servicebus/schema"
	"github.com/coreos/etcd/clientv3"
	"github.com/gogo/protobuf/proto"
	"github.com/smartystreets/gunit"
	"log"
	"time"
)

// EtcdServiceOperationsFixture Test structure framework define
type EtcdServiceOperationsFixture struct {
	*gunit.Fixture

	servOper *etcd.EtcdServiceOperations
}

// SetupAgent
func (myself *EtcdServiceOperationsFixture) Setup() {
	// --- config

	var cli *clientv3.Client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		log.Println(err)
	}

	// --- create agent ---
	myself.servOper = etcd.NewEtcdServiceOperations(cli, nil)
}

func (myself *EtcdServiceOperationsFixture) Teardown() {

}

func (myself *EtcdServiceOperationsFixture) Test_GetMessage() {

	req, err := myself.servOper.GetMesssage("reqm/1129925668075737088")

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

}
