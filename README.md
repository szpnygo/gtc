# gtc: Golang Terminal Chat
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/szpnygo/gtc)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/szpnygo/gtc?label=version)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/neosu/gtc?label=docker%20version)
![GitHub last commit](https://img.shields.io/github/last-commit/szpnygo/gtc)
![GitHub Release Date](https://img.shields.io/github/release-date/szpnygo/gtc)
![GitHub issues](https://img.shields.io/github/issues-raw/szpnygo/gtc)
![GitHub top language](https://img.shields.io/github/languages/top/szpnygo/gtc)

gtc is a p2p terminal chat app, using webrtc. All conversation content will not pass through the server.

## Getting Started

### Install

```bash
go install github.com/szpnygo/gtc@v0.2.2
```

or you can download from release

### Quick Run
```bash
gtc -s wss://gogs.tcodestudio.com
```

### Use your own service
```bash
gtc server -r golang,movie,meta,job,tech,idea
gtc -s ws://127.0.0.1:8888
```

### Docker

[Docker](https://hub.docker.com/r/neosu/gtc)

## License
[Apache License Version 2.0](./LICENSE)
