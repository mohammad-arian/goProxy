# goProxy
goProxy implements a simple HTTP proxy with a built-in firewall to block certain domains from being accessed. It intercepts HTTP and HTTPS requests, applies the firewall rules, and forwards the allowed requests to the destination server.


Usage
---
1. Set up the Block List: Create a file named "BlockList.txt" in the same directory as the code. List the domain names (one per line) that you want to block access to in this file.
2. Build and Run: Execute the following commands in the terminal:
```bash
$ go build proxy.go
$ ./proxy
```
3. write the ip and port you want the server to listen to in this format: ip:port

Logging
---
All incoming requests and blocked domains are logged in "log.txt". Make sure the program has write permissions to create and update the log file.

Disclaimer
---
The code provided is a simple example and may contain bugs or security vulnerabilities. Use it at your own risk.
