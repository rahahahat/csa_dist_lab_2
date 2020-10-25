package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
	//	"net"
)

var nextAddr string

var Call = "BottlesOfBeer.Call"

type Response struct {
	Bottles int
}

type Request struct {
	Bottles int
}

type BottlesOfBeer struct{}

func (s *BottlesOfBeer) Call(req Request, res *Response) (err error) {
	if req.Bottles != 0 {
		newBottles := req.Bottles - 1
		res.Bottles = newBottles
		request := Request{Bottles: newBottles}
		response := new(Response)
		client, _ := rpc.Dial("tcp", nextAddr)
		fmt.Println(string(newBottles) + " bottles of beer on the wall." + string(newBottles) + " Bottles of beer. Take one down and pass it around")
		time.Sleep(2 * time.Second)
		client.Call(Call, request, response)
		return
	} else {
		return
	}
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()
	rpc.Register(&BottlesOfBeer{})
	if *bottles != 0 {
		fmt.Println(*bottles)
		request := Request{Bottles: *bottles}
		response := new(Response)
		client, _ := rpc.Dial("tcp", *&nextAddr)
		fmt.Println(string(*bottles) + " bottles of beer on the wall." + string(*bottles) + " Bottles of beer. Take one down and pass it around")
		time.Sleep(2 * time.Second)
		client.Call(Call, request, response)
	}
	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer listener.Close()
	rpc.Accept(listener)
	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
}
