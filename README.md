#Life Game Bot


##Description
Life Game Bot is a twitter bot based on Life Game with cellular automaton and Malcov chain.

##Install
```
go get github.com/Rompei/lgb
```

##Usage

```
Usage:
  lgb [OPTIONS]

Application Options:
  -w, --width=               Field width (default: 160)
  -h, --height=              Field height (default: 90)
  -a, --alive-rate=          The rate of alive cells in initialization (default: 50)
  -m, --mut-rate=            The rate of mutation (default: 1)
  -g, --generation=          Reguration of generations
  -k, --keyword=             Keyword for twitter
  -l, --location=            Tweet location(Default: Japan) (default: 132.2,29.9,146.2,39.0,138.4,33.5,146.1,46.20)
  -d, --debug                Debug mode
  -t, --table                View type
  -s, --speed                Spead mode
      --consumer-key=        Twitter consumer key
      --consumer-secret=     Twitter consumer secret
      --access-token=        Twitter access token
      --access-token-secret= Twitter access token secret
```
