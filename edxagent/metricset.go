/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxagent

import (
	"context"
)

type CollectFn func(ctx context.Context, source Publisher)
type UpdateFn func(items []WatchKey)

type MetricSet struct {
	Name    string
	Metrics []Metric
	Collect CollectFn
	Update  UpdateFn
}
