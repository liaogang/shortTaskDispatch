package main

import "root/cmd_server"

func main() {
	err := cmd_server.Work()
	if err != nil {
		panic(err)
	}
}
