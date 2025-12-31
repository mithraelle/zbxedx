/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import "golang.zabbix.com/agent2/pkg/watch"

type WatchKey struct {
	Key    string
	Params []string
}

type WatchFilter interface {
	watch.EventFilter
}
