# gojam

Go implementation of a Jamulus client.

To be able to use this client, we need to build a custom version of Opus from Jamulus source tree.

```sh
git clone https://github.com/jamulussoftware/jamulus.git
bash -c 'cd jamulus/libs/opus && ./configure --enable-static --disable-shared --enable-custom-modes --disable-hardening && make'
```

Then provide these environment variables to `go build`:

```sh
export CGO_CFLAGS="-I$PWD/jamulus/libs/opus/include"
export CGO_LDFLAGS="-L$PWD/jamulus/libs/opus/.libs"
```