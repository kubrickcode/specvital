package ratelimit

import (
	"testing"
	"time"
)

func TestIPRateLimiter_Allow(t *testing.T) {
	limiter := NewIPRateLimiter(3, 100*time.Millisecond)
	defer func() { _ = limiter.Close() }()

	ip := "192.168.1.1"

	if !limiter.Allow(ip) {
		t.Error("first request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("second request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("third request should be allowed")
	}

	if limiter.Allow(ip) {
		t.Error("fourth request should be rate limited")
	}
}

func TestIPRateLimiter_WindowReset(t *testing.T) {
	limiter := NewIPRateLimiter(2, 50*time.Millisecond)
	defer func() { _ = limiter.Close() }()

	ip := "192.168.1.1"

	if !limiter.Allow(ip) {
		t.Error("first request should be allowed")
	}
	if !limiter.Allow(ip) {
		t.Error("second request should be allowed")
	}
	if limiter.Allow(ip) {
		t.Error("third request should be rate limited")
	}

	time.Sleep(60 * time.Millisecond)

	if !limiter.Allow(ip) {
		t.Error("request after window reset should be allowed")
	}
}

func TestIPRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewIPRateLimiter(1, 100*time.Millisecond)
	defer func() { _ = limiter.Close() }()

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	if !limiter.Allow(ip1) {
		t.Error("first IP first request should be allowed")
	}
	if limiter.Allow(ip1) {
		t.Error("first IP second request should be rate limited")
	}

	if !limiter.Allow(ip2) {
		t.Error("second IP should be allowed (separate limit)")
	}
}
