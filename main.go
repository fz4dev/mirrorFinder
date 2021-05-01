package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fzuchows/mirrors"
)

// https://beta.packtpub.com/reader/book/web_development/9781838643577/1/ch01lvl1sec07
type response struct {
	FastestURL string        `json:"fastest_url"`
	Latency    time.Duration `json:"latency"`
}

/**
*
 */
func findFastest(urls []string) response {
	urlChan := make(chan string)
	latencyChan := make(chan time.Duration)

	for _, url := range urls {
		mirrorUrl := url
		go func() {
			start := time.Now()
			_, err := http.Get(mirrorUrl + "/README")
			latency := time.Now().Sub(start) / time.Millisecond

			if err == nil {
				urlChan <- mirrorUrl
				latencyChan <- latency
			}
		}()
	}
	return response(<-urlChan, <-latencyChan)
}

/**
* MAIN
 */
func main() {
	http.HandleFunc("/fastest-mirror", func(w http.ResponseWriter, r *http.Request) {
		response := findFastest(mirrors.MirrorList)
		respJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.write(respJSON)
	})

	port := ":9000"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Starting server on port %sn", port)
	log.Fatal(server.ListenAndServe())
}
