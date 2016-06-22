package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	readAgent()
}

func readAgent() {
	r := regexp.MustCompile("--([a-z]+)$")
	agentName := r.FindStringSubmatch(os.Args[1])[1]

	if agentName == "alice" {
		A := NewSocket("127.0.0.1:7000")
		B := NewSocket("172.19.0.1:5555")
		fmt.Printf("Agent %s started.\n", agentName)
		EstablishSocket(*A, *B)
	} else if agentName == "bob" {
		A := NewSocket("172.20.0.1:5555")
		B := NewSocket("127.0.0.1:8000")
		fmt.Printf("Agent %s started.\n", agentName)
		EstablishSocket(*A, *B)
	} else {
		fmt.Println("Error: Unknown agent. Abort.")
		os.Exit(-1)
	}
}