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

// PluginOptions defines configuration parameters for the Dice plugin.
// Relates to conf/plugin_dice.conf
type PluginOptions struct {
	// Defines the force-push interval (in seconds) for plugin-collected metrics.
	// These metrics are collected asynchronously and flushed independently of the agent's default reporting cycle.
	// Configured via Plugins.DicePlugin.Interval.
	Interval int `conf:"default=10"`

	// Sets the number of dice faces.
	// Configured via Plugins.DicePlugin.Faces.
	Faces int `conf:"default=6"`

	// Defines the interval (in seconds) between dice throws.
	ThrowInterval int `conf:"default=5"`
}

// init registers the Dice plugin and its metric sets in the Zabbix Agent2.
func init() {
	err := edxagent.Register(pluginName,
		[]*edxagent.MetricSet{&DiceMetrics},       //List of plugin metric sets
		edxagent.WithValidateFunction(Validate),   // Optional configuration validation
		edxagent.WithConfigureFunction(Configure), // Optional configuration file parser
	)

	// Registration failures are fatal and prevent agent startup.
	if err != nil {
		panic(edxerr.Wrap(err, "failed to register metrics"))
	}
}

// Validate performs basic sanity checks for plugin configuration options.
func Validate(options interface{}) error {
	var opts PluginOptions

	//Unwrap configuration parameters
	if err := edxagent.Unmarshal(options, &opts); err != nil {
		return edxerr.New("cannot unmarshal configuration options: " + err.Error())
	}

	//Sample parameter validation
	if opts.Faces < 2 {
		return edxerr.New("number of dice faces must be at least 2")
	}

	if opts.ThrowInterval < 1 {
		return edxerr.New("throw interval must be at least 1 second")
	}

	if opts.Interval < 1 {
		return edxerr.New("force-push interval must be at least 1 second")
	}

	return nil
}

// Configure parses plugin-specific configuration from the Zabbix agent config file,
// applies validated configuration values and initializes runtime state for the Dice metric set.
func Configure(p *edxagent.Plugin, global *edxagent.GlobalOptions, options interface{}) {
	config := PluginOptions{}

	p.Infof("DICE plugin Configure called")
	if err := edxagent.Unmarshal(options, &config); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}
	p.Infof(
		"Dice plugin configuration: faces=%d, throw interval=%d, force-push interval=%d",
		config.Faces,
		config.ThrowInterval,
		config.Interval,
	)

	// Pass configuration to the dice metric state
	err := dice.ConfigureDice(&config, p.Logger("dice"))
	if err != nil {
		p.Errf("Can't configure dice plugin: %s", err)
		return
	}
	p.SetFlushInterval(config.Interval)
}
