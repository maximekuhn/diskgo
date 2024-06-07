# diskgo

TODO:
- [x] basic CI (fmt, lint, build, test, ...)
- [ ] mdns discovery for client
- [ ] mdns announce for server
- [ ] timeout handling (using contexts ?)
- [ ] disk file storage (instead of in memory)
- [ ] server configuration (max disk space provided, store directory)
- [ ] client configuration (encryption, replicas)
- [ ] more "security" (send peer name when performing requests, ...)
- [ ] docs
- [ ] optimisations (allocations, buffers, ...)
- [ ] provide more errors from server to clients when an operation fails
- [ ] handle big files (not fitting in RAM)

Future ideas:
- [ ] native desktop app for macOS (using SwiftUI and GRPC to interact with Golang)
- [ ] integrate with the filesystem and mount a volume
