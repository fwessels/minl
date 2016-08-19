package main

import (
	"bufio"
	"fmt"
	l "github.com/minio/minl/webrpc"
	"os"
)

func onInvoke(msg []byte) {
	fmt.Println("serve invoked:", string(msg))

	// Invoke function and get result
	cmd := "./sandbox"
	cmd += " node myfunc.js input=`" + string(msg) + "`"
	fmt.Println("Simulating:", cmd)

	result := `{"result":"response"}`
	l.SendData(result)
}

func main() {
	l.SetLoggingVerbosity(1)
	l.Mode = l.ModeInit
	l.OnData = onInvoke

	reader := bufio.NewReader(os.Stdin)

	wait := make(chan int, 1)

	fmt.Println("\n ---- Please paste offer from peer ---- \n")

	// Input loop.
	for {
		text, _ := reader.ReadString('\n')
		switch l.Mode {
		case l.ModeInit:
			l.SignalReceive(text)
		}
		text = ""
	}

	<-wait
	fmt.Println("done")
}
