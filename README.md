# gorequests-proxy

[![Go](https://github.com/memclutter/gorequests-proxy/actions/workflows/go.yml/badge.svg)](https://github.com/memclutter/gorequests-proxy/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/memclutter/gorequests-proxy/branch/main/graph/badge.svg?token=1IWTNCLCAQ)](https://codecov.io/gh/memclutter/gorequests-proxy)

Proxy middleware for [memclutter/gorequests](https://github.com/memclutter/gorequests).

## Install

Installation using the go package system

```shell
go get github.com/memclutter/gorequets-proxy
```

## Use

To use, pass to the `Use()` method of the `RequestInstance`

```go
package main

import (
	"github.com/memclutter/gorequests"
	"github.com/memclutter/gorequests-proxy"
	"time"
)

func main() {
	proxy := &gorequests_proxy.Proxy{Proxies: []string{"http://example.com:5554", "http://user:pass@secure.example.com:33333"}}
	err := gorequests.Get("http://example.com").Use(proxy).Exec()
	// ...
}
```