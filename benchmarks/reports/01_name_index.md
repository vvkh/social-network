# Adding index for searching profiles by name

## Base line

### Generate requests for load testing
```
make gen-bench
```

Sample requests are saved to `benchmarks/requests` directory.

`/register/` requests create users with random first name and last name picked from `benchmark/data/names.txt`.
Requests are saved as url-encoded form body.

Examples:
```
> head benchmarks/requests/register.txt
about=&age=55&first_name=Kristbj%C3%B6rn&last_name=Fein&location=Russia&password=kristbj%C3%B6rnfein&sex=female&username=kristbj%C3%B6rnfein
about=&age=37&first_name=Kristbj%C3%B6rn&last_name=Stonelake&location=Germany&password=kristbj%C3%B6rnstonelake&sex=female&username=kristbj%C3%B6rnstonelake
about=&age=51&first_name=Kristbj%C3%B6rn&last_name=Lepage&location=Germany&password=kristbj%C3%B6rnlepage&sex=male&username=kristbj%C3%B6rnlepage
about=&age=22&first_name=Kristbj%C3%B6rn&last_name=Grona&location=Poland&password=kristbj%C3%B6rngrona&sex=male&username=kristbj%C3%B6rngrona
about=&age=22&first_name=Kristbj%C3%B6rn&last_name=Tuma&location=Russia&password=kristbj%C3%B6rntuma&sex=female&username=kristbj%C3%B6rntuma
about=&age=24&first_name=Kristbj%C3%B6rn&last_name=Froeliger&location=Poland&password=kristbj%C3%B6rnfroeliger&sex=female&username=kristbj%C3%B6rnfroeliger
about=&age=77&first_name=Kristbj%C3%B6rn&last_name=Royea&location=Canada&password=kristbj%C3%B6rnroyea&sex=male&username=kristbj%C3%B6rnroyea
about=&age=48&first_name=Kristbj%C3%B6rn&last_name=Kiehn&location=Poland&password=kristbj%C3%B6rnkiehn&sex=male&username=kristbj%C3%B6rnkiehn
about=&age=71&first_name=Kristbj%C3%B6rn&last_name=Wineinger&location=Russia&password=kristbj%C3%B6rnwineinger&sex=female&username=kristbj%C3%B6rnwineinger
about=&age=32&first_name=Kristbj%C3%B6rn&last_name=Lollar&location=UK&password=kristbj%C3%B6rnlollar&sex=male&username=kristbj%C3%B6rnlollar

```
`/profiles/` requests search for random last name and first name prefix, either can be empty.
Prefixes are randomly generated from the names in `benchmarks/data/names.txt`. Requests are saved as url-encoded query parameters.

Examples:
```
> head benchmarks/requests/search.txt
first_name=K&last_name=
first_name=Kristbj&last_name=Sto
first_name=K&last_name=Lep
first_name=Krist&last_name=
first_name=Kristb&last_name=
first_name=Kris&last_name=Froelig
first_name=Kr&last_name=R
first_name=Kristbj%C3&last_name=Ki
first_name=K&last_name=Winei
first_name=Kristbj&last_name=L
```

### Start server
```
make env
make up
```

### `/register/` performance

Run `/register/` benchmark.
```
make REGISTER_N_CONN=50 bench-register
```
Results:
```
> wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   377.57ms   87.84ms 974.41ms   72.18%
    Req/Sec    42.82     22.49   161.00     60.82%
  Latency Distribution
     50%  371.36ms
     75%  427.96ms
     90%  487.13ms
     99%  623.24ms
  38185 requests in 5.00m, 3.53MB read
  Socket errors: connect 0, read 0, write 0, timeout 49
Requests/sec:    127.26
Transfer/sec:     12.05KB
```

### `/profiles/` performance
Seed users data through running `/register/` benchmark longer.
You may want to restart service with weaker bcrypt cost to temporary speed up seeding data.
```
export BCRYPT_COST=4
make up
```

```
make BENCH_DURATION=30m BENCH_N_CONN=110 bench-register
```

As a result, there are 980K users in the database.

Now run benchmark for `/profiles/` page.
```
make BENCH_DURATION=10s BENCH_N_CONN=50 BENCH_TIMEOUT=10s
```
Results:
```
> wrk --latency --timeout 10s -d 10s -t 3 -c 50 -s benchmarks/search.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.92s     3.38s    9.04s    70.83%
    Req/Sec    21.07     55.67   290.00     93.18%
  Latency Distribution
     50%   70.35ms
     75%    6.65s
     90%    7.66s
     99%    8.99s
  120 requests in 10.09s, 30.99MB read
  Non-2xx or 3xx responses: 48
Requests/sec:     11.90
Transfer/sec:      3.07MB
```
