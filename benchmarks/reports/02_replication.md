# Replication


## Baseline

### Resources usage while idle
```
CONTAINER ID   NAME                  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O        PIDS
f76507060844   social-network_db_1   0.61%     407.1MiB / 1.941GiB   20.49%    16.1kB / 14.4kB   94.9MB / 286MB   38
```

### Resources usage while writing
```
> make REGISTER_DURATION=5m BENCH_N_CONN=50 bench-register

CONTAINER ID   NAME                  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O          BLOCK I/O        PIDS
f76507060844   social-network_db_1   211.28%   533.1MiB / 1.941GiB   26.82%    144MB / 133MB   95.2MB / 1.28GB   90
```

### Benchmark results
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

### Resources usage while reading
```
> make BENCH_DURATION=5m BENCH_N_CONN=50 BENCH_TIMEOUT=10s bench-search

CONTAINER ID   NAME                  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O         BLOCK I/O         PIDS
f76507060844   social-network_db_1   509.02%   546.8MiB / 1.941GiB   27.51%    182MB / 344MB   95.5MB / 1.28GB   90
```

### Benchmark results
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

## Async replication

Async replication with one replica was configured. Read queries were routed to replica, write queries were routed to main instance.
Query routing was configured through `proxysql` rule:
```
mysql_query_rules:
(
    {
        rule_id=1
        active=1
        match_pattern="^SELECT \* FROM `profiles`"
        destination_hostgroup=1
        apply=1
    },
)
```

Main node config:
```
[mysqld]
server-id=1
binlog_format=ROW
log-bin
```

### Resources usage while writing
```
> wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost

CONTAINER ID   NAME                          CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O        PIDS
3bfa9158cc07   social-network_db_proxy_1     66.62%    18MiB / 1.939GiB      0.91%     106MB / 126MB     16.4kB / 0B      19
15faa3501e76   social-network_db_replica_1   56.59%    382.3MiB / 1.939GiB   19.26%    55.7MB / 534kB    43.9MB / 773MB   43
24718d51a7b1   social-network_db_1           111.96%   447.2MiB / 1.939GiB   22.53%    59.5MB / 96.6MB   33.7MB / 743MB   116
```

### Benchmark results
```
wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   149.48ms   45.78ms 480.25ms   71.54%
    Req/Sec   107.65     32.67   303.00     70.23%
  Latency Distribution
     50%  142.12ms
     75%  174.11ms
     90%  210.77ms
     99%  287.49ms
  95890 requests in 5.00m, 8.87MB read
  Socket errors: connect 0, read 0, write 0, timeout 144
Requests/sec:    319.54
Transfer/sec:     30.27KB
```
So async replication did affect write performance quite a bit — it went down to 320 RPS from inital 400 RPS.

### Resources usage while reading
```
CONTAINER ID   NAME                          CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
15faa3501e76   social-network_db_replica_1   355.66%   413.5MiB / 1.939GiB   20.83%    84.7MB / 22.6MB   115MB / 4.04GB    89
24718d51a7b1   social-network_db_1           17.35%    388.3MiB / 1.939GiB   19.56%    95.8MB / 185MB    44.3MB / 1.03GB   56
4b928805224f   social-network_db_proxy_1     39.27%    23.79MiB / 1.939GiB   1.20%     41.1MB / 51.1MB   14MB / 381kB      19
```

### Benchmark results
```
wrk --latency --timeout 10s -d 5m -t 3 -c 50 -s benchmarks/search.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   322.85ms  399.15ms   3.86s    84.83%
    Req/Sec    84.82     59.00   696.00     75.14%
  Latency Distribution
     50%  185.97ms
     75%  519.21ms
     90%  875.61ms
     99%    1.68s
  74638 requests in 5.00m, 149.77MB read
  Non-2xx or 3xx responses: 48
Requests/sec:    248.71
Transfer/sec:    511.04KB
```
Read performance was roughly the same as well but with way less load on main db instance.

## Semi-sync replication
Semi-sync replication was configured with following config on the main instance:
```
[mysqld]
server-id=1
binlog_format=ROW
log-bin
gtid_mode=ON
enforce_gtid_consistency=ON
rpl_semi_sync_master_enabled=ON
rpl_semi_sync_master_timeout=1000
```

And following config for replica instance:
```
[mysqld]
server-id=2
rpl_semi_sync_slave_enabled=1
gtid_mode=ON
enforce_gtid_consistency=ON
```

### Resources usage while idle
```
CONTAINER ID   NAME                            CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O        PIDS
7840087dffad   social-network_db_proxy_1       1.14%     17.91MiB / 1.939GiB   0.90%     49.7kB / 59.3kB   72.6MB / 381kB   19
51d8b693bab6   social-network_db_1             0.69%     357.1MiB / 1.939GiB   17.99%    53.4kB / 6.46MB   25.5MB / 281MB   42
136ce0c99645   social-network_db_replica_1     0.65%     378MiB / 1.939GiB     19.04%    3.24MB / 29.9kB   30.6MB / 327MB   41
f3f9b3c436ac   social-network_db_replica_2_1   0.55%     386.8MiB / 1.939GiB   19.48%    3.23MB / 22.5kB   55.8MB / 329MB   40
```

### Shutting down main instance test
Write load

```
> make REGISTER_DURATION=5m BENCH_N_CONN=50 bench-register
```

Resources usage while writing:
```
CONTAINER ID   NAME                            CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O        PIDS
7840087dffad   social-network_db_proxy_1       65.70%    21.9MiB / 1.939GiB    1.10%     14MB / 16.5MB     72.6MB / 381kB   19
51d8b693bab6   social-network_db_1             116.76%   427.8MiB / 1.939GiB   21.55%    7.49MB / 26.4MB   26.9MB / 384MB   109
136ce0c99645   social-network_db_replica_1     46.34%    372.4MiB / 1.939GiB   18.76%    10.6MB / 104kB    30.8MB / 384MB   41
f3f9b3c436ac   social-network_db_replica_2_1   44.87%    379.2MiB / 1.939GiB   19.10%    10.5MB / 96.4kB   56.2MB / 399MB   40
```

Main instance was stopped during benchmark
```
docker kill social-network_db_1 
```

Benchmark results:
```
wrk --latency --timeout 1s -d 5m -t 3 -c 50 -s benchmarks/register.lua http://localhost
Running 5m test @ http://localhost
  3 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   114.32ms   62.55ms 638.58ms   88.43%
    Req/Sec   156.39     47.95   393.00     77.12%
  Latency Distribution
     50%   94.94ms
     75%  118.63ms
     90%  184.91ms
     99%  346.30ms
  95277 requests in 3.69m, 8.89MB read
  Socket errors: connect 0, read 0, write 0, timeout 48
  Non-2xx or 3xx responses: 44
Requests/sec:    430.86
Transfer/sec:     41.15KB
```

Following snippet was added to register.lua in order to count successful `/register` requests
```lua
successful_requests = 0
response = function(status, headers, body)
    if status == 302 then
        successful_requests = successful_requests + 1
        print("thread #", thread_id, "performed", successful_requests, "requests")
    end
end
```

Log tail
```
thread #	1	performed	31665	requests
thread #	2	performed	31342	requests
thread #	0	performed	32226	requests
```

So (31665 + 31342 + 32226) = 95233 users must have been saved in at least one replica.

Replica instance 1:
```
mysql> select count(1) from profiles;
+----------+
| count(1) |
+----------+
|    95233 |
+----------+
1 row in set (0.03 sec)
```

Replica instance 2:
```
mysql> select count(1) from profiles;
+----------+
| count(1) |
+----------+
|    95233 |
+----------+
1 row in set (0.00 sec)
```

So none of successful requests were lost after shutting down master. 

Main instance was started up again and the number of rows was checked.
```
mysql> select count(1) from profiles;
+----------+
| count(1) |
+----------+
|    95237 |
+----------+
1 row in set (0.00 sec)
```
So there were 4 requests that made it to mysql binlog but were not confirmed for the client at the moment of shutting down. 
This means that if we had promoted any replica as the new main, it would not be safe to use that old main instance without trimming it's binlog.
This effect is explained in details in [the following post](https://percona.community/blog/2018/08/23/question-about-semi-synchronous-replication-answer-with-all-the-details/).