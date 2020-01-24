/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ghmetrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// ghTokenUntilResetGaugeVec provides the 'github_token_reset' gauge that
// enables keeping track of GitHub reset times.
var ghTokenUntilResetGaugeVec = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "github_token_reset",
		Help: "Last reported GitHub token reset time.",
	},
	[]string{"token_hash", "api_version"},
)

// ghTokenUsageGaugeVec provides the 'github_token_usage' gauge that
// enables keeping track of GitHub calls and quotas.
var ghTokenUsageGaugeVec = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "github_token_usage",
		Help: "How many GitHub token requets are remaining for the current hour.",
	},
	[]string{"token_hash", "api_version"},
)

// ghRequestDurationHistVec provides the 'github_request_duration' histogram that keeps track
// of the duration of GitHub requests by API path.
var ghRequestDurationHistVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "github_request_duration",
		Help:    "GitHub request duration by API path.",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	},
	[]string{"token_hash", "path", "status", "user_agent"},
)

// cacheCounter provides the 'ghcache_responses' counter vec that is indexed
// by the cache response mode.
var cacheCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ghcache_responses",
		Help: "How many cache responses of each cache response mode there are.",
	},
	[]string{"mode", "path"},
)

// timeoutDuration provides the 'github_request_timeouts' histogram that keeps
// track of the timeouts of GitHub requests by API path.
var timeoutDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "github_request_timeouts",
		Help:    "GitHub request timeout by API path.",
		Buckets: []float64{45, 60, 90, 120, 300},
	},
	[]string{"token_hash", "path", "user_agent"},
)

var muxTokenUsage, muxRequestMetrics sync.Mutex
var lastGitHubResponse time.Time

func init() {
	prometheus.MustRegister(ghTokenUntilResetGaugeVec)
	prometheus.MustRegister(ghTokenUsageGaugeVec)
	prometheus.MustRegister(ghRequestDurationHistVec)
	prometheus.MustRegister(cacheCounter)
	prometheus.MustRegister(timeoutDuration)
}

// CollectGitHubTokenMetrics publishes the rate limits of the github api to
// `github_token_usage` as well as `github_token_reset` on prometheus.
func CollectGitHubTokenMetrics(tokenHash, apiVersion string, headers http.Header, reqStartTime, responseTime time.Time) {
	remaining := headers.Get("X-RateLimit-Remaining")
	timeUntilReset := timestampStringToTime(headers.Get("X-RateLimit-Reset"))
	durationUntilReset := timeUntilReset.Sub(reqStartTime)

	remainingFloat, err := strconv.ParseFloat(remaining, 64)
	if err != nil {
		logrus.WithError(err).Infof("Couldn't convert number of remaining token requests into gauge value (float)")
	}
	if remainingFloat == 0 {
		logrus.WithFields(logrus.Fields{
			"header":     remaining,
			"user-agent": headers.Get("User-Agent"),
		}).Debug("Parsed GitHub header as indicating no remaining rate-limit.")
	}

	muxTokenUsage.Lock()
	isAfter := lastGitHubResponse.After(responseTime)
	if !isAfter {
		lastGitHubResponse = responseTime
	}
	muxTokenUsage.Unlock()
	if isAfter {
		logrus.WithField("last-github-response", lastGitHubResponse).WithField("response-time", responseTime).Debug("Previously pushed metrics of a newer response, skipping old metrics")
	} else {
		ghTokenUntilResetGaugeVec.With(prometheus.Labels{"token_hash": tokenHash, "api_version": apiVersion}).Set(float64(durationUntilReset.Nanoseconds()))
		ghTokenUsageGaugeVec.With(prometheus.Labels{"token_hash": tokenHash, "api_version": apiVersion}).Set(remainingFloat)
	}
}

// CollectGitHubRequestMetrics publishes the number of requests by API path to
// `github_requests` on prometheus.
func CollectGitHubRequestMetrics(tokenHash, path, statusCode, userAgent string, roundTripTime float64) {
	ghRequestDurationHistVec.With(prometheus.Labels{"token_hash": tokenHash, "path": simplifier.Simplify(path), "status": statusCode, "user_agent": userAgent}).Observe(roundTripTime)
}

// timestampStringToTime takes a unix timestamp and returns a `time.Time`
// from the given time.
func timestampStringToTime(tstamp string) time.Time {
	timestamp, err := strconv.ParseInt(tstamp, 10, 64)
	if err != nil {
		logrus.WithField("timestamp", tstamp).Info("Couldn't convert unix timestamp")
	}
	return time.Unix(timestamp, 0)
}

// CollectCacheRequestMetrics records a cache outcome for a specific path
func CollectCacheRequestMetrics(mode, path string) {
	cacheCounter.With(prometheus.Labels{"mode": mode, "path": simplifier.Simplify(path)}).Inc()
}

// CollectRequestTimeoutMetrics publishes the duration of timed-out requests by
// API path to 'github_request_timeouts' on prometheus.
func CollectRequestTimeoutMetrics(tokenHash, path, userAgent string, reqStartTime, responseTime time.Time) {
	timeoutDuration.With(prometheus.Labels{"token_hash": tokenHash, "path": simplifier.Simplify(path), "user_agent": userAgent}).Observe(float64(responseTime.Sub(reqStartTime).Seconds()))
}