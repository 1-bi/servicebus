package servicebus

// Agent define service bus agent proxy
type Agent struct {
}

func (myself *Agent) Start() {

}

func (myself *Agent) Stop() {

}

func (myself *Agent) startRegisterServer() {

}

func (myself *Agent) startWatchServer() {

}

func (myself *Agent) checkRegCenterConnect() {

}

func (myself *Agent) getNodeAgentList() {

}

// NewAgent check agent
func NewAgent(conf *Config) *Agent {

	var agent = new(Agent)

	//  start scheduler

	// check health

	// start describe queue to nats
	return agent

}
