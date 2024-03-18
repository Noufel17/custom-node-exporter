package exporter

import (
	"os/exec"
	"regexp"
	"strconv"
)

// metrics struct
type Metrics struct {
	latency float64
	packetLoss float64
	jitter float64
	destination string
}

func GetMetrics() (Metrics, error){

	var metrics Metrics

	// Sample network metrics using ping
	destinations := []string{"namla.cloud","google.com",}
	for _, dest := range destinations {
		cmd := exec.Command("ping", "-c", "10", dest)
		output, err := cmd.CombinedOutput()
		//fmt.Print(string(output))
		if err != nil {
			continue
		}

		// Use regular expressions to extract latency, packet loss, and jitter
		reLatency := regexp.MustCompile(`min/avg/max/mdev = [\d.]+/([\d.]+)/[\d.]+/[\d.]+ ms`)
		rePacketLoss := regexp.MustCompile(`(\d+)% packet loss`)
		reJitter := regexp.MustCompile(`= [\d.]+/[\d.]+/([\d.]+)/([\d.]+) ms`)

		matchLatency := reLatency.FindStringSubmatch(string(output))
		matchPacketLoss := rePacketLoss.FindStringSubmatch(string(output))
		matchJitter := reJitter.FindStringSubmatch(string(output))

		if len(matchLatency) < 2 || len(matchPacketLoss) < 2 || len(matchJitter) < 2 {
			continue
		}

		latency, err := strconv.ParseFloat(matchLatency[1], 64)
		if err != nil {
			continue
		}
		packetLoss, err := strconv.ParseFloat(matchPacketLoss[1], 64)
		if err != nil {
			continue
		}
		jitter, err := strconv.ParseFloat(matchJitter[1], 64)
		if err != nil {
			continue
		}
		metrics.latency = latency
		metrics.packetLoss  =packetLoss
		metrics.jitter  = jitter
		metrics.destination = dest
		
	}
	return metrics, nil
}