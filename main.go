package main

func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}

//curl.exe --http0.9 127.0.0.1:8888
//更改运行时编译设置编译整个包，可避免undefined问题
