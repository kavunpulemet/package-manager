package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Use: pm create packet.json / pm update packages.json")
		return
	}

	var (
		cmd = os.Args[1]
		arg = os.Args[2]
	)

	switch cmd {
	case "create":
		err := CreatePacket(arg)
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "update":
		err := UpdatePackages(arg)
		if err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
