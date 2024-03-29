package metrics

import (
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"log"
	"time"
)

var peerNotifee = false

type PingStatus struct {
	UUID   string    `json:"uuid"`
	Source string    `json:"source"`
	Target string    `json:"target"`
	Alive  bool      `json:"alive"`
	RTT    int64     `json:"rtt"`
	Time   time.Time `json:"time"`
}

func NewPingStatus(source, target string) *PingStatus {
	return &PingStatus{
		Source: source,
		Target: target}
}

func (status *PingStatus) SetPingStatus(res ping.Result, deadline *int) {
	if res.Error == nil {
		status.Alive = true
		status.RTT = res.RTT.Milliseconds()
		*deadline = 0
		log.Println("pinged", status.Target, "in", res.RTT)
	} else {
		status.Alive = false
		status.RTT = 0
		log.Println("pinged", status.Target, "without success", res.Error)
		*deadline++
	}
	status.Time = time.Now()
	status.UUID = uuid.New().String()
}
