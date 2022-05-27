package utils

import (
	"os"
	"fmt"
	"os/signal"
	stanAPI "github.com/nats-io/stan.go"
)

// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
func CleanupAfterSIGINT(cleanupDone chan bool, sub stanAPI.Subscription, sc stanAPI.Conn) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nUnsubscribing and closing connection...\n")
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
}