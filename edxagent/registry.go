/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import (
	"golang.zabbix.com/agent2/pkg/watch"
	"golang.zabbix.com/sdk/errs"
	"golang.zabbix.com/sdk/plugin"
)

var registry = make(map[string]*Plugin)

func Register(name string, metricSets []*MetricSet, options ...Option) error {
	var err error
	var metricSetSource *MetricSetSource

	if _, ok := registry[name]; ok {
		return errs.New("plugin " + name + " already registered")
	}

	p := &Plugin{}
	registry[name] = p

	p.name = name
	p.SetFlushInterval(0)

	for _, option := range options {
		option.set(p)
	}

	manager := watch.NewManager(p)
	p.SetManager(manager)

	metrics := make([]string, 0)

	usedKeys := make(map[string]bool)

	for _, ms := range metricSets {
		metricSetSource = NewMetricSetSource(manager, ms.Collect)
		if ms.Collect != nil {
			p.AddMetricSource(metricSetSource)
		}

		if ms.Update != nil {
			metricSetSource.SetWatchUpdate(ms.Update)
		}

		for _, m := range ms.Metrics {
			if len(m.Key) == 0 {
				return errs.New("metric key cannot be empty")
			}

			if len(m.Description) == 0 {
				return errs.New("metric description cannot be empty for key: " + m.Key)
			}

			if m.Description[len(m.Description)-1] != '.' {
				m.Description += "."
			}

			if _, exists := usedKeys[m.Key]; exists {
				return errs.New("duplicate metric key: " + m.Key)
			}

			if m.Export != nil {
				p.AddExporter(m.Key, m.Export)
			} else if metricSetSource != nil && m.NewWatchFilter != nil {
				metricSetSource.AddMetric(m.Key, m.NewWatchFilter)
				p.AssignKeyToMetricSource(m.Key, metricSetSource)
			} else {
				return errs.New("no export function or filter defined for metric key: " + m.Key)
			}

			metrics = append(metrics, m.Key, m.Description)
		}

		/*
			if metricSetSource != nil && !metricSetSource.HasMetrics() {
					return errs.New("Plugin " + name + ": Collector defined, but no suitable metrics in " + ms.Name)
			}
		*/
	}

	err = plugin.RegisterMetrics(p, name, metrics...)
	return err
}
