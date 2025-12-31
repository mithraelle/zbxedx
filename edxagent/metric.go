/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

type ExportFn func(params []string) (any, error)
type WatchFilterFn func(params []string) (WatchFilter, error)

type Metric struct {
	Key            string
	Description    string
	Export         ExportFn
	NewWatchFilter WatchFilterFn
}
