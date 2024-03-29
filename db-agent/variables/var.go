package variables

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

type SystemInfo struct {
	Id         peer.ID   `json:"node_id"`
	UUID       string    `json:"uuid"`
	Ip         string    `json:"ip"`
	Hostname   string    `json:"hostname"`
	OS         string    `json:"os"`
	Platform   string    `json:"platform"`
	Version    string    `json:"version"`
	Role       string    `json:"role"`
	Network    string    `json:"network"`
	Time       time.Time `json:"time"`
	OnlineUser int       `json:"online_user"`
	Latency    int64     `json:"latency"`
}

type Cpu struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"node"`
	Model    string    `json:"model"`
	Usage    int       `json:"usage"`
	Time     time.Time `json:"time"`
}

type Ram struct {
	Id       string    `json:"node_id"`
	UUID     string    `json:"uuid"`
	Ip       string    `json:"ip"`
	Hostname string    `json:"node,omitempty"`
	Usage    int       `json:"usage, omitempty"`
	Time     time.Time `json:"time, omitempty"`
}

type PingStatus struct {
	UUID     string    `json:"uuid"`
	Source   string    `json:"source"`
	Target   string    `json:"target"`
	Alive    bool      `json:"alive"`
	RTT      int64     `json:"rtt"`
	SourceIp string    `json:"source_ip"`
	TargetIp string    `json:"target_ip"`
	Time     time.Time `json:"time"`
}

//incoming and outgoing data transferred by the local peer.
type Bandwidth struct {
	UUID     string    `json:"uuid"`
	Id       string    `json:"id"`
	Source   string    `json:"source"`
	Target   string    `json:"target"`
	TotalIn  int64     `json:"total_in"`
	TotalOut int64     `json:"total_out"`
	RateIn   int       `json:"rate_in"`
	RateOut  int       `json:"rate_out"`
	Time     time.Time `json:"time"`
}

//Queue Size = number of open TCP connections
//Received = number of segments received
//Sent = number of segments sent
type TCPstatus struct {
	UUID      string    `json:"uuid"`
	Hostname  string    `json:"hostname"`
	Ip        string    `json:"ip"`
	QueueSize int       `json:"tcp_queue_size"`
	Received  int       `json:"segments_received"`
	Sent      int       `json:"segments_sent"`
	Time      time.Time `json:"time"`
}
