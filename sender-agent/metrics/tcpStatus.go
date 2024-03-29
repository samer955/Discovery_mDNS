package metrics

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

func NewTCPstatus(ip string) *TCPstatus {
	var host, _ = os.Hostname()
	return &TCPstatus{
		Ip:       ip,
		Hostname: host,
	}
}

var execCommand = exec.Command

//Working on Windows and Linux in order to get the number of open tcp queue of the node.
//Execution of the "netstat -na" Command in order to get all the ESTABLISHED Queue
func (t *TCPstatus) TcpQueueSize() {
	out, err := execCommand("netstat", "-na").Output()
	if err != nil {
		log.Println(err)
		return
	}
	output := string(out)
	tcpQueue, err := numberOfTcpQueue(output)
	if err != nil {
		log.Println(err)
		return
	}
	t.QueueSize = tcpQueue
}

//This function run the command netstat -s or netstat -st in order to get the number of
//TCP segments received and sent
func (t *TCPstatus) TcpSegmentsNumber() {

	if runtime.GOOS == "windows" {
		pr, err := execCommand("netstat", "-s").Output()
		if err != nil {
			log.Println(err)
			return
		}
		received, sent, err := numberOfSegmentsWindows(string(pr))

		if err != nil {
			log.Println(err)
		}
		t.Received = received
		t.Sent = sent
		return
	}

	if runtime.GOOS == "linux" {
		pr, err := execCommand("netstat", "-st").Output()
		if err != nil {
			log.Println(err)
			return
		}
		received, sent, err := numbersOfSegmentsLinux(string(pr))

		if err != nil {
			log.Println(err)
		}
		t.Received = received
		t.Sent = sent
		return
	}
}

//Format the output of "netstat -na" to find the ESTAB tcp queue
func numberOfTcpQueue(s string) (tcpConn int, err error) {

	var lines [][]string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if (strings.HasPrefix(words[0], "TCP") || strings.HasPrefix(words[0], "tcp")) &&
			strings.HasPrefix(words[len(words)-1], "ESTAB") {
			lines = append(lines, words)
		}
	}
	err = scanner.Err()
	return len(lines), err
}

func numbersOfSegmentsLinux(s string) (int, int, error) {

	var segmentsReceived = 0
	var segmentsSent = 0

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if len(words) == 3 {
			if strings.Contains(words[1], "segments") && strings.Contains(words[2], "received") {
				value, err := strconv.Atoi(words[0])
				if err != nil {
					fmt.Println(err)
				} else {
					segmentsReceived = value
				}
			}
		}
		if len(words) == 4 {
			if strings.Contains(words[1], "segments") && strings.Contains(words[2], "sent") {
				value, err := strconv.Atoi(words[0])
				if err != nil {
					fmt.Println(err)
				} else {
					segmentsSent = value
				}
			}
		}
	}
	err := scanner.Err()
	return segmentsReceived, segmentsSent, err

}

func numberOfSegmentsWindows(s string) (int, int, error) {

	var segmentsReceived = 0
	var segmentsSent = 0

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if strings.Contains(words[0], "Segments") && strings.Contains(words[1], "Sent") {
			value, err := strconv.Atoi(words[3])
			if err != nil {
				fmt.Println(err)
			} else {
				segmentsSent += value
			}
		}
		if strings.Contains(words[0], "Segments") && strings.Contains(words[1], "Received") {
			value, err := strconv.Atoi(words[3])
			if err != nil {
				fmt.Println(err)
			} else {
				segmentsReceived += value
			}
		}
	}
	err := scanner.Err()
	return segmentsReceived, segmentsSent, err
}
