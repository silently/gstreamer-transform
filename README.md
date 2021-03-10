# gstreamer-transform

gstreamer-send is a simple WebRTC application that shows how to:

- send audio and video from the browser to a pion-WebRTC remote peer
- transform audio/video through GStreamer
- display them back in the browser

It's mainly a refactor of [gstreamer-send](https://github.com/pion/example-webrtc-applications/tree/master/gstreamer-send), [gstreamer-receive](https://github.com/pion/example-webrtc-applications/tree/master/gstreamer-receive) and [reflect](https://github.com/pion/webrtc/tree/master/examples/reflect) pion example projects.

## Instructions

### Install GStreamer

This example requires you have GStreamer installed, these are the supported platforms

#### Debian/Ubuntu

`sudo apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good`

#### Windows MinGW64/MSYS2

`pacman -S mingw-w64-x86_64-gstreamer mingw-w64-x86_64-gst-libav mingw-w64-x86_64-gst-plugins-good mingw-w64-x86_64-gst-plugins-bad mingw-w64-x86_64-gst-plugins-ugly`

#### macOS

`brew install gst-plugins-good gst-plugins-ugly pkg-config && export PKG_CONFIG_PATH="/usr/local/opt/libffi/lib/pkgconfig"`

### Build

```
go build
```

### Open gstreamer-send example page

Serve the `web/index.html` static page.

For instance with [http-server](https://github.com/http-party/http-server)

```
http-server web -p 8080
```

### Run gstreamer-transform with your browsers SessionDescription as stdin

Open http://localhost:8080/ and copy the contents of the top textarea and:

#### Linux/macOS

Run `echo $BROWSER_SDP | ./gstreamer-transform`

#### Windows

1. Paste the SessionDescription into a file.
1. Run `gstreamer-send < my_file`

### Input gstreamer-send's SessionDescription into your browser

Copy the text that `gstreamer-send` just emitted and copy into second text area

### Hit 'Start Session' in jsfiddle, enjoy your video!

The audio and video captured from your browser should be transformed (audio echo and warp video effects).

Congrats, you have used pion-WebRTC! Now start building something cool

## About GStreamer

When prototyping with GStreamer it is highly recommended that you enable debug output, this is done by setting the `GST_DEBUG` enviroment variable.
You can read about that [here](https://gstreamer.freedesktop.org/data/doc/gstreamer/head/gstreamer/html/gst-running.html) a good default value is `GST_DEBUG=*:3`

You can also prototype a GStreamer pipeline by using `gst-launch-1.0` to see how things look before trying them with `gstreamer-send` for the examples below you
also may need additional setup to enable extra video codecs like H264. The output from GST_DEBUG should give you hints
