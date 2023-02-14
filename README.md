## oai

An openai cli written in go, supports proxying traffic over
tor by default.

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

