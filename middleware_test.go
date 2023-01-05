package gorequests_proxy

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ProxySuite struct {
	suite.Suite
}

func (suite ProxySuite) TestEmptyProxies() {
	// Data
	client := &http.Client{}
	var proxyUrl []string

	// Run tests
	_, err := (&Proxy{Proxies: proxyUrl}).ClientOverride(client)

	// Assertions
	assert.Error(suite.T(), err, "should be throw error")
}

func (suite ProxySuite) TestInvalidUrl() {
	// Data
	client := &http.Client{}
	proxyUrl := []string{string([]byte{0x7f})}

	// Run tests
	_, err := (&Proxy{Proxies: proxyUrl}).ClientOverride(client)

	// Assertions
	assert.Error(suite.T(), err, "should be throw error")
}

func (suite ProxySuite) TestNoTransportHttpProxies() {
	tests := []struct {
		proxyUrl []string
	}{
		{
			proxyUrl: []string{"http://127.0.0.1:8888"},
		},
		{
			proxyUrl: []string{"https://127.0.0.1:8888"},
		},
	}

	for _, test := range tests {
		suite.Run(test.proxyUrl[0], func() {
			// Data
			client := &http.Client{}

			// Run tests
			overrideClient, err := (&Proxy{Proxies: test.proxyUrl}).ClientOverride(client)

			// Asserts
			assert.NoError(suite.T(), err, "should be no errors")
			assert.NotNilf(suite.T(), overrideClient.Transport, "transport must be set")
			assert.NotNilf(suite.T(), overrideClient.Transport.(*http.Transport).Proxy, "should be set Proxy func")
		})
	}
}

func (suite ProxySuite) TestNoTransportSocksProxies() {
	tests := []struct {
		proxyUrl []string
	}{
		{
			proxyUrl: []string{"socks4://127.0.0.1:8888"},
		},
		{
			proxyUrl: []string{"socks5://127.0.0.1:8888"},
		},
	}

	for _, test := range tests {
		suite.Run(test.proxyUrl[0], func() {
			// Data
			client := &http.Client{}

			// Run tests
			overrideClient, err := (&Proxy{Proxies: test.proxyUrl}).ClientOverride(client)

			// Asserts
			assert.NoError(suite.T(), err, "should be no errors")
			assert.NotNilf(suite.T(), overrideClient.Transport, "transport must be set")
			assert.NotNilf(suite.T(), overrideClient.Transport.(*http.Transport).DialContext, "should be set DialContext func")
		})
	}
}

func (suite ProxySuite) TestWithTransportHttpProxies() {
	tests := []struct {
		proxyUrl []string
	}{
		{
			proxyUrl: []string{"http://127.0.0.1:8888"},
		},
		{
			proxyUrl: []string{"https://127.0.0.1:8888"},
		},
	}

	for _, test := range tests {
		suite.Run(test.proxyUrl[0], func() {
			// Data
			client := &http.Client{Transport: &http.Transport{}}

			// Run tests
			overrideClient, err := (&Proxy{Proxies: test.proxyUrl}).ClientOverride(client)

			// Asserts
			assert.NoError(suite.T(), err, "should be no errors")
			assert.NotNilf(suite.T(), overrideClient.Transport.(*http.Transport).Proxy, "should be set Proxy func")
		})
	}
}

func (suite ProxySuite) TestWithTransportSocksProxies() {
	tests := []struct {
		proxyUrl []string
	}{
		{
			proxyUrl: []string{"socks4://127.0.0.1:8888"},
		},
		{
			proxyUrl: []string{"socks5://127.0.0.1:8888"},
		},
	}

	for _, test := range tests {
		suite.Run(test.proxyUrl[0], func() {
			// Data
			client := &http.Client{Transport: &http.Transport{}}

			// Run tests
			overrideClient, err := (&Proxy{Proxies: test.proxyUrl}).ClientOverride(client)

			// Asserts
			assert.NoError(suite.T(), err, "should be no errors")
			assert.NotNilf(suite.T(), overrideClient.Transport.(*http.Transport).DialContext, "should be set DialContext func")
		})
	}
}

func (suite ProxySuite) TestWithRetryableTransportHttpProxies() {
	tests := []struct {
		proxyUrl []string
	}{
		{
			proxyUrl: []string{"http://127.0.0.1:8888"},
		},
		{
			proxyUrl: []string{"https://127.0.0.1:8888"},
		},
	}

	for _, test := range tests {
		suite.Run(test.proxyUrl[0], func() {
			// Data
			client := &http.Client{}
			rc := retryablehttp.NewClient()
			rc.HTTPClient = client

			// Run tests
			overrideClient, err := (&Proxy{Proxies: test.proxyUrl}).ClientOverride(rc.StandardClient())

			// Asserts
			assert.NoError(suite.T(), err, "should be no errors")
			assert.NotNilf(suite.T(), overrideClient.Transport, "transport must be set")
			assert.NotNilf(suite.T(), overrideClient.Transport.(*retryablehttp.RoundTripper).Client.HTTPClient.Transport.(*http.Transport).Proxy, "should be set Proxy func")
		})
	}
}

func (suite ProxySuite) TestWithRetryableTransportSocksProxies() {
	tests := []struct {
		proxyUrl []string
	}{
		{
			proxyUrl: []string{"socks4://127.0.0.1:8888"},
		},
		{
			proxyUrl: []string{"socks5://127.0.0.1:8888"},
		},
	}

	for _, test := range tests {
		suite.Run(test.proxyUrl[0], func() {
			// Data
			client := &http.Client{}
			rc := retryablehttp.NewClient()
			rc.HTTPClient = client

			// Run tests
			overrideClient, err := (&Proxy{Proxies: test.proxyUrl}).ClientOverride(rc.StandardClient())

			// Asserts
			assert.NoError(suite.T(), err, "should be no errors")
			assert.NotNilf(suite.T(), overrideClient.Transport, "transport must be set")
			assert.NotNilf(suite.T(), overrideClient.Transport.(*retryablehttp.RoundTripper).Client.HTTPClient.Transport.(*http.Transport).DialContext, "should be set DialContext func")
		})
	}
}

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxySuite))
}
