package firewall

import (
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/limit"
	"github.com/openlyinc/pointy"
)

// GetGroupsByUserID ...
//
// required: userID
func GetGroupsByUserID(userID int) (*[]group.Group, error) {
	groups, err := group.Gets()
	if err != nil {
		return nil, err
	}

	oriLimits, err := limit.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	limits := map[int]*limit.Limit{}
	for idx := range *oriLimits {
		limits[*(*oriLimits)[idx].Group.ID] = &(*oriLimits)[idx]
	}

	for idx := range *groups {
		if limit, ok := limits[*(*groups)[idx].ID]; ok {
			(*groups)[idx].MaxCount = limit.RecordMaxCount
		} else {
			(*groups)[idx].MaxCount = pointy.Int(0)
		}
	}

	return groups, nil
}
