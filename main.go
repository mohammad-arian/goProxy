package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var logger, _ = os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

func main() {
	defer logger.Close()
	log.SetOutput(logger)

	handle := http.HandlerFunc(handler)
	log.Fatal(http.ListenAndServe("192.168.1.209:555", handle))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ok := firewall(w, r)

	if ok {
		c := context.TODO()
		timeOutContext, cancel := context.WithTimeout(c, 40*time.Second)
		if r.Method == http.MethodConnect {
			go handleConnect(w, r, timeOutContext)
		} else {
			go handleHTTP(w, r, timeOutContext)
		}

		select {
		case <-timeOutContext.Done():
			cancel()
			log.Println("Connection was canceled. Host was: " + r.Host)
		}
	}
}

func firewall(w http.ResponseWriter, r *http.Request) bool {
	if strings.Contains(r.Host, "digikala") {
		fmt.Println("digikala was found")

		desConn, err := net.Dial("tcp", "127.0.0.1:9000")

		w.WriteHeader(http.StatusOK)
		hijacker, _ := w.(http.Hijacker)

		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			log.Println(w, err.Error(), http.StatusServiceUnavailable)
		}
		go transfer(desConn, clientConn)
		go transfer(clientConn, desConn)
		time.Sleep(4000)
		fmt.Println("yess")
		return false
	}
	return true
}

func handleConnect(w http.ResponseWriter, r *http.Request, c context.Context) {
	desConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println(err, "error dialing: "+r.Host, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		log.Println(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		log.Println(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(desConn, clientConn)
	go transfer(clientConn, desConn)
}

func handleHTTP(w http.ResponseWriter, r *http.Request, c context.Context) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
