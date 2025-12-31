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
	"time"

	"golang.zabbix.com/agent2/zbxedx/edxagent"
)

var Faces = 6                 // Number of faces on a single dice.
var Throws = 0                // Total number of dice throws.
var Dices = make([]string, 0) // List of configured dice identifiers.
var ThrowInterval = 5         // Interval in seconds between dice throws.

var DiceMetrics = edxagent.MetricSet{
	Name: "dice",
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
				return Faces, nil
			},
		},
	},

	// Collect prepares metric values and publishes them to the Zabbix collector.
	// It runs in a dedicated goroutine and flushes data independently
	// from the agent's default reporting cycle.
	Collect: rollDice,
}

// rollDice runs the asynchronous collection loop.
// It periodically generates dice values and publishes them using the provided Publisher.
func rollDice(ctx context.Context, source edxagent.Publisher) {
	ticker := time.NewTicker(time.Duration(ThrowInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			res := make(map[string]int)
			for _, diceID := range Dices {
				res[diceID] = 1 + rand.Intn(Faces)
			}

			Throws++
			source.Publish(res)
		}
	}
}
