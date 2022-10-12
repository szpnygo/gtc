# gtc: Golang Terminal Chat

gtc is a p2p terminal chat app, using webrtc. All conversation content will not pass through the server.

## TODO
- [ ] Support config the ice server
- [ ] Room list shows the number of people online

## Getting Started

### Install

```bash
go install github.com/szpnygo/gtc@v0.0.1
```

or you can download from release

### Quick Run
```bash
gtc -s wss://gogs.tcodestudio.com
```

### Use your own service
```bash
gtc server
gtc -s ws://127.0.0.1:8888
```

## License
[Apache License Version 2.0](./LICENSE)