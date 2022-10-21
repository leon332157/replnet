# Replnet design

## Introduction

This document describes the design of `replnet` which is a command line application that aims to open TCP socket into a repl. It is a work in progress and is not yet ready for use.

## Motivation

The main motivation for this project is to allow for raw socket connections from and to a repl. This is useful for learning about socket programming and for creating a repl that can be used as a server. There are also other use cases such as allowing for inter-repl communication.

## Technical Design 

Currerrently, replit only allows for the opening of port proxied over http and a server capable of handling http requests is required. To take advantage of that, the websocket protocol is a good solution for this problem. This also allows for other extensions based on http to be used such as WebDav to allow for directory listing and file transfer to a repl, which is potentially useful for syncing a repl with local files.

## Architecture

The visual representation of the project is as follows:
```
┌──────────────────────────────────┐                          ┌─────────────────────────────────┐
│                                  │                          │                                 │
│                                  │                          │                                 │
│                                  │  WebSocket (TCP traffic) │       In a Repl (Server)        │
│           PC (Client)            | ───────────────────────► │                                 │
│ Open local port for forwarding   │                          │       Websocket, WebDav         │
│                                  │                          │                                 │
│                                  │                          │                                 │
└──────────────────────────────────┘                          └─────────────────────────────────┘
```

### Client detail

Similar to ssh port forwarding, the client should expect to be configured to listen on a local port, a host to connect to, and a remote port to forward to. The client must be able to run in the background and should be able to be configured to run on startup.

#### The client must contain the following configurations:

- `replnet.client.remote-url` Host http url to connect to, must be `http(s) or ws(s)` eg:"wss://replnet.leon332157.repl.co"
- `replnet.client.remote-port` A remote port to connect to (1-65535)
- `replnet.client.local-port` The local port can be equal to the remote port, but should also be configurable (1-65535)

### Server detail

#### The server must contain the following configurations:

- `replnet.server.listen-port` OPTIONAL: A port to listen on, default to 0, which allows kernel to automatically assign a port (0-65535)
- `replnet.server.reverse-proxy-port` OPTIONAL: A port to reverse proxy http traffic to if the repl is hosting a http server
