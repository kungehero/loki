package wal

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type WatcherMetrics struct {
	recordsRead               *prometheus.CounterVec
	recordDecodeFails         *prometheus.CounterVec
	droppedWriteNotifications *prometheus.CounterVec
	segmentRead               *prometheus.CounterVec
	currentSegment            *prometheus.GaugeVec
	watchersRunning           *prometheus.GaugeVec
}

func NewWatcherMetrics(reg prometheus.Registerer) *WatcherMetrics {
	m := &WatcherMetrics{
		recordsRead: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "loki",
				Subsystem: "wal_watcher",
				Name:      "records_read_total",
				Help:      "Number of records read by the WAL watcher from the WAL.",
			},
			[]string{"id"},
		),
		recordDecodeFails: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "loki",
				Subsystem: "wal_watcher",
				Name:      "record_decode_failures_total",
				Help:      "Number of records read by the WAL watcher that resulted in an error when decoding.",
			},
			[]string{"id"},
		),
		droppedWriteNotifications: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "loki",
				Subsystem: "wal_watcher",
				Name:      "dropped_write_notifications_total",
				Help:      "Number of dropped write notifications due to having one already buffered.",
			},
			[]string{"id"},
		),
		segmentRead: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "loki",
				Subsystem: "wal_watcher",
				Name:      "segment_read_total",
				Help:      "Number of segment reads triggered by the backup timer firing.",
			},
			[]string{"id", "reason"},
		),
		currentSegment: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "loki",
				Subsystem: "wal_watcher",
				Name:      "current_segment",
				Help:      "Current segment the WAL watcher is reading records from.",
			},
			[]string{"id"},
		),
		watchersRunning: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "loki",
				Subsystem: "wal_watcher",
				Name:      "running",
				Help:      "Number of WAL watchers running.",
			},
			nil,
		),
	}

	// Collectors will be re-registered to registry if it's got reloaded
	// Reuse the old collectors instead of panicking out.
	if reg != nil {
		if err := reg.Register(m.recordsRead); err != nil {
			are := &prometheus.AlreadyRegisteredError{}
			if errors.As(err, are) {
				m.recordsRead = are.ExistingCollector.(*prometheus.CounterVec)
			}
		}
		if err := reg.Register(m.recordDecodeFails); err != nil {
			are := &prometheus.AlreadyRegisteredError{}
			if errors.As(err, are) {
				m.recordDecodeFails = are.ExistingCollector.(*prometheus.CounterVec)
			}
		}
		if err := reg.Register(m.droppedWriteNotifications); err != nil {
			are := &prometheus.AlreadyRegisteredError{}
			if errors.As(err, are) {
				m.droppedWriteNotifications = are.ExistingCollector.(*prometheus.CounterVec)
			}
		}
		if err := reg.Register(m.segmentRead); err != nil {
			are := &prometheus.AlreadyRegisteredError{}
			if errors.As(err, are) {
				m.segmentRead = are.ExistingCollector.(*prometheus.CounterVec)
			}
		}
		if err := reg.Register(m.currentSegment); err != nil {
			are := &prometheus.AlreadyRegisteredError{}
			if errors.As(err, are) {
				m.currentSegment = are.ExistingCollector.(*prometheus.GaugeVec)
			}
		}
		if err := reg.Register(m.watchersRunning); err != nil {
			are := &prometheus.AlreadyRegisteredError{}
			if errors.As(err, are) {
				m.watchersRunning = are.ExistingCollector.(*prometheus.GaugeVec)
			}
		}
	}

	return m
}
