import socket
import threading
import socketserver as SocketServer
from socket import error as SocketError
import errno


class SingleTCPHandler(SocketServer.BaseRequestHandler):
    def handle(self):
        while True:
            try:
                data = self.request.recv(1024)  # clip input at 1Kb
            except SocketError as e:
                if e.errno != errno.ECONNRESET:
                    raise  # Not error we are looking for
                break
            if data == '':
                self.request.close()
                break
            print(f'recvd from: {self.client_address} data: {data.decode("utf-8")}')
            response = f'answer : {data.decode("utf-8")}'
            self.request.sendall(response.encode('utf-8'))


class SimpleServer(SocketServer.ThreadingMixIn, SocketServer.TCPServer):
    # Ctrl-C will cleanly kill all spawned threads
    daemon_threads = True
    # much faster rebinding
    allow_reuse_address = True

    def __init__(self, server_address, RequestHandlerClass):
        SocketServer.TCPServer.__init__(
            self, server_address, RequestHandlerClass)


if __name__ == "__main__":
    # Port 0 means to select an arbitrary unused port
    HOST, PORT = "127.0.0.1", 20000
    print("Server starting on {}:{}".format(HOST, PORT))
    server = SimpleServer((HOST, PORT), SingleTCPHandler)
    # terminate with Ctrl-C
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        exit(0)
