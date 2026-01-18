/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import (
	"context"
	"time"

	"golang.zabbix.com/agent2/pkg/itemutil"
	"golang.zabbix.com/agent2/pkg/watch"
	"golang.zabbix.com/sdk/log"
	"golang.zabbix.com/sdk/plugin"
)

const DefaultFlushInterval = 120

var dSource = &dummySource{}

type ConfigureFn func(p *Plugin, global *GlobalOptions, options interface{})
type ValidateFn func(options interface{}) error

type Plugin struct {
	plugin.Base
	manager *watch.Manager

	name          string
	FlushInterval int

	Context    context.Context
	CancelFunc context.CancelFunc

	ConfigureFunc ConfigureFn
	ValidateFunc  ValidateFn

	exportByKey     map[string]ExportFn
	metricSetSource []*MetricSetSource
	eventSources    map[string]*MetricSetSource
}

type Logger interface {
	log.Logger
}

func (p *Plugin) AddMetricSource(metricSetSource *MetricSetSource) {
	if p.metricSetSource == nil {
		p.metricSetSource = make([]*MetricSetSource, 0)
	}
	p.metricSetSource = append(p.metricSetSource, metricSetSource)
}

func (p *Plugin) AssignKeyToMetricSource(key string, metricSetSource *MetricSetSource) {
	if p.eventSources == nil {
		p.eventSources = make(map[string]*MetricSetSource)
	}
	p.eventSources[key] = metricSetSource
}

func (p *Plugin) SetFlushInterval(interval int) {
	if interval > 0 {
		p.FlushInterval = interval
	} else {
		p.FlushInterval = DefaultFlushInterval
	}
}

func (p *Plugin) AddExporter(key string, fn ExportFn) {
	if p.exportByKey == nil {
		p.exportByKey = make(map[string]ExportFn)
	}
	p.exportByKey[key] = fn
}

func (p *Plugin) SetManager(manager *watch.Manager) {
	p.manager = manager
}

func (p *Plugin) Start() {
	p.Debugf("Starting " + p.name + " plugin")
	p.Context, p.CancelFunc = context.WithCancel(context.Background())

	for _, ms := range p.metricSetSource {
		go ms.collector(p.Context, ms)
	}

	go p.Flush()
}

func (p *Plugin) Stop() {
	p.Debugf("Stoping " + p.name + " plugin")
	p.CancelFunc()
}

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	if exporter, exists := p.exportByKey[key]; exists {
		return exporter(params)
	}
	return nil, nil
}

func (p *Plugin) Watch(items []*plugin.Item, ctx plugin.ContextProvider) {
	metricItemsMap := make(map[*MetricSetSource][]WatchKey)

	for _, item := range items {
		key, params, err := itemutil.ParseKey(item.Key)
		if err != nil {
			p.Errf("Failed to parse item key %s: %v", item.Key, err)
			continue
		} else {
			metricItemsMap[p.eventSources[key]] = append(metricItemsMap[p.eventSources[key]], WatchKey{Key: key, Params: params})
		}
	}

	for ms, watchItems := range metricItemsMap {
		if ms != nil && ms.watchUpdate != nil {
			ms.watchUpdate(watchItems)
		}
	}

	p.manager.Lock()
	p.manager.Update(ctx.ClientID(), ctx.Output(), items)
	p.manager.Unlock()
}

func (p *Plugin) EventSourceByKey(key string) (watch.EventSource, error) {
	p.Debugf("Event key: %v", key)

	if p.eventSources == nil {
		p.eventSources = make(map[string]*MetricSetSource)
	}

	name, _, err := itemutil.ParseKey(key)
	if err != nil {
		return nil, err
	}

	if es, exists := p.eventSources[name]; exists {
		return es, nil
	}

	return dSource, nil
}

func (p *Plugin) Flush() {
	ticker := time.NewTicker(time.Duration(p.FlushInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.Context.Done():
			return
		case <-ticker.C:
			p.manager.Lock()
			for _, es := range p.metricSetSource {
				p.manager.Flush(es)
			}
			p.manager.Unlock()
		}
	}
}

func (p *Plugin) Configure(global *plugin.GlobalOptions, options interface{}) {
	if p.ConfigureFunc != nil {
		p.ConfigureFunc(p, &GlobalOptions{*global}, options)
	}
}

func (p *Plugin) Validate(options interface{}) error {
	if p.ValidateFunc != nil {
		return p.ValidateFunc(options)
	}
	return nil
}

func (p *Plugin) Logger(name string) Logger {
	loggerName := p.name
	if name != "" {
		loggerName += "." + name
	}
	return log.New(loggerName)
}
