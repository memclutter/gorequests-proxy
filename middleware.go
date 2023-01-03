package gorequests_proxy

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"math/rand"
	"net/http"
	"net/url"
)

type Proxy struct{ Proxies []string }

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
		c = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	} else if t, ok := c.Transport.(*http.Transport); ok {
		t.Proxy = http.ProxyURL(proxyUrl)
	} else if t, ok := c.Transport.(*retryablehttp.RoundTripper); ok {
		if t.Client.HTTPClient.Transport == nil {
			t.Client.HTTPClient.Transport = &http.Transport{}
		}
		t.Client.HTTPClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
	} else {
		return c, fmt.Errorf("unsupported http transport: %T", c.Transport)
	}
	return c, nil
}
