# Adding index for searching profiles by name

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
    Latency     7.29s     2.41s   10.00s    84.00%
    Req/Sec     3.91     11.15   282.00     97.46%
  Latency Distribution
     50%    7.91s
     75%    8.73s
     90%    9.45s
     99%    9.95s
  1496 requests in 5.00m, 2.61GB read
  Socket errors: connect 0, read 0, write 0, timeout 496
  Non-2xx or 3xx responses: 48
Requests/sec:      4.98
Transfer/sec:      8.89MB
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

create index profile_last_name on profiles(last_name);
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

Now index scan is used, but there is a useless first_name filter. Rewrite backend to produce following query for last-name only search
```
> EXPLAIN FORMAT = tree SELECT * FROM profiles WHERE last_name LIKE "Lep%"

-> Index range scan on profiles using profile_last_name, with index condition: (`profiles`.last_name like 'Lep%')  (cost=441.26 rows=980)
```

Run benchmarks again:
```
> make BENCH_DURATION=5m BENCH_N_CONN=50 BENCH_TIMEOUT=10s bench-search

Running 5m test @ http://localhost
  3 threads and 50 connections
user wrk created
authenticated wrk with cookie 	token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjg1OTUxMTcsInByb2ZpbGUiOjExNDg3MzAsInN1YiI6MTE0OTk3Mn0.EkqiRJJSd5N_sfE-aZw9ZXVW5MNvLSg2U8dsbcT_emY; Path=/; HttpOnly	 status 	302
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   865.99ms    1.49s   10.00s    93.07%
    Req/Sec     8.49     10.94   254.00     89.53%
  Latency Distribution
     50%  453.15ms
     75%  677.75ms
     90%    1.45s
     99%    8.63s
  4879 requests in 5.00m, 6.84GB read
  Socket errors: connect 0, read 0, write 0, timeout 179
  Non-2xx or 3xx responses: 48
Requests/sec:     16.26
Transfer/sec:     23.36MB
```
