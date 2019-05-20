package servicebus

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

type QueueWatcher struct {
	client   *clientv3.Client
	eventKey string
	cbs      []Callback
}

func NewQueueWatcher(client *clientv3.Client) *QueueWatcher {
	watcher := new(QueueWatcher)
	watcher.client = client
	return watcher
}

func (myself *QueueWatcher) SetEventKey(eventKey string) {
	myself.eventKey = eventKey
}

func (myself *QueueWatcher) SetCallbacks(newCbs ...Callback) {
	myself.cbs = newCbs
}

func (myself *QueueWatcher) run() error {
	myself.watchNodeChange()

	return nil
}

func (myself *QueueWatcher) watchNodeChange() {

	var stopChan chan struct{}
	ctx, cancel := context.WithCancel(context.Background())
	cancelRoutine := make(chan struct{})
	defer close(cancelRoutine)

	go func() {

		for {
			select {
			// 当外部传来 stopChan 时， cancel watcher 的 context
			case <-stopChan:
				fmt.Println("0000 ---")
				cancel()
				//return
			case <-cancelRoutine:
				fmt.Println("another ------------")
				//return
			}

			fmt.Println("-------------- anotieo ----")

		}

	}()

	// --- watch message of node changed
	rch := myself.client.Watch(ctx, myself.eventKey)

	for wresp := range rch {

		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%fixture] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(ev)

				fmt.Println(info)

			case clientv3.EventTypeDelete:
				fmt.Printf("[%fixture] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				fmt.Println(string(ev.Kv.Key))
			}
		}
	}

}
