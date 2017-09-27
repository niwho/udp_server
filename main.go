package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"

	"udp_server/logs"
	"udp_server/message"
)

func main() {
	var port int
	var ignore string
	flag.IntVar(&port, "port", 5678, "")
	flag.StringVar(&ignore, "conf", "", "")
	flag.StringVar(&ignore, "rpc", "", "")
	flag.StringVar(&ignore, "log", "", "")
	flag.StringVar(&ignore, "svc", "", "")
	flag.Parse()
	logs.InitLog("udp_ss.log", logs.L(5))
	af := NewAsyncFrame(1024, handleMsg)
	_ = af
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	fmt.Printf("pc=%+v, err=%+v", pc, err)
	for {
		buffer := make([]byte, 102400)
		n, addr, err := pc.ReadFrom(buffer)
		fmt.Println(addr)
		if err == nil {
			af.Write(buffer[:n])
		} else {
			fmt.Println("err", err)
		}
	}

}

func handleMsg(b []byte) {
	var am message.ApplicationMessage
	err := json.Unmarshal(b, &am)
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("handleMsg")
		return
	}
	switch am.MType {
	case "dingding":
		message.SendDD(am.Title, am.Content, am.At)
	case "kafka":
		break
	case "influx":
		break
	default:
		logs.Log(logs.F{"am": am}).Error("unkonw type!!!!")
	}

}
