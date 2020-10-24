package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"secretstrings/stubs"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("wordlist")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	sliceData := strings.Split(string(data), "\n")
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	for _, word := range sliceData {
		request := stubs.Request{Message: string(word)}
		response := new(stubs.Response)
		client.Call(stubs.PremiumReverseHandler, request, response)
		fmt.Println("Responded :" + response.Message)
	} //TODO: connect to the RPC server and send the request(s)
}
