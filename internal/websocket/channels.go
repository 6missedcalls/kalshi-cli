package websocket

import (
	"errors"
	"sync"
)

// CommandType represents a WebSocket command type
type CommandType string

const (
	CmdAuth        CommandType = "auth"
	CmdSubscribe   CommandType = "subscribe"
	CmdUnsubscribe CommandType = "unsubscribe"
	CmdPing        CommandType = "ping"
)

// AuthCommandID is the reserved command ID for authentication
const AuthCommandID = 1

// Command represents a WebSocket command to send to the server
type Command struct {
	ID     int                    `json:"id"`
	Cmd    CommandType            `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

// Subscription represents an active channel subscription
type Subscription struct {
	Channel Channel
	Params  map[string]string
}

// SubscriptionManager manages channel subscriptions
type SubscriptionManager struct {
	subscriptions map[Channel]Subscription
	nextID        int
	mu            sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[Channel]Subscription),
		nextID:        2, // Start at 2 since 1 is reserved for auth
	}
}

// Subscribe creates a subscription command for a channel
func (sm *SubscriptionManager) Subscribe(channel Channel, params map[string]string) (*Command, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := sm.nextID
	sm.nextID++

	cmdParams := map[string]interface{}{
		"channels": []Channel{channel},
	}

	// Add channel-specific params
	for k, v := range params {
		cmdParams[k] = v
	}

	sm.subscriptions[channel] = Subscription{
		Channel: channel,
		Params:  params,
	}

	return &Command{
		ID:     id,
		Cmd:    CmdSubscribe,
		Params: cmdParams,
	}, nil
}

// Unsubscribe creates an unsubscribe command for a channel
func (sm *SubscriptionManager) Unsubscribe(channel Channel) (*Command, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.subscriptions[channel]; !ok {
		return nil, errors.New("not subscribed to channel")
	}

	id := sm.nextID
	sm.nextID++

	delete(sm.subscriptions, channel)

	return &Command{
		ID:  id,
		Cmd: CmdUnsubscribe,
		Params: map[string]interface{}{
			"channels": []Channel{channel},
		},
	}, nil
}

// IsSubscribed checks if a channel is currently subscribed
func (sm *SubscriptionManager) IsSubscribed(channel Channel) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.subscriptions[channel]
	return ok
}

// GetSubscriptions returns all current subscriptions
func (sm *SubscriptionManager) GetSubscriptions() []Subscription {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	subs := make([]Subscription, 0, len(sm.subscriptions))
	for _, sub := range sm.subscriptions {
		subs = append(subs, sub)
	}
	return subs
}

// RestoreSubscriptions creates subscribe commands for restoring subscriptions after reconnect
func (sm *SubscriptionManager) RestoreSubscriptions(subs []Subscription) []*Command {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	commands := make([]*Command, 0, len(subs))

	for _, sub := range subs {
		id := sm.nextID
		sm.nextID++

		cmdParams := map[string]interface{}{
			"channels": []Channel{sub.Channel},
		}

		for k, v := range sub.Params {
			cmdParams[k] = v
		}

		sm.subscriptions[sub.Channel] = sub

		commands = append(commands, &Command{
			ID:     id,
			Cmd:    CmdSubscribe,
			Params: cmdParams,
		})
	}

	return commands
}

// Clear removes all subscriptions (used during disconnect)
func (sm *SubscriptionManager) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.subscriptions = make(map[Channel]Subscription)
}

// BuildAuthCommand creates an authentication command
func BuildAuthCommand(apiKeyID, signature, timestamp string) *Command {
	return &Command{
		ID:  AuthCommandID,
		Cmd: CmdAuth,
		Params: map[string]interface{}{
			"api_key":   apiKeyID,
			"signature": signature,
			"timestamp": timestamp,
		},
	}
}

// BuildPingCommand creates a ping command
func BuildPingCommand(id int) *Command {
	return &Command{
		ID:     id,
		Cmd:    CmdPing,
		Params: map[string]interface{}{},
	}
}
