/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import (
	"errors"

	"golang.zabbix.com/agent2/pkg/itemutil"
	"golang.zabbix.com/agent2/pkg/watch"
)

type Publisher interface {
	watch.EventSource
	Publish(any)
}

type MetricSetSource struct {
	manager      *watch.Manager
	newFilterMap map[string]WatchFilterFn
	collector    CollectFn
	watchUpdate  UpdateFn
}

func NewMetricSetSource(manager *watch.Manager, fn CollectFn) *MetricSetSource {
	return &MetricSetSource{
		manager:      manager,
		collector:    fn,
		newFilterMap: make(map[string]WatchFilterFn),
	}
}

func (e *MetricSetSource) AddMetric(key string, fn WatchFilterFn) {
	if e.newFilterMap == nil {
		e.newFilterMap = make(map[string]WatchFilterFn)
	}
	e.newFilterMap[key] = fn
}

func (e *MetricSetSource) SetWatchUpdate(fn UpdateFn) {
	e.watchUpdate = fn
}

func (e *MetricSetSource) Initialize() (err error) {
	return nil
}

func (e *MetricSetSource) Release() {
}

func (e *MetricSetSource) NewFilter(key string) (watch.EventFilter, error) {
	if e.newFilterMap != nil {
		name, params, err := itemutil.ParseKey(key)
		if err != nil {
			return nil, err
		}
		if f, ok := e.newFilterMap[name]; ok {
			return f(params)
		} else {
			return nil, errors.New("Unknown event key: " + name)
		}
	}
	return nil, errors.New("no metrics registered for this set")
}

func (e *MetricSetSource) Publish(data interface{}) {
	if e.manager != nil {
		e.manager.Lock()
		e.manager.Notify(e, data)
		e.manager.Unlock()
	}
}

func (e *MetricSetSource) HasMetrics() bool {
	return len(e.newFilterMap) > 0
}
