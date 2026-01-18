/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

// This file defines the metric set and the asynchronous collection loop used by the Dice example plugin.
//
// Metrics in this plugin are collected in a dedicated goroutine and
// published independently from the Zabbix agent's default reporting cycle.
// This demonstrates how plugins can force metric flushing at custom intervals.
package dice

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.zabbix.com/agent2/zbxedx/edxagent"
)

// diceState holds runtime state for the Dice metric set.
// It is shared between the asynchronous collector and watch key updates.
type diceState struct {
	edxagent.Logger
	Faces         int
	Throws        int
	Dices         []string
	ThrowInterval int
	mu            sync.Mutex
}

// dice is the global runtime instance used by the Dice metric set.
// It is initialized with default values and later configured in Configure().
var dice = diceState{
	Faces:         6,                 // Number of faces on a single dice.
	Throws:        0,                 // Total number of dice throws.
	Dices:         make([]string, 0), // List of configured dice identifiers.
	ThrowInterval: 5,                 // Interval in seconds between dice throws.
}

// DiceMetrics defines the Dice metric set.
//
// It contains:
//   - metric declarations
//   - asynchronous collector
//   - watch key update handler
//
// This structure is declarative and registered as part of the plugin.
var DiceMetrics = edxagent.MetricSet{
	Name: "dice",
	// Metric definitions exposed by the Dice plugin.
	Metrics: []edxagent.Metric{
		// Metrics collected by the custom asynchronous collector.
		// Collection frequency is controlled by the plugin itself.
		{
			Key:         "dice",
			Description: "The last dice roll result.",
			// NewWatchFilter constructs a filter instance specific to the item key and parameters.
			// See filters.go for details on how watch filters extract values from collected events.
			NewWatchFilter: NewDiceFilter,
		},
		{
			Key:            "dice.throws",
			Description:    "The number of dice throws.",
			NewWatchFilter: NewThrowsFilter,
		},

		// Metric exported via the built-in Zabbix agent collection cycle.
		// The collection interval is defined by the Zabbix server configuration.
		{
			Key:         "dice.faces",
			Description: "The number of faces on the dice.",
			Export: func(_ []string) (any, error) {
				return dice.Faces, nil
			},
		},
	},

	// Collect prepares metric values and publishes them to the Zabbix collector.
	// It runs in a dedicated goroutine and flushes data independently to the agent's default reporting cycle.
	Collect: rollDice,

	// Update is called when the agent updates active item keys for this metric set.
	Update: updateDices,
}

// rollDice runs the asynchronous collection loop.
// It periodically generates dice values and publishes them using the provided Publisher.
func rollDice(ctx context.Context, source edxagent.Publisher) {
	ticker := time.NewTicker(time.Duration(dice.ThrowInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			res := make(map[string]int)
			dice.mu.Lock() // Lock the dice state while generating new values
			for _, diceID := range dice.Dices {
				res[diceID] = 1 + rand.Intn(dice.Faces)
			}
			dice.Throws++
			dice.mu.Unlock()

			source.Publish(res)
		}
	}
}

// updateDices updates the list of active dice identifiers
// based on the current set of watched item keys.
func updateDices(items []edxagent.WatchKey) {
	dice.mu.Lock() // Lock the dice state for update
	defer dice.mu.Unlock()
	dice.Dices = make([]string, 0) //Recreate the dice list
	for _, item := range items {
		dice.Infof("Updated item key: %s. Params: %s", item.Key, strings.Join(item.Params, " | "))
		if item.Key == "dice" && len(item.Params) > 0 {
			dice.Dices = append(dice.Dices, item.Params[0])
		}
	}
}

// ConfigureDice applies plugin configuration values to the dice runtime state
// and assigns the metric-set scoped logger.
func (d *diceState) ConfigureDice(opts *PluginOptions, logger edxagent.Logger) error {
	d.ThrowInterval = opts.ThrowInterval
	d.Faces = opts.Faces

	d.Logger = logger
	return nil
}
