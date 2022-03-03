package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/esirk/net-monitor/data_access"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const target = "google.com"

func main() {
	for {
		time.Sleep(time.Second * 3)
		state, ping_time := Ping(target)
		data_access.SaveState(state, ping_time)
	}
}

func Ping(target string) (state int, ping_time float64) {
	ip, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		panic(err)
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
		state = 0
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
		state = 0
		return
	}
	duration := time.Since(start)
	x := duration.String()[:strings.Index(duration.String(), "ms")]
	ping_time, err = strconv.ParseFloat(x, 64)
	if err != nil {
		fmt.Printf("Error on ParseFloat %v\n", err)
		state = 0
		return
	}
	parsed_reply, err := icmp.ParseMessage(1, reply[:n])

	if err != nil {
		state = 0
		fmt.Printf("Error on ParseMessage %v after %v\n", err, ping_time)
		return
	}

	switch parsed_reply.Code {
	case 0:
		// Got a reply so we can save this
		fmt.Printf("Got Reply from %s after %v\n", target, ping_time)
		state = 1
	case 3:
		fmt.Printf("Host %s is unreachable after %v\n", target, ping_time)
		state = 0
		// Given that we don't expect google to be unreachable, we can assume that our network is down
	case 11:
		state = 0
		// Time Exceeded so we can assume our network is slow
		fmt.Printf("Host %s is slow %v\n", target, ping_time)
	default:
		state = 0
		// We don't know what this is so we can assume it's unreachable
		fmt.Printf("Host %s is unreachable %v\n", target, ping_time)
	}
	return
}
