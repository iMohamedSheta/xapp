package utils

import (
	"fmt"
	"net"
	"time"
)

// queryNTPTime queries multiple NTP servers using raw UDP and returns consensus time.
// No external dependencies — uses only the standard net package.
func QueryNTPTime() (time.Time, bool) {
	servers := []string{
		"time.google.com:123",
		"time.cloudflare.com:123",
		"pool.ntp.org:123",
		"time.windows.com:123",
	}

	type result struct {
		t   time.Time
		err error
	}

	results := make(chan result, len(servers))
	for _, srv := range servers {
		srv := srv
		SafeGo(func() {
			t, err := ntpQuery(srv)
			results <- result{t, err}
		})
	}

	var times []time.Time
	deadline := time.After(5 * time.Second)
	for i := 0; i < len(servers); i++ {
		select {
		case r := <-results:
			if r.err == nil {
				times = append(times, r.t)
			}
		case <-deadline:
		}
	}

	if len(times) < 2 {
		return time.Time{}, false
	}

	for _, t := range times[1:] {
		diff := times[0].Sub(t)
		if diff < 0 {
			diff = -diff
		}
		if diff > 10*time.Second {
			return time.Time{}, false // inconsistent responses
		}
	}

	return times[0], true
}

// ntpQuery sends a single NTP request to the given server and returns the time.
// NTP packet is 48 bytes. We only need the transmit timestamp at bytes 40-47.
func ntpQuery(server string) (time.Time, error) {
	conn, err := net.DialTimeout("udp", server, 5*time.Second)
	if err != nil {
		return time.Time{}, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// NTP request packet — 48 bytes
	// First byte: LI=0, VN=4, Mode=3 (client)
	req := make([]byte, 48)
	req[0] = 0b00100011 // LI=0, Version=4, Mode=3

	if _, err := conn.Write(req); err != nil {
		return time.Time{}, err
	}

	resp := make([]byte, 48)
	if _, err := conn.Read(resp); err != nil {
		return time.Time{}, err
	}

	// Transmit timestamp is at bytes 40-47
	// NTP timestamp = seconds since Jan 1, 1900
	// Unix timestamp = seconds since Jan 1, 1970
	// Difference = 70 years in seconds
	const ntpEpochOffset = 2208988800

	secs := uint64(resp[40])<<24 | uint64(resp[41])<<16 | uint64(resp[42])<<8 | uint64(resp[43])
	if secs == 0 {
		return time.Time{}, fmt.Errorf("invalid NTP response")
	}

	unixSecs := int64(secs) - ntpEpochOffset
	return time.Unix(unixSecs, 0), nil
}
