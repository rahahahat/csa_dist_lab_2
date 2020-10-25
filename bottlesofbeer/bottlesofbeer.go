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
	fmt.Print(bottles)
	fmt.Print(" bottles of beer on the wall, ")
	fmt.Print(bottles)
	fmt.Print(" bottles of beer.")
	fmt.Println(" Take one down, pass it around...")
}

func (s *BottlesOfBeer) Call(req Request, res *Response) (err error) {
	if req.Bottles != 0 {
		newBottles := req.Bottles - 1
		res.Bottles = newBottles
		request := Request{Bottles: newBottles}
		response := new(Response)
		client1, _ := rpc.Dial("tcp", *&nextAddr)
		defer client1.Close()
		print(newBottles)
		time.Sleep(1 * time.Second)
		client1.Call("BottlesOfBeer.Call", request, &response)
		return
	} else {
		listener.Close()
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
		request := Request{Bottles: *bottles}
		response := new(Response)
		client, _ := rpc.Dial("tcp", *&nextAddr)
		defer client.Close()
		print(*bottles)
		time.Sleep(1 * time.Second)
		client.Go(Call, request, &response, nil)
	}
	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer listener.Close()
	rpc.Accept(listener)
	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
}
