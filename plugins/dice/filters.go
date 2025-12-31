/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

// This file defines watch filters used by the Dice example plugin.
//
// Watch filters are responsible for converting collected event payloads into
// string values returned to Zabbix items. A filter instance is created per item key
// (and per key parameters) via Metric.NewWatchFilter.
package dice

import (
	"fmt"

	"golang.zabbix.com/agent2/zbxedx/edxagent"
	"golang.zabbix.com/agent2/zbxedx/edxagent/edxerr"
)

// throwsFilter returns a plugin-wide counter value for the given item key.
type throwsFilter struct {
}

func NewThrowsFilter(_ []string) (edxagent.WatchFilter, error) {
	return &throwsFilter{}, nil
}

func (d *throwsFilter) Process(v interface{}) (value *string, err error) {
	s := fmt.Sprintf("%d", Throws)
	return &s, nil
}

// diceFilter extracts a parameter-specific dice value from the event payload.
type diceFilter struct {
	diceID string
}

// NewDiceFilter creates a watch filter instance for a specific dice ID.
// It is called when the agent creates a handler for an item key with parameters
// (e.g., on startup, configuration reload, or when a new item appears in the Zabbix server configuration).
func NewDiceFilter(params []string) (edxagent.WatchFilter, error) {
	if len(params) != 1 {
		return nil, edxerr.New("invalid number of parameters for dice event key")
	}
	return &diceFilter{diceID: params[0]}, nil
}

// Process extracts the dice value from the event payload published
// by the asynchronous metrics collector.
func (d *diceFilter) Process(v interface{}) (*string, error) {
	var s string
	res, ok := v.(map[string]int)
	if !ok {
		return &s, edxerr.New("unexpected type for dice result: " + fmt.Sprintf("%T", v))
	}
	dice, ok := res[d.diceID]
	if !ok {
		return nil, edxerr.New("unknown dice: " + d.diceID)
	} else {
		s = fmt.Sprintf("%d", dice)
	}
	return &s, nil
}
