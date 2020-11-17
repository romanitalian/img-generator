# img-generator

[http://localhost:8080](http://localhost:8080)

[http://localhost:8080/600/200](http://localhost:8080/600/200)

[http://localhost:8080/1500/1700/1D7373/Привет с Марса/FFFFFF/100](http://localhost:8080/1500/1700/1D7373/%D0%9F%D1%80%D0%B8%D0%B2%D0%B5%D1%82%20%D1%81%20%D0%9C%D0%B0%D1%80%D1%81%D0%B0/FFFFFF/100)

[http://localhost:8080/1500/1700/1D7373/Привет с Марса/FFFFFF/100](http://localhost:8080/1500/1700/1D7373/%D0%9F%D1%80%D0%B8%D0%B2%D0%B5%D1%82%20%D1%81%20%D0%9C%D0%B0%D1%80%D1%81%D0%B0/FFFFFF/100)

[http://localhost:8080/600/150/1D7373/Привет с Марса/FFFFFF](http://localhost:8080/600/150/1D7373/%D0%9F%D1%80%D0%B8%D0%B2%D0%B5%D1%82%20%D1%81%20%D0%9C%D0%B0%D1%80%D1%81%D0%B0/FFFFFF)

[http://localhost:8080/600/150/1D7373/Lorem/FFFFFF](http://localhost:8080/600/150/1D7373/Lorem/FFFFFF)


```shell script
 wrk http://localhost:8080/100/100
Running 10s test @ http://localhost:8080/100/100
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.47ms   15.31ms 134.09ms   79.14%
    Req/Sec   226.44     45.48   310.00     72.50%
  4535 requests in 10.06s, 4.65MB read
Requests/sec:    450.64
Transfer/sec:    473.53KB
```

```shell script
 wrk http://localhost:8080/1/1
Running 10s test @ http://localhost:8080/1/1
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.57ms   10.71ms  85.48ms   71.96%
    Req/Sec   263.57     32.53   330.00     71.00%
  5280 requests in 10.07s, 3.52MB read
Requests/sec:    524.48
Transfer/sec:    358.53KB
```


```shell script
wrk http://localhost:8080/ping
Running 10s test @ http://localhost:8080/ping
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   199.14us  634.76us  33.92ms   98.19%
    Req/Sec    30.93k     5.41k   49.58k    80.60%
  618533 requests in 10.10s, 70.79MB read
Requests/sec:  61235.01
Transfer/sec:      7.01MB
```
