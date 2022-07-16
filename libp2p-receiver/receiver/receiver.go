package receiver

import (
	"context"
	"encoding/json"
	host2 "github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"libp2p-receiver/database"
	"libp2p-receiver/variables"
	"log"
	"time"
)

var db = database.Database()

func ReadTimeMessages(subscribe *pubsub.Subscription, context context.Context, node host2.Host) {
	for {
		message, err := subscribe.Next(context)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {
				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())

				peer := new(variables.PeerInfo)

				//Unmarshal the file into the peer struct
				json.Unmarshal(message.Data, peer)

				//Latency is calculated from the time when the peer send the message
				//and when the receiver reads it (in millis)

				latency := time.Since(peer.Time).Milliseconds()

				log.Println("latency :", latency)

				//Here we store latency of the peer in the database as well as node_id, ip_address
				//in order to identify it
				_, err = db.Exec("INSERT INTO latency(hostname,node_id,ip,latency,time) "+
					"VALUES($1,$2,$3,$4,$5)", peer.Hostname, peer.Id, peer.Ip, latency, time.Now())

				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func ReadRamInformation(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host) {
	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				ram := new(variables.Ram)

				//parse the JSON-encoded data and store the result into ram
				json.Unmarshal(message.Data, ram)

				//Here we store cpu usage percentage of the peer in the database as well
				//as node_id, ip_address to identify the peer
				_, err = db.Exec("INSERT INTO ram(hostname,node_id,ip,usage,time)"+
					" VALUES($1,$2,$3,$4,$5)",
					ram.Hostname,
					ram.Id,
					ram.Ip,
					ram.Usage,
					ram.Time)

				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}
	}
}

func ReadCpuInformation(subscribe *pubsub.Subscription, ctx context.Context, node host2.Host) {
	for {
		message, err := subscribe.Next(ctx)
		if err != nil {
			log.Println("cannot read from topic")
		} else {
			if message.ReceivedFrom.String() != node.ID().Pretty() {

				cpu := new(variables.Cpu)

				//parse the JSON-encoded data and store the result into cpu
				json.Unmarshal(message.Data, cpu)

				//Here we store cpu usage percentage of the peer in the database as well
				//as node_id, ip_address to identify the peer
				_, err = db.Exec("INSERT INTO cpu(hostname,node_id,ip,usage,time)"+
					" VALUES($1,$2,$3,$4,$5)",
					cpu.Hostname,
					cpu.Id,
					cpu.Ip,
					cpu.Usage,
					cpu.Time)

				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Message: <%s> %s", message.Data, message.ReceivedFrom.String())
			}
		}
	}
}

func latencyCalculate(actual, source int64) int64 {
	return actual - source
}
