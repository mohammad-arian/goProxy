package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var logFile, _ = os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
var blockList = readBlockList()

func main() {
	defer logFile.Close()
	log.SetOutput(logFile)

	handle := http.HandlerFunc(handler)
	log.Fatal(http.ListenAndServe("192.168.1.209:555", handle))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ok := firewall(w, r)

	if ok {
		if r.Method == http.MethodConnect {
			handleConnect(w, r)
		} else {
			handleHTTP(w, r)
		}
	}
}

func firewall(w http.ResponseWriter, r *http.Request) bool {
	for _, j := range blockList {
		if strings.Contains(r.Host, j) {
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				log.Println(w, "Hijacking not supported", http.StatusInternalServerError)
			}
			clientConn, _, _ := hijacker.Hijack()
			_, err := clientConn.Write([]byte("HTTP/1.1 200 OK \n\n <h1 dir=\"rtl\" style=\"font-size:9vw\">&#1575;&#1587;&#1578;&#1594;&#1601;&#1585;&#1575;&#1604;&#1604;&#1607;&#1563; &#1576;&#1585;&#1608; &#1578;&#1608;&#1576;&#1607; &#1705;&#1606;&#1563; &#1576;&#1583;&#1608;</h1>"))
			if err != nil {
				log.Println(err)
			}
			clientConn.Close()
			return false
		}
	}
	return true
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	desConn, err := net.DialTimeout("tcp", r.Host, 30*time.Second)
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
		return
	}
	go transfer(desConn, clientConn)
	go transfer(clientConn, desConn)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
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

func readBlockList() []string {
	var list []string
	file, err := os.Open("BlockList.txt")
	defer file.Close()
	if err != nil {
		return list
	}
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		list = append(list, s.Text())
	}
	return list
}
