package firewall

import (
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
)

func AddRecordsController(r *record.Record) (*record.Record, error) {
	err := AddRecordWithNusoft(r)
	if err != nil {
		return &record.Record{}, err
	}

	return r, nil
}

func DelRecordsController(r *record.Record) (*record.Record, error) {
	err := DelRecordWithNusoft(r)
	if err != nil {
		return &record.Record{}, err
	}

	return r, nil
}
