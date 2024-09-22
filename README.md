# Image Generator

## Purpose:

This project provides a RESTful API for generating images dynamically based on various parameters. The API allows users to specify the image dimensions, background color, text content, text color, and font size to create customized images.

## Features:

- Dynamic image generation: Creates images on-the-fly based on user-provided parameters.
- Customizable parameters: Allows users to specify dimensions, colors, text, and font size.
- Flexible API endpoints: Offers multiple endpoints for different image sizes and parameters.
- Performance optimization: Leverages efficient image generation techniques for optimal performance.


## Usage:

To generate an image, send an HTTP GET request to the appropriate API endpoint, providing the desired parameters in the URL query string. For example:

```
http://localhost:8080/img/800/250/1D7373/Hello%20World/FFFFFF/100
```

This request will generate an image with the following specifications:

- Width: 800 pixels
- Height: 250 pixels
- Background color: #1D7373
- Text: "Hello World"
- Text color: #FFFFFF
- Font size: 100

![alt text](./docs/imgs/example-hello-world.png)


```
http://localhost:8080/img/500/450
```

![alt text](./docs/imgs/example-simple.png)

## Performance:

The API has been designed with performance in mind. Benchmark tests using wrk have demonstrated high request rates and low latency, even under load.

## Additional Considerations:

Error handling: The API should implement appropriate error handling mechanisms to provide informative responses for invalid requests.
Security: Consider implementing security measures to protect against potential vulnerabilities, such as input validation and rate limiting.
Scalability: If the API is expected to handle a large number of requests, explore options for scaling and load balancing.

## Benchmarking:

```bash
$ wrk http://localhost:8080/img/100/100
Running 10s test @ http://localhost:8080/100/100
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.47ms   15.31ms 134.09ms   79.14%
    Req/Sec   226.44     45.48   310.00     72.50%
  4535 requests in 10.06s, 4.65MB read
Requests/sec:    450.64
Transfer/sec:    473.53KB
```

```bash
$ wrk http://localhost:8080/img/1/1
Running 10s test @ http://localhost:8080/1/1
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.57ms   10.71ms  85.48ms   71.96%
    Req/Sec   263.57     32.53   330.00     71.00%
  5280 requests in 10.07s, 3.52MB read
Requests/sec:    524.48
Transfer/sec:    358.53KB
```


```bash
$ wrk http://localhost:8080/ping
Running 10s test @ http://localhost:8080/ping
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   199.14us  634.76us  33.92ms   98.19%
    Req/Sec    30.93k     5.41k   49.58k    80.60%
  618533 requests in 10.10s, 70.79MB read
Requests/sec:  61235.01
Transfer/sec:      7.01MB
```