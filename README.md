# replnet [![Run on Replit](https://replit.com/badge/github/ReplDepot/replnet)](https://replit.com/github/ReplDepot/replnet)

## What?
* A repl to repl networking tool
* Also allow you to connect from local computer 
  
## Why?
* In order to connect to raw socket on a repl, you need to proxy it through http/websocket
* This makes it harder for a learner to learn about socket programming

## Goals:

* Single binary for convinence and compatibility. 
* Compatiblity with OS-native webdav mount
* Easy to use, self-explainatory, beginner friendly.

## Functionality:

* ✅ Basic webdav mounting of a repl with basic auth
* ❌ Health check route, ping route 
* ⚠️ WebSocket Proxy for port forwarding similar to ngrok, with a local client. 