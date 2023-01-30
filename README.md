# parkDNS
A Golang-based DNS server for unified domain management.

## Install

```
$ go get github.com/miekg/dns
$ go build parkdns.go
$ sudo ./parkdns
```

## records.json sample

```
{
  "records": [
    {
      "type": 1,
      "ttl": 3600,
      "data": "192.168.1.1"
    },
    {
      "type": 15,
      "ttl": 3600,
      "data": "mail.example.com.",
      "priority": 10
    },
    {
      "type": 15,
      "ttl": 3600,
      "data": "mail2.example.com.",
      "priority": 20
    }
  ]
}
```
