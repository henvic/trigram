# trigram
After compiling, you can see a list of options for running the trigram server with:

`./trigram -h`

You might want to use -race to evaluate the program.

## Technical considerations and implementation details

* No external libraries were used. Using [ctxsignal](https://github.com/henvic/ctxsignal) and [logrus](https://github.com/sirupsen/logrus) on a few points would help, though.
* Minor change to punctuation were made to make text presentation clearer.
* An assumption of serving this service behind a reverse proxy such as nginx was made.

## Running
It is recommended to use the `-expose-debug` flag to expose debugging data (from packages expvar and pprof) on HTTP local port 8081 (including on production environments), allowing you to run commands such as:

```
$ go tool pprof -web http://localhost:8081/debug/pprof/heap
$ curl http://localhost:8081/debug/vars
```

## Concurrency and stress testing

```
ab -r -p free-trade.txt -T text/plain -n 1 -c 1 http://127.0.0.1:8080/learn
ab -n 100000 -c 50 http://127.0.0.1:8080/generate
```

