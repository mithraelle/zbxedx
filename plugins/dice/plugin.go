/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

// Package dice is a minimal example Zabbix Agent2 plugin built on top of
// the zbxedx framework.
//
// It demonstrates:
//   - plugin registration
//   - metric set definition
//   - error handling via edxerr
//   - periodic metric collection
//
// This plugin is intentionally simple and serves as executable documentation.
package dice

import (
	"golang.zabbix.com/agent2/zbxedx/edxagent"
	"golang.zabbix.com/agent2/zbxedx/edxagent/edxerr"
)

// pluginName is the identifier used to register the plugin in Zabbix Agent2.
// It must be unique within the agent process.
const pluginName = "DicePlugin"

// Sample plugin config structure relates to conf/plugin_dice.conf
type PluginOptions struct {
	// Defines the force-push interval (in seconds) for plugin-collected metrics.
	// These metrics are collected asynchronously and flushed independently from the agent's default reporting cycle.
	// Configured via Plugins.DicePlugin.Interval.
	Interval int `conf:"default=5"`

	// Sets the number of dice faces.
	// Configured via Plugins.DicePlugin.Faces.
	Faces int `conf:"default=6"`

	// Defines the list of available dice identifiers.
	// Configured via Plugins.DicePlugin.Dice.
	Dice []string
}

// Sample plugin registration
func init() {
	err := edxagent.Register(pluginName,
		[]edxagent.MetricSet{DiceMetrics},         //List of plugin metric sets
		edxagent.WithValidateFunction(Validate),   // Optional configuration validation
		edxagent.WithConfigureFunction(Configure), // Optional configuration file parser
	)

	// Registration failures are fatal and prevent agent startup.
	if err != nil {
		panic(edxerr.Wrap(err, "failed to register metrics"))
	}
}

// Validate validates plugin configuration options.
func Validate(options interface{}) error {
	var opts PluginOptions

	//Unwrap configuration parameters
	if err := edxagent.Unmarshal(options, &opts); err != nil {
		return edxerr.New("cannot unmarshal configuration options: " + err.Error())
	}

	//Simple check and error reporting
	if len(opts.Dice) == 0 {
		return edxerr.New("no dice specified in configuration")
	}

	return nil
}

// Configure parses plugin-specific configuration from the Zabbix agent config file.
func Configure(p *edxagent.Plugin, global *edxagent.GlobalOptions, options interface{}) {
	config := PluginOptions{}

	p.Infof("DICE plugin Configure called")
	if err := edxagent.Unmarshal(options, &config); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}
	p.Infof(
		"Dice plugin configuration: faces=%d, interval=%d, dices=%s",
		config.Faces,
		config.Interval,
		config.Dice,
	)

	Dices = config.Dice
	Faces = config.Faces
	p.SetFlushInterval(config.Interval)
}
