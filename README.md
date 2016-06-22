# Reliable Data Transmission Under High Latency And High Packet Loss Rate Network

A Reliable Data Transmission (rdt) bridge written in [Go](https://golang.org/). Described in [htc's lab document](https://github.com/FduSS/network-pj1/blob/master/README.md).

## Overview

![lab1](https://cloud.githubusercontent.com/assets/7262715/16257926/ea72a1a8-388d-11e6-90c3-bf751127956f.png)

**Alice** listens on `127.0.0.1:7000`, accepts a request and forwards it to **tunA** on `172.19.0.1:5555`.

**Bob** listens **tunB** on `172.20.0.1:5555`, accepts a request from **Alice** and forwards it to the **Web Server** on `127.0.0.1:8000`.

When **Web Server** accepts the request, it sends back the response reversely.

## How to run

**Step0.** To clone and run this repository you'll need [Git](https://git-scm.com/) and [golang environment](https://golang.org/doc/install) installed on your computer. From your command line:

```shell
# Clone this repository
git clone https://github.com/geeeeeeeeek/rdt.git
# Go into the repository
cd rdt
```

To continue working on the code, you can open the project in IntelliJ IDEA with this [plugin](https://plugins.jetbrains.com/plugin/5047).

**Step1.** [Set up](https://github.com/FduSS/network-pj1/blob/master/README.md#environment-setting-up) and run `transfer` to simulate a tunnel with latency and packet loss:

```shell
sudo ./network-pj1/Release/transfer-on-your-platform
```

**Step2.** Build the Go package:

```shell
go build -o lab1 ./src
```

There should be a binary file `lab1` generated in your working directory.

**Step3.** Run Alice and Bob:

```shell
./lab1 --bob
./lab1 --alice
```

Remember to run **Bob** ahead of **Alice**. Or the socket could not be initialized.

**Step 4.** Use `wget` or `axel` to test downloading performance (suppose your [nginx server](https://github.com/FduSS/network-pj1/blob/master/README.md#run-download-test) has been running on `127.0.0.1:8000` with a 100M file):

```shell
# Single-thread download
wget 127.0.0.1:7000/100M -O 100M.test
# Multi-thread download
axel -n 64 http://127.0.0.1:7000/100M -O 100M.test
```

You can use `diff` to verify the result.

**A sample test with all commands running.**

![image](https://cloud.githubusercontent.com/assets/7262715/16259333/b484423e-3894-11e6-831d-040da1988dfa.png)

### Implementation

So easy with golang. The code is less than 100 lines. 

Take **Alice** as an example, establish a TCP socket that listens `127.0.0.1:7000` and dials `172.19.0.1:5555`. Use goroutine to handle multiple connections.

```go
go HandleTransmission(A, B, "A to B")
go HandleTransmission(B, A, "B to A")
```

A goroutine is a lightweight thread managed by the Go runtime. go f(x, y, z) starts a new goroutine running f(x, y, z) The evaluation of f , x , y , and z happens in the current goroutine and the execution of f happens in the new goroutine.

### Performance

Test environment:

```yaml
os: Mac OS X 10.11.4
go: go1.6.2 darwin/amd64
axel: Axel version 2.5 (Darwin)
wget: GNU Wget 1.17 built on darwin15.0.0
```

Test result:

| Speed Limit | Delay | Packet Drop Rate | Threads | Overall Speed | Efficiency |
| :---------: | :---: | :--------------: | :-----: | :-----------: | :--------: |
|     --      |  --   |        --        |    1    |   580 Mbps    |     --     |
|     --      |  --   |        --        |   64    |   306 Mbps    |     --     |
|  100 Mbps   |  --   |        --        |    1    |    82 Mbps    |    82%     |
|  100 Mbps   |  --   |        --        |   64    |    95 Mbps    |    95%     |
|     --      | 10 ms |        --        |    1    |    47 Mbps    |     --     |
|     --      | 10 ms |        --        |   64    |   152 Mbps    |     --     |
|     --      |  --   |       10%        |    1    |   1.2 Mbps    |     --     |
|     --      |  --   |       10%        |   64    |   1.2 Mbps    |     --     |
|  100 Mbps   | 10 ms |       10%        |    1    |   0.6 Mbps    |    0.6%    |
|  100 Mbps   | 10 ms |       10%        |   64    |   0.5 Mbps    |    0.5%    |