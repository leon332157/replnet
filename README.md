# replit-code-server

Use VSCode on [Repl.it](https://replit.com)

### This should be a proof of concept which allows the usage of ssh-based vscode remote editing with replit

Implementation idea is as follow:
### Server
- Installes openssh server/run it
- Proxy this TCP connection over the https port with websocket
- Let the extension do the rest

### Client
- A client on a local PC which has a CLI (maybe gui)
- This client acts like a ssh server but actually proxies that traffic over websockets to the server in a repl.
#### A visual implementation
```
┌──────────────────────────────────┐                          ┌─────────────────────────────────┐
│                                  │                          │                                 │
│                                  │                          │                                 │
│      PC (Client)                 │  WebSocket (Proxy SSH)   │       In a Repl (Server)        │
│                                  ┼────────────────────────► │                                 │
│                                  │                          │       With an SSH server        │
│                                  │                          │                                 │
│                                  │                          │                                 │
└──────────────────────────────────┘                          └─────────────────────────────────┘
```

This poses a few problems:
- Will the forwarded TCP traffic be safe and trusted?
- Will the server be able to install and run an openssh-server
- And if it doesn't, how should the ssh-server be implemented?
- How will it ensure that the person connecting to the server is the repl owner/ securty concerns. 
- Performance issues?