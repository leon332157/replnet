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
│                                  │  WebSocket (Proxy SSH)   │       In a Repl (Server)        │
│           PC (Client)            ┼────────────────────────► │                                 │
│                                  │                          │       With an SSH server (Go)   │
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

## Update 1:
After some exprimenting, it seems that running an ssh-server would require a lot of effort and it's not guranteed that it would work with the vscode extension. 

## Update 2:
I have looked into some ssh server implementation written in go, it does seem like the extension works with an GO implementaion. However, the example I used did not have `direct-tcpip` channel, which is needed for port forwarding. So now, instead of trying to proxy an TCP connection to an openssh-server, the server itself should be an ssh server but with websockt over an SSH protocol. 
PS. I have only tested the ssh server with my own setup, I haven't tested it with replit in an repl. 

## Update 3:
The example ssh server I used does seem to partially work with the vscode extension however without the support of `direct-tcpip`, most of the functionality is not possble.
