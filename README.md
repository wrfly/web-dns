# Web Dns

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

- [x] dns lib
- [ ] cacher
- [ ] docker-compose
- [ ] metrics and debug
- [ ] performance test
- [ ] rate limit
- [ ] hijack
- [x] api list