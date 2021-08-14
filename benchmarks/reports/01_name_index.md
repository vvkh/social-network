# Adding index for searching profiles by name

All tests were performed on the following specs

| cpu | ram | storage | os | docker | 
| --- | --- | --- | --- | --- | 
| i7-9750H | 32GB | 512GB SSD | macOS 11.5 | 3.5.2 |

## Baseline

### Generate requests for load testing
```
> make gen-bench
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
> make REGISTER_N_CONN=50 bench-register

wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost

Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   450.63ms   97.74ms 989.36ms   71.23%
    Req/Sec    36.09     20.04   140.00     68.83%
  Latency Distribution
     50%  446.89ms
     75%  509.99ms
     90%  573.48ms
     99%  703.88ms
  31913 requests in 5.00m, 2.95MB read
  Socket errors: connect 0, read 0, write 0, timeout 96
Requests/sec:    106.34
Transfer/sec:     10.07KB
```

### `/profiles/` performance
Seed users data through running `/register/` benchmark longer.
You may want to restart service with weaker bcrypt cost to temporary speed up seeding data.
```
> export BCRYPT_COST=4
> make up
```

```
> make BENCH_DURATION=30m BENCH_N_CONN=110 bench-register
```

As a result, there are 980K users in the database.

Now run benchmark for `/profiles/` page.
```
> make BENCH_DURATION=5m BENCH_N_CONN=50 BENCH_TIMEOUT=10s bench-search

wrk --latency --timeout 10s -d 5m -t 3 -c 50 -s benchmarks/search.lua http://localhost

Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.12s     2.87s    8.07s    34.93%
    Req/Sec    11.34     14.84   250.00     87.89%
  Latency Distribution
     50%    3.59s
     75%    6.27s
     90%    7.06s
     99%    7.72s
  5042 requests in 5.00m, 10.20MB read
  Non-2xx or 3xx responses: 48
Requests/sec:     16.80
Transfer/sec:     34.81KB
```

## Profiles index
The reason of such performance is that db has to use full seqscan through `profiles` table.

```
> EXPLAIN SELECT * FROM profiles WHERE first_name LIKE "K" AND last_name LIKE "Lep"

-> Filter: ((`profiles`.first_name like 'K') and (`profiles`.last_name like 'Lep'))  (cost=97572.85 rows=11894)
    -> Table scan on profiles  (cost=97572.85 rows=963631)
```

Add index to address the issue.
```
> cat migrations/005_profiles_name_index.up.sql

create index profile_names on profiles(first_name, last_name);
```

Migrate db.
```
> make migrate
```

Now we got query plan with index scan for the query.
```
> EXPLAIN FORMAT = tree SELECT * FROM profiles WHERE first_name LIKE "K%" AND last_name LIKE "Lep%"

-> Index range scan on profiles using profile_names, with index condition: ((`profiles`.first_name like 'K%') and (`profiles`.last_name like 'Lep%'))  (cost=57560.66 rows=127912)
```

First name only query.
```
> EXPLAIN FORMAT = tree SELECT * FROM profiles WHERE first_name LIKE "K%" AND last_name LIKE "%"

-> Index range scan on profiles using profile_names, with index condition: ((`profiles`.first_name like 'K%') and (`profiles`.last_name like '%'))  (cost=57560.66 rows=127912)
```

Last name only query.
```
> EXPLAIN FORMAT = tree SELECT * FROM profiles WHERE first_name LIKE "%" AND last_name LIKE "Lep%"

-> Filter: ((`profiles`.first_name like '%') and (`profiles`.last_name like 'Lep%'))  (cost=97572.85 rows=11894)
    -> Table scan on profiles  (cost=97572.85 rows=963631)
```

Created index can't be used for last name only search so db falls back to seqscan.

Add last name index to address the issue.
```
> cat migrations/006_profiles_last_name_index.up.sql

create index profile_last_names on profiles(last_name, first_name);
```

Migrate db.
```
> make migrate
```

Query plan for last name only filter.
```
> EXPLAIN FORMAT = tree SELECT * FROM profiles WHERE first_name LIKE "%" AND last_name LIKE "Lep%"

-> Filter: (`profiles`.first_name like '%')  (cost=441.26 rows=109)
    -> Index range scan on profiles using profile_last_name, with index condition: (`profiles`.last_name like 'Lep%')  (cost=441.26 rows=980)
```

Run benchmarks again:
```
> make BENCH_DURATION=5m BENCH_N_CONN=50 BENCH_TIMEOUT=10s bench-search

wrk --latency --timeout 10s -d 5m -t 3 -c 50 -s benchmarks/search.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   153.98ms  237.47ms   1.78s    85.71%
    Req/Sec   336.58    123.49     1.12k    77.27%
  Latency Distribution
     50%   47.36ms
     75%  129.34ms
     90%  519.08ms
     99%    1.06s
  220149 requests in 5.00m, 469.12MB read
  Non-2xx or 3xx responses: 48
Requests/sec:    733.56
Transfer/sec:      1.56MB
```

Recheck register performance after adding index
```
> make REGISTER_N_CONN=50 bench-register

wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   463.00ms  107.32ms 998.83ms   71.05%
    Req/Sec    35.10     20.05   151.00     69.25%
  Latency Distribution
     50%  459.08ms
     75%  527.42ms
     90%  596.87ms
     99%  745.55ms
  30767 requests in 5.00m, 2.85MB read
  Socket errors: connect 0, read 0, write 0, timeout 246
Requests/sec:    102.53
Transfer/sec:      9.71KB
```

## Final results
Test results show that adding index significantly increased search performance while not affecting register performance that much.

| Test | Search, Latency 90% | Search, Throughput | Register, Latency 90% | Register, Throughput |
| --- | --- | --- | --- | --- |
| Baseline | ~7 s | ~17 rps | ~600 ms | ~100 rps | 
| With index | ~500 ms | ~733 rps |  ~600ms |   ~100 rps |
