package main

import (
	"bufio"
	"flag"
	"fmt"
	l "github.com/minio/minl/webrpc"
	"os"
	"time"
)

var (
	payloadFlag = flag.String("payload", "", "JSON as input for your Lambda function")
)

func onResult(msg []byte) {
	fmt.Println("Server result:", string(msg))
	Done <- 1
}

var Done chan int

func main() {
	flag.Parse()
	l.SetLoggingVerbosity(1)
	l.Mode = l.ModeInit
	l.OnData = onResult
	reader := bufio.NewReader(os.Stdin)

	Done = make(chan int, 1)

	// Initiate connection
	l.Start(true)

	// Input loop.
	for l.Mode != l.ModeChat {
		switch l.Mode {
		case l.ModeConnect:
			text, _ := reader.ReadString('\n')
			l.SignalReceive(text)

			// Sleep for short while to have mode changed
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("Invoking server with payload:", *payloadFlag)
	l.SendData(*payloadFlag)

	<-Done
}
