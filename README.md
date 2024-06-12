# diskgo

TODO:
- [x] basic CI (fmt, lint, build, test, ...)
- [x] mdns discovery for client
- [x] mdns announce for server
- [ ] timeout handling (using contexts ?)
- [x] disk file storage (instead of in memory)
- [x] server configuration (max disk space provided, store directory)
- [ ] client configuration (encryption, replicas)
  - [x] encryption
  - [ ] replicas
- [ ] docs
- [ ] optimisations (allocations, buffers, ...)
- [x] provide more errors from server to clients when an operation fails
- [ ] handle big files (not fitting in RAM)
- [ ] protocol version
- [ ] handle server restart (retrieve root dir's current size, ...)
- [ ] handle client restart (retrieve where files are, ...)
- [ ] node (combine client and server)

Future ideas:
- [ ] native desktop app for macOS (using SwiftUI and GRPC to interact with Golang)
- [ ] integrate with the filesystem and mount a volume
