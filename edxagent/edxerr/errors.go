/*
 * Copyright (c) 2025 Aleksei Dorozhkin
 * SPDX-License-Identifier: MIT
 */

package edxerr

import "golang.zabbix.com/sdk/errs"

func New(msg string) error {
	return errs.New(msg)
}

func Wrap(err error, msg string) error {
	return errs.Wrap(err, msg)
}
