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

## Update 4:
The proof of concept was successful. I was able to forward SSH/TCP traffic over a [websocket runnel](https://github.com/derhuerst/tcp-over-websockets). Now, I think the project should be broken into two parts. 
- A websocket tunnel to forward raw tcp traffic from a repl to a computer. Integrated with [replit-desktop](https://github.com/replit-discord/replit-desktop) OR Golang with native GUI. 
- A custom SSH server compatible with openssh client written in golang that listens on 127.0.0.1 with support for vscode remote. 

But it can also be **One unified Binary with client/server commandline argument and/or gui combining websocket tunnel and the ssh server into one binary. Which is harder to do but maybe easier to setup on replit. 

# Server Implementation idea
- Hybrid HTTP Proxy/SSH server
- Listen on 0.0.0.0 in a repl and allow config file/commandline argument to forward raw http traffic (repl.co) to the underlying applicaion if applicable. 
- Use path /ssh for the path of the websocket which forwards to the ssh server with io.Copy
- Implement the ssh server first to make sure it's compatible with vscode, then implement the proxies. 
