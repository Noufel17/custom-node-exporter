package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &Collector{}

// a promehteus.collector 
type Collector struct {
	latencyMetric *prometheus.GaugeVec
	packetLossMetric *prometheus.GaugeVec
	jitterMetric *prometheus.GaugeVec
	getMetrics             func() (Metrics, error)
}

func NewCollector(getMetrics func() (Metrics, error)) *Collector {
	return &Collector{
		latencyMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "network_latency",
				Help: "Network latency in milliseconds",
			},
			[]string{"destination"},
		),
		packetLossMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "network_packet_loss",
				Help: "Packet loss percentage",
			},
			[]string{"destination"},
		),
		jitterMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "network_jitter",
				Help: "Network jitter in milliseconds",
			},
			[]string{"destination"},
		),
		getMetrics:getMetrics,
	}
}


func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.latencyMetric.Describe(ch)
	c.packetLossMetric.Describe(ch)
	c.jitterMetric.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.latencyMetric.Collect(ch)
	c.packetLossMetric.Collect(ch)
	c.jitterMetric.Collect(ch)
}

func (c *Collector) Update() {
	metrics,err := c.getMetrics()
	if err != nil {
		panic(err)
	}
	c.latencyMetric.WithLabelValues(metrics.destination).Set(metrics.latency)
	c.packetLossMetric.WithLabelValues(metrics.destination).Set(metrics.packetLoss)
	c.jitterMetric.WithLabelValues(metrics.destination).Set(metrics.jitter)
}
