// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"time"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/base/scalars"
)

type AddRecordInput struct {
	GroupID string           `json:"groupID"`
	Name    string           `json:"name"`
	MacAddr *scalars.MacAddr `json:"macAddr"`
}

type AutoSign struct {
	ID        string     `json:"id"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	Content   *string    `json:"content"`
}

func (AutoSign) IsNode() {}

type DelRecordInput struct {
	ID string `json:"id"`
}

type SetRecordInput struct {
	ID      string           `json:"id"`
	Name    *string          `json:"name"`
	GroupID *string          `json:"groupID"`
	MacAddr *scalars.MacAddr `json:"macAddr"`
}