# Broadcast Server

A real-time TCP broadcast chat server written in Go with Message encryption using AES-GCM. Multiple clients can connect simultaneously and send encrypted messages that are broadcast to all other connected users.

## Features

- 🔐 **Message Encryption**: Messages are encrypted using AES-256-GCM
- 👥 **Multi-Client Support**: Multiple users can connect and chat simultaneously
- 🔄 **Real-Time Broadcasting**: Messages are instantly broadcast to all connected clients
- 🏷️ **Unique Usernames**: Each user has a unique username to identify their messages
- 🔒 **Thread-Safe**: Concurrent connection handling with proper mutex locks

## Architecture

The application consists of three main components:

1. **Server** (`servers/server.go`): Handles incoming connections, manages clients, and broadcasts messages
2. **Client** (`client/client.go`): Connects to the server, sends encrypted messages, and receives broadcasts
3. **Encryption** (`encryption/`): Provides AES-GCM encryption and decryption utilities

## Prerequisites

- Go 1.16 or higher
- `github.com/joho/godotenv` package

## Installation

1. Clone the repository:
```bash
git clone https://github.com/SanskarMali726/Broadcast-server.git
cd Broadcast-server
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory:
```env
PORT=:9000
```

## Usage

### Starting the Server

Run the server using:
```bash
go run main.go start
```

The server will start listening on the port specified in your `.env` file (default: 9000).

### Starting a Client

In a new terminal, run:
```bash
go run main.go connect
```

You'll be prompted to enter a username. Once connected, you can start sending messages.

### Running Multiple Clients

Open multiple terminals and run the client command in each to simulate multiple users chatting.

## How It Works

### Message Flow

1. **Client connects** to the server and provides a unique username
2. **Client sends a message**:
   - Generates a random 32-byte AES key
   - Encrypts the message using AES-GCM (produces ciphertext + 12-byte nonce)
   - Sends: `[4-byte length][32-byte key][12-byte nonce][encrypted message]`
3. **Server receives and processes**:
   - Reads the message length
   - Extracts the key, nonce, and encrypted message
   - Decrypts the message
   - Broadcasts the plaintext to all other connected clients

### Encryption Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 32 bytes (256 bits)
- **Nonce Size**: 12 bytes
- **Key Generation**: New random key per message (sent with the message)

## Project Structure

```
Broadcast-server/
├── client/
│   └── client.go          # Client implementation
├── servers/
│   └── server.go          # Server implementation
├── encryption/
│   └── encryption.go      # Encryption utilities
├── main.go                # Entry point
├── .env                   # Environment configuration
├── go.mod                 # Go module file
└── README.md              # This file
```

## Security Considerations

⚠️ **Important Security Limitations**: 

This is **NOT** end-to-end encryption. The current implementation has significant security limitations:

1. **Server Can Read Messages**: The server decrypts all messages before broadcasting them, meaning the server has access to plaintext messages
2. **Key Transmitted With Message**: Each message includes its encryption key, so anyone intercepting the network traffic can decrypt the messages
3. **No Authentication**: There's no verification of client identity
4. **No Message Integrity Between Clients**: While AES-GCM provides integrity for the client-to-server transmission, the server re-transmits messages in plaintext

**What This Implementation Provides:**
- Protection against basic packet sniffing (if keys aren't extracted)
- Learning example for encryption implementation in Go
- Foundation for building more secure systems

**For Production Use, You Need:**
- **True E2E Encryption**: Clients should encrypt messages with recipient's public key
- **TLS/SSL**: Encrypt the transport layer itself
- **Key Exchange Protocol**: Implement Diffie-Hellman or similar for secure key establishment
- **Authentication**: Use tokens, certificates, or OAuth for client authentication
- **Perfect Forward Secrecy**: Ensure past messages remain secure even if keys are compromised

## Example Session

```
Server Terminal:
$ go run main.go start
connect to the server on port :9000
Client connected 127.0.0.1:54321
[Alice]: Hello everyone!
[Bob]: Hi Alice!

Client 1 Terminal:
$ go run main.go connect
Connected to the server
Enter your username: Alice
Hello everyone!
[Bob]: Hi Alice!

Client 2 Terminal:
$ go run main.go connect
Connected to the server
Enter your username: Bob
[Alice]: Hello everyone!
Hi Alice!
```

## License

This project is open source and available under the [MIT License](LICENSE).

## Author

**Sanskar Mali**
- GitHub: [@SanskarMali726](https://github.com/SanskarMali726)

## Acknowledgments

- Uses [godotenv](https://github.com/joho/godotenv) for environment variable management
- Built with Go's standard `crypto/aes` and `crypto/cipher` packages

## Future Scope & Improvements

### Security Enhancements
- [ ] **Implement True End-to-End Encryption**: Use public-key cryptography (RSA/ECC) so only intended recipients can decrypt messages
- [ ] **Add TLS/SSL Support**: Secure the transport layer with certificates
- [ ] **Diffie-Hellman Key Exchange**: Implement secure key exchange protocol
- [ ] **User Authentication**: Add JWT tokens or certificate-based authentication
- [ ] **Message Signing**: Use digital signatures to verify message authenticity
- [ ] **Perfect Forward Secrecy**: Implement session keys that change regularly

### Features
- [ ] **Private Messaging**: Direct messages between specific users
- [ ] **Group Chats**: Create separate chat rooms/channels
- [ ] **Message History**: Store and retrieve chat history
- [ ] **File Sharing**: Support for sending files and images
- [ ] **User Status**: Show online/offline/typing indicators
- [ ] **Message Acknowledgments**: Delivery and read receipts
- [ ] **Rate Limiting**: Prevent spam and abuse
- [ ] **User Roles & Permissions**: Admin, moderator, regular user roles

### Technical Improvements
- [ ] **Database Integration**: Persist users and messages (PostgreSQL/MongoDB)
- [ ] **WebSocket Support**: Add WebSocket protocol alongside TCP
- [ ] **REST API**: Add HTTP endpoints for web client integration
- [ ] **Logging**: Structured logging with levels (info, warn, error)
- [ ] **Metrics & Monitoring**: Add Prometheus metrics and health checks
- [ ] **Docker Support**: Containerize the application
- [ ] **Configuration Management**: Better config handling (YAML/TOML)
- [ ] **Graceful Shutdown**: Handle SIGTERM/SIGINT properly
- [ ] **Connection Pooling**: Optimize resource usage
- [ ] **Message Queue**: Use Redis/RabbitMQ for scalability


