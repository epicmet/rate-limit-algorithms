package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// cli [--cooldown <seconds>] url ...
	cooldown := flag.Int("cooldown", 1, "how many seconds it takes between each request")
	flag.Parse()

	fmt.Printf("Interval: %v\n", *cooldown)

	urls := flag.Args()

	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)

		go func() {
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(*cooldown) * time.Second)
			defer ticker.Stop()

			sigintChan := make(chan os.Signal, 1)
			signal.Notify(sigintChan, syscall.SIGINT, syscall.SIGTERM)

			for {
				select {
				case _ = <-sigintChan:
					{
						fmt.Println("Received interupt signal...")
						return
					}
				case _ = <-ticker.C:
					{
						resp, err := http.Get(u)
						if err != nil {
							log.Fatalf("Couldn't GET %v\n", u)
						}

						log.Printf("Called %v got status -> %v\n", u, resp.StatusCode)
					}
				}
			}
		}()
	}

	wg.Wait()
}
