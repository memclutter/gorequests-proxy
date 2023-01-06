package gorequests_proxy

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"h12.io/socks"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Proxy struct {
	Proxies    []string
	AllowEmpty bool
}

func (m *Proxy) ClientOverride(c *http.Client) (*http.Client, error) {
	if len(m.Proxies) == 0 {
		if m.AllowEmpty {
			return c, nil
		}
		return nil, fmt.Errorf("nil proxies")
	}
	proxy := m.Proxies[rand.Intn(len(m.Proxies))]
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("error parse proxy url %s: %v", proxy, err)
	}

	var dialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error) = nil
	var proxyFunc func(*http.Request) (*url.URL, error) = nil

	if strings.Contains(proxyUrl.Scheme, "socks") {
		dialContextFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return socks.Dial(proxyUrl.String())(network, addr)
		}
	} else {
		proxyFunc = http.ProxyURL(proxyUrl)
	}

	// Preset default transport
	if c.Transport == nil {
		c.Transport = &http.Transport{}
	}

	switch t := c.Transport.(type) {
	case *http.Transport:
		t.DialContext = dialContextFunc
		t.Proxy = proxyFunc
	case *retryablehttp.RoundTripper:
		if t.Client.HTTPClient.Transport == nil {
			t.Client.HTTPClient.Transport = &http.Transport{}
		}
		t.Client.HTTPClient.Transport.(*http.Transport).DialContext = dialContextFunc
		t.Client.HTTPClient.Transport.(*http.Transport).Proxy = proxyFunc
	default:
		return nil, fmt.Errorf("unsupported http transport type '%T'", t)
	}

	return c, nil
}
