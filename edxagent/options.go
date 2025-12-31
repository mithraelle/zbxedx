/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import (
	"golang.zabbix.com/sdk/conf"
	"golang.zabbix.com/sdk/plugin"
)

type GlobalOptions struct {
	plugin.GlobalOptions
}

type Option interface {
	set(p *Plugin)
}

func (o ConfigureFn) set(p *Plugin) {
	p.ConfigureFunc = o
}

func WithConfigureFunction(fn ConfigureFn) Option {
	return fn
}

func (o ValidateFn) set(p *Plugin) {
	p.ValidateFunc = o
}

func WithValidateFunction(fn ValidateFn) Option {
	return fn
}

func Unmarshal(data any, v any) error {
	return conf.Unmarshal(data, v)
}
