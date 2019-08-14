package main

import (
	"fmt"
	"os"

	gameServer "../gameServer"
	gatewayServer "../gatewayServer"
)

func main() {
	go gatewayServer.Run()
	gameServer.Run()

	fmt.Println("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
