package subscriber

import (
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"
	"sender-agent/node"
	"testing"
)

func setupEnvironment(t *testing.T) (node.Node, context.Context) {
	ctx := context.Background()
	host, _ := libp2p.New()
	node := node.Node{Host: host}

	t.Cleanup(func() {
		node.Host.Close()
		ctx.Done()
	})
	return node, ctx
}

func TestNewPubSubService(t *testing.T) {
	node, ctx := setupEnvironment(t)
	psub := NewPubSubService(ctx, node)

	assert.NotEqual(t, psub, nil)
}

func TestJoinTopic(t *testing.T) {
	const roomtest = "test"
	node, ctx := setupEnvironment(t)
	psub := NewPubSubService(ctx, node)

	testTopic := psub.JoinTopic(roomtest)

	assert.Equal(t, testTopic.String(), roomtest)
}

func TestSubscribe(t *testing.T) {
	const roomtest = "test"
	node, ctx := setupEnvironment(t)
	ps := NewPubSubService(ctx, node)
	testTopic := ps.JoinTopic(roomtest)

	ps.Subscribe(testTopic)

	assert.Contains(t, ps.psub.GetTopics(), testTopic.String())
}

func TestPublish(t *testing.T) {

	type Message struct {
		Data string
	}
	helloMessage := new(Message)
	helloMessage.Data = "Hello World"
	const roomtest = "test"
	node, ctx := setupEnvironment(t)
	ps := NewPubSubService(ctx, node)
	testTopic := ps.JoinTopic(roomtest)
	subscr := ps.Subscribe(testTopic)

	Publish(helloMessage, ctx, testTopic)
	//read the message published
	receivedMessageByte, _ := subscr.Next(ctx)
	//unmarshal message data
	receivedMess := new(Message)
	json.Unmarshal(receivedMessageByte.Data, receivedMess)

	assert.Equal(t, receivedMess.Data, helloMessage.Data)
}
