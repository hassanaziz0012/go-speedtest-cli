package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const fileUrl = "https://releases.ubuntu.com/24.04.2/ubuntu-24.04.2-desktop-amd64.iso"

// const fileUrl = "https://images.unsplash.com/photo-1720884413532-59289875c3e1?q=80&w=735&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"

func main() {
	res, err := http.Get(fileUrl)
	if err != nil {
		log.Fatal("failed to request file")
	}
	defer res.Body.Close()

	done := make(chan bool)

	go func() {
		defer close(done)
		start := time.Now()
		var bytesRead int
		ticker := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-ticker.C:
				timeTaken := time.Since(start).Seconds()

				speedInBytes := float64(bytesRead) / timeTaken
				speed := speedInBytes / 1024 / 1024 // MBs
				fmt.Printf("Download speed: %.2fMB/s\n(%f) / (%f)\n", speed, speedInBytes, timeTaken)
				start = time.Now()
				bytesRead = 0
			case <-done:
				return
			default:
				buf := make([]byte, 1024*256) // 256KB buffer

				n, err := res.Body.Read(buf)
				bytesRead += n
				if err == io.EOF {
					done <- true
					return
				}
				if err != nil {
					fmt.Println(err)
				}

			}
		}

	}()

	<-done
	fmt.Println("Finished.")
}
