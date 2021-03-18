package firewall

import "errors"

var (
	ErrorRecordCountExceed = errors.New("the group's record count is full")
)
