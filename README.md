# gojam

An implementation of a Jamulus client written in Go, intended for headless streaming use-cases.

To be able to use this client, we need to build a custom version of Opus from Jamulus source tree.

```sh
git clone https://github.com/jamulussoftware/jamulus.git
bash -c 'cd jamulus/libs/opus && ./configure --enable-static --disable-shared --enable-custom-modes --disable-hardening && make'
```

Then provide these environment variables to `go build` or `go run` in order to successfully build and run your code:

```sh
export CGO_CFLAGS="-I$PWD/jamulus/libs/opus/include"
export CGO_LDFLAGS="-L$PWD/jamulus/libs/opus/.libs"
```

## Usage

A simple CLI, `gojamclient` is provided that can connect to a Jamulus server and send received audio to another TCP server as raw PCM samples (16-bit signed, little endian, 48kHz, 2 channels).

First, you need to run a TCP server that will receive the audio samples. For example, you can use a combination of `nc` and `play` (from `sox` package) to listen to the audio:

```sh
nc -l 22222 | play -r 48000 -c 2 -b 16 -e signed -t raw -
```

Then you can run `gojamclient`:

```sh
go run ./cmd/gojamclient/ -server <ip>:<port> -pcmout localhost:22222
```

### REST API

When you run an API server, e.g. `-apiserver localhost:9999`, you can control the client, for example, using [HTTPie](https://httpie.io/):

```sh
# Get the client info
http get localhost:9999/channel-info

# Update the client info
http patch localhost:9999/channel-info name=newname
```

### Streaming Jamulus audio to Discord

There is a separate project [pcm2discord](https://github.com/dtinth/pcm2discord) that receives raw PCM samples via TCP, and sends them to Discord.

<details><summary>To make the Jamulus change its name to display the number of listeners in Discord, the following Python script can be used:</summary>

```python
import requests
import time

last_name = None

while True:
    try:
        # Get the count of listeners from Discord
        r = requests.get('http://localhost:28280/count')

        # Response is in form: { "listening": 2 }
        # Get the number of listeners
        listeners = r.json()['listening']

        # Submit the number of listeners to channel info endpoint
        name = ' Discord[' + str(listeners) + ']'
        r = requests.patch('http://localhost:28281/channel-info', json={'name': name})

        if last_name != name:
            print('Updated channel name to ' + name + ' at ' + time.strftime('%H:%M:%S'))
            last_name = name
    except Exception as e:
        print('Error: {}'.format(e))
    finally:
        # Wait 2 seconds
        time.sleep(2)
```

</details>
