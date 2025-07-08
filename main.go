package main

func main() {
	server := NewServer()
	defer server.Close()

	server.Run()
}