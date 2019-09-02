# Task Tracker

Inspired by [@campoy](https://github.com/campoy)'s JustForFunc series #30 video: https://www.youtube.com/watch?v=_jQ3i_fyqGA

## Setup

1. Download protoc binary and add it to your PATH

See [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) for
the binary download. Then just add the `bin/protoc` binary to your PATH.

2. Install protoc-gen-go plugin
  * `go get -u github.com/golang/protobuf/protoc-gen-go`
  * Add this binary or `${GOBIN}` to your path so `protoc` can find the plugin

## Debugging

If you have `protoc` installed, you can view your protobuf data at 
anytime with the following command:

`cat db.pb | protoc --decode_raw`
