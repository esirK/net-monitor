package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/esirk/net-monitor/data_access"
	"github.com/esirk/net-monitor/types"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const target = "google.com"

func main() {
	ch := make(chan types.PingResult)

	var targ string

	flag.StringVar(&targ, "target", target, "Target to ping")
	flag.Parse()

	for {
		time.Sleep(time.Second * 3)
		go Ping(targ, ch)
		go data_access.SaveState(ch)
	}
}

func Ping(target string, ch chan types.PingResult) {
	ip, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		ch <- types.PingResult{State: 0, Ping_time: -1}
		return
	}
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		fmt.Println("Error on ListenPacket")
		panic(err)
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}
	msg_bytes, err := msg.Marshal(nil)
	if err != nil {
		ch <- types.PingResult{State: 0, Ping_time: -1}
		fmt.Printf("Error on Marshal %v\n", msg_bytes)
		return
	}

	start := time.Now()
	// Write the message to the listening connection
	if _, err := conn.WriteTo(msg_bytes, &net.UDPAddr{IP: net.ParseIP(ip.String())}); err != nil {
		fmt.Printf("Error on WriteTo %v", err)
		panic(err)
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	if err != nil {
		fmt.Printf("Error on SetReadDeadline %v\n", err)
		panic(err)
	}
	reply := make([]byte, 1500)
	n, _, err := conn.ReadFrom(reply)

	if err != nil {
		fmt.Printf("Got an Error on ReadFrom %v\n", err)
		ch <- types.PingResult{State: 0, Ping_time: -1}
		return
	}
	duration := time.Since(start)
	index := strings.Index(duration.String(), "ms")
	if index == -1{
		ch <- types.PingResult{State: 0, Ping_time: -1}
		return
	}

	x := duration.String()[:index]
	ping_time, err := strconv.ParseFloat(x, 64)
	if err != nil {
		fmt.Printf("Error on ParseFloat %v\n", err)
		ch <- types.PingResult{State: 0, Ping_time: ping_time}
		return
	}
	parsed_reply, err := icmp.ParseMessage(1, reply[:n])

	if err != nil {
		ch <- types.PingResult{State: 0, Ping_time: ping_time}
		fmt.Printf("Error on ParseMessage %v after %v\n", err, ping_time)
		return
	}

	switch parsed_reply.Code {
	case 0:
		// Got a reply so we can save this
		fmt.Printf("Got Reply from %s after %v\n", target, ping_time)
		ch <- types.PingResult{State: 1, Ping_time: ping_time}
	case 3:
		fmt.Printf("Host %s is unreachable after %v\n", target, ping_time)
		ch <- types.PingResult{State: 0, Ping_time: ping_time}
		// Given that we don't expect google to be unreachable, we can assume that our network is down
	case 11:
		ch <- types.PingResult{State: 0, Ping_time: ping_time}
		// Time Exceeded so we can assume our network is slow
		fmt.Printf("Host %s is slow %v\n", target, ping_time)
	default:
		ch <- types.PingResult{State: 0, Ping_time: ping_time}
		// We don't know what this is so we can assume it's unreachable
		fmt.Printf("Host %s is unreachable %v\n", target, ping_time)
	}
	return
}
