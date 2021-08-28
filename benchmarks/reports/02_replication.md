# Replication


## Baseline

Empty db
```
CONTAINER ID   NAME                  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O        PIDS
f76507060844   social-network_db_1   0.61%     407.1MiB / 1.941GiB   20.49%    16.1kB / 14.4kB   94.9MB / 286MB   38
```

While loading data
```
> make REGISTER_DURATION=5m BENCH_N_CONN=50 bench-register

CONTAINER ID   NAME                  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O          BLOCK I/O        PIDS
f76507060844   social-network_db_1   211.28%   533.1MiB / 1.941GiB   26.82%    144MB / 133MB   95.2MB / 1.28GB   90
```

Benchmark results
```
> wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   126.45ms   73.06ms 691.44ms   86.98%
    Req/Sec   145.57     48.37   420.00     75.87%
  Latency Distribution
     50%  101.52ms
     75%  130.68ms
     90%  227.78ms
     99%  379.77ms
  121030 requests in 5.00m, 11.20MB read
  Socket errors: connect 0, read 0, write 0, timeout 48
Requests/sec:    403.31
Transfer/sec:     38.20KB
```

While reading
```
> make BENCH_DURATION=5m BENCH_N_CONN=50 BENCH_TIMEOUT=10s bench-search

CONTAINER ID   NAME                  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O         BLOCK I/O         PIDS
f76507060844   social-network_db_1   509.02%   546.8MiB / 1.941GiB   27.51%    182MB / 344MB   95.5MB / 1.28GB   90
```

Benchmark results
```
wrk --latency --timeout 10s -d 5m -t 3 -c 50 -s benchmarks/search.lua http://localhost

Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   297.99ms  376.08ms   3.05s    85.83%
    Req/Sec    91.60     58.81   740.00     72.65%
  Latency Distribution
     50%  158.11ms
     75%  464.36ms
     90%  792.06ms
     99%    1.66s
  81620 requests in 5.00m, 165.47MB read
  Non-2xx or 3xx responses: 48
Requests/sec:    271.98
Transfer/sec:    564.62KB
```