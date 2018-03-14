# Web Dns

[![Build Status](https://travis-ci.org/wrfly/web-dns.svg?branch=master)](https://travis-ci.org/wrfly/web-dns)

curl https://dns.kfd.me/www.google.com

## Usage

```txt
NAME:
   web-dns - Query domain via HTTP(S)

USAGE:
   web-dns [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command
GLOBAL OPTIONS:
   --port value        port to listen on (default: 8080) [$WDC_PORT]
   --dns value         dns servers (default: "8.8.8.8:53", "8.8.4.4:53") [$WDC_DNS]  --timeout value     dig timeout (second) (default: 1) [$WDC_TIMEOUT]
   --rate value        rate of requests per minute per IP (default: 1000) [$WDC_RATE]
   --cache value       cache type: mem|redis|bolt (default: "mem") [$WDC_CACHE]
   --redis-addr value  this flag is used for redis cacher (default: "localhost:6379") [$WDC_REDIS_ADDR]
   --black-list value  blacklist of clients (default: "8.8.8.8", "4.4.4.4") [$WDC_BLACK_LIST]
   --hosts value       hijack hosts file path [$WDC_HOSTS]
   --debug             debug log-level, metrics and pprof debug (default: false) [$WDC_DEBUG]
   --debug-port value  server debug port (default: 8081) [$WDC_DEBUG_PORT]
   --help, -h          show help (default: false)
   --version, -v       print the version (default: false)
```

## API List

### Default

```sh
curl https://dns.kfd.me/www.google.com
216.58.199.228
```

### With Type

```sh
curl https://dns.kfd.me/www.google.com/AAAA
2404:6800:400a:807::2004

curl https://dns.kfd.me/github.com/A
192.30.253.112

curl https://dns.kfd.me/github.com/MX
aspmx.l.google.com.

curl https://dns.kfd.me/github.com/NS
ns4.p16.dynect.net.
```

### Reurn Json

```sh
curl https://dns.kfd.me/github.io/json -s | python -mjson.tool
[
    {
        "host": "151.101.1.147",
        "ttl": 600,
        "type": "A"
    },
    {
        "host": "151.101.129.147",
        "ttl": 600,
        "type": "A"
    },
    {
        "host": "151.101.65.147",
        "ttl": 600,
        "type": "A"
    },
    {
        "host": "151.101.193.147",
        "ttl": 600,
        "type": "A"
    }
]

curl https://dns.kfd.me/github.com/MX/json -s | python -mjson.tool
[
    {
        "host": "ALT1.ASPMX.L.GOOGLE.COM.",
        "ttl": 3600,
        "type": "MX"
    },
    {
        "host": "ASPMX.L.GOOGLE.COM.",
        "ttl": 3600,
        "type": "MX"
    },
    {
        "host": "ALT4.ASPMX.L.GOOGLE.COM.",
        "ttl": 3600,
        "type": "MX"
    },
    {
        "host": "ALT2.ASPMX.L.GOOGLE.COM.",
        "ttl": 3600,
        "type": "MX"
    },
    {
        "host": "ALT3.ASPMX.L.GOOGLE.COM.",
        "ttl": 3600,
        "type": "MX"
    }
]

curl https://dns.kfd.me/www.google.com/AAAA/json -s | python -mjson.tool
[
    {
        "host": "2404:6800:4008:800::2004",
        "ttl": 36395,
        "type": "AAAA"
    }
]
```

## TODO

- [x] it works
- [x] dns lib
- [x] cacher
    - [x] mem
    - [x] redis
    - [x] blot
- [x] metrics and debug
- [x] rate limit(per IP)
- [x] hijack
- [x] api list