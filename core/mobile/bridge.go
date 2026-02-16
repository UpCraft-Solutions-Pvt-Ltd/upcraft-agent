package mobile

import "github.com/UpCraft-Solutions-Pvt-Ltd/upcraft-agent/core/engine"

// UpCraftBridge is a gomobile-exported entrypoint that accepts simple string I/O.
type UpCraftBridge struct {
	agent *engine.Agent
}

func NewBridge() *UpCraftBridge {
	return &UpCraftBridge{
		agent: engine.NewAgent(),
	}
}

func (b *UpCraftBridge) Start() {
	go b.agent.Start()
}

// ProcessScreenEvent receives screen-tree JSON and returns JSON action commands.
func (b *UpCraftBridge) ProcessScreenEvent(inputJSON string) string {
	return b.agent.HandleScreenInput(inputJSON)
}
