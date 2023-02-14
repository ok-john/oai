## oai

An openai cli written in go, supports proxying traffic over tor. 


## install

1) Grab the [latest release](https://github.com/ok-john/oai/releases/tag/v0.0.1), there are precompiled binaries available for the following architechtures.

```
darwin-amd64  freebsd-amd64  freebsd-arm64  linux-amd64  linux-arm64  linux-mips64    linux-ppc64    linux-s390x    windows-arm5   windows-arm7
darwin-arm64  freebsd-arm5   freebsd-arm7   linux-arm5   linux-arm7   linux-mips64le  linux-ppc64le  windows-386    windows-arm6
freebsd-386   freebsd-arm6   linux-386      linux-arm6   linux-mips   linux-mipsle    linux-riscv64  windows-amd64  windows-arm64
```

2) Extract the tarball and make a directory to store your API keys, if you don't want to pass that directory
   into the cli each time then make a folder at the default location `/etc/openai`.

```bash
tar xf $OS-$ARCH.tar.zst
mkdir -p /etc/openai
echo 'MY_OPENAI_KEY' > /etc/openai/.api
```

3) Put `oai` in a `$PATH` directory:

```bash
mv oai /usr/local/bin
```

And you're all set! Have fun!


## Be amazed

As openai implements repeating key xor in pure bash:

<p float="left" align="middle">
<img src="./oai.svg" />
</p>

### usage
```bash
Usage of ./oai:
  -env string
        absolute path to the environment directory. (default "/etc/openai")
  -max int
        model to use (default 2000)
  -model string
        model to use (default "text-davinci-003")
  -o string
        optionally specify an output file, defaults to stdout.
  -org string
        optionally specify an organization id.
  -socks5-hostname string
        optionally override the default tor proxy address. (default "localhost:9050")
  -temp float
        model to use (default 0.5)
  -tor
        toggles the use of a socks5 tor proxy. (default true)
```

