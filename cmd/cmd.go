package main

import (
	"context"
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

	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigintChan
		fmt.Println("Received interupt signal...")
		cancel()
	}()

	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)

		go func() {
			defer wg.Done()
			pollURL(ctx, u, *cooldown)
		}()
	}

	wg.Wait()
}

func pollURL(ctx context.Context, url string, cooldown int) {
	ticker := time.NewTicker(time.Duration(cooldown) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			{
				fmt.Printf("Stopping polling for %s\n", url)
				return
			}
		case _ = <-ticker.C:
			{
				resp, err := http.Get(url)
				if err != nil {
					log.Fatalf("Couldn't GET %v\n", url)
				}

				log.Printf("Called %v got status -> %v\n", url, resp.StatusCode)
			}
		}
	}
}
