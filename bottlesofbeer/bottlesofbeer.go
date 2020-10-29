package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
	//	"net"
)

var listener net.Listener
var count = 0
var nextAddr string

var Call = "BottlesOfBeer.Call"

type Response struct {
	Bottles int
}

type Request struct {
	Bottles int
}

type BottlesOfBeer struct{}

func print(bottles int) {
	if bottles != 0 {
		fmt.Print(bottles)
		fmt.Print(" bottles of beer on the wall, ")
		fmt.Print(bottles)
		fmt.Print(" bottles of beer.")
		fmt.Println(" Take one down, pass it around...")
	} else {
		fmt.Println("No more bottles of beer on the wall, no more bottles of beer.")
		fmt.Println("There's nothing else to fall, because there's no more bottles of beer on the wall.")
	}
}

func (s *BottlesOfBeer) Call(req Request, res *Response) (err error) {
	if req.Bottles > 0 {
		newBottles := req.Bottles - 1
		res.Bottles = newBottles
		request := Request{Bottles: newBottles}
		response := new(Response)
		client1, _ := rpc.Dial("tcp", *&nextAddr)
		defer client1.Close()
		print(req.Bottles)
		time.Sleep(1 * time.Second)
		client1.Go("BottlesOfBeer.Call", request, response, nil)
		return
	} else if req.Bottles == 0 {
		print(0)
		request := Request{Bottles: -1}
		response := new(Response)
		client1, _ := rpc.Dial("tcp", *&nextAddr)
		defer client1.Close()
		client1.Go("BottlesOfBeer.Call", request, response, nil)
		count = -1
		return
	} else {
		if count != -1 {
			request := Request{Bottles: -1}
			response := new(Response)
			client1, _ := rpc.Dial("tcp", *&nextAddr)
			defer client1.Close()
			client1.Go("BottlesOfBeer.Call", request, response, nil)
			defer listener.Close()
		} else {
			defer listener.Close()
		}
		return
	}
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "127.0.0.1:8030", "IP:Port string for next member of the round.")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()
	rpc.Register(&BottlesOfBeer{})
	if *bottles != 0 {
		request := Request{Bottles: *bottles - 1}
		response := new(Response)
		client, _ := rpc.Dial("tcp", *&nextAddr)
		defer client.Close()
		print(*bottles)
		time.Sleep(1 * time.Second)
		client.Go(Call, request, response, nil)
	}
	listener, _ = net.Listen("tcp", ":"+*thisPort)
	rpc.Accept(listener)
	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
}
