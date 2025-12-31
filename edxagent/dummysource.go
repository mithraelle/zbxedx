/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import "golang.zabbix.com/agent2/pkg/watch"

type dummySource struct {
}

func (e *dummySource) Initialize() (err error) {
	return nil
}

func (e *dummySource) Release() {
}

func (e *dummySource) NewFilter(key string) (watch.EventFilter, error) {
	return nil, nil
}
