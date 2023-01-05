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
	Proxies []string
}

func (m *Proxy) ClientOverride(c *http.Client) (*http.Client, error) {
	if len(m.Proxies) == 0 {
		return nil, fmt.Errorf("nil proxies")
	}
	proxy := m.Proxies[rand.Intn(len(m.Proxies))]
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("error parse proxy url %s: %v", proxy, err)
	}
	if c.Transport == nil {
		c = &http.Client{}
		if strings.Contains(proxyUrl.Scheme, "socks") {
			c.Transport = &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return socks.Dial(proxyUrl.String())(network, addr)
				},
			}
		} else {
			c = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
		}
	} else if t, ok := c.Transport.(*http.Transport); ok {
		if strings.Contains(proxyUrl.Scheme, "socks") {
			t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return socks.Dial(proxyUrl.String())(network, addr)
			}
		} else {
			t.Proxy = http.ProxyURL(proxyUrl)
		}
	} else if t, ok := c.Transport.(*retryablehttp.RoundTripper); ok {
		if t.Client.HTTPClient.Transport == nil {
			t.Client.HTTPClient.Transport = &http.Transport{}
		}
		if strings.Contains(proxyUrl.Scheme, "socks") {
			t.Client.HTTPClient.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return socks.Dial(proxyUrl.String())(network, addr)
			}
		} else {
			t.Client.HTTPClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
		}
	} else {
		return c, fmt.Errorf("unsupported http transport: %T", c.Transport)
	}
	return c, nil
}
