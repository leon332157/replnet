"""
 Implements a simple HTTP/1.0 Server

"""

import socket


# Define socket host and port
SERVER_HOST = '127.0.0.1'
SERVER_PORT = 7777

# Create socket
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server_socket.bind((SERVER_HOST, SERVER_PORT))
server_socket.listen(1)
print('Listening on port %s ...' % SERVER_PORT)

while True:    
    # Wait for client connections
    client_connection, client_address = server_socket.accept()
    print(f'accepted from :{client_address}')
    # Get the client request
    request = client_connection.recv(1024).decode()
    print(request.split('\n')[0])

    # Send HTTP response
    response = f'HTTP/1.0 200 OK\r\nConnection: close\r\n\r\nHello {client_address} \n\n\nrequest: {request}\r\n'
    client_connection.sendall(response.encode())
    client_connection.close()

# Close socket
server_socket.close()