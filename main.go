package main

import "IM-System/Server"

func main() {
	server := Server.NewServer("127.0.0.1", 8888)
	server.Start()
}

//curl.exe --http0.9 127.0.0.1:8888
