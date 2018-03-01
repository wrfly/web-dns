# Web Dns

[![Build Status](https://travis-ci.org/wrfly/web-dns.svg?branch=master)](https://travis-ci.org/wrfly/web-dns)


curl https://dns.kfd.me/www.google.com

## API List

### Default

```text
/www.google.com
```

### With Type

```text
/www.google.com/AAAA
/www.google.com/A
/www.google.com/MX
```

### Reurn Json

```text
/www.google.com/json
/www.google.com/MX/json
/www.google.com/AAAA/json
```

## TODO

- [x] it works
- [x] dns lib
- [x] cacher
    - [x] mem
    - [x] redis
    - [x] blot
- [ ] metrics and debug
- [ ] performance test
- [x] rate limit(per IP)
- [ ] hijack
- [x] api list