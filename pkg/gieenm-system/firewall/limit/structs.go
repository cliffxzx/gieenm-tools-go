package limit

import (
	"fmt"
	"time"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Limit ...
type Limit struct {
	ID             *int         `db:"id"`
	UID            *string      `db:"uid"`
	RecordMaxCount *int         `db:"record_max_count"`
	User           *user.User   `db:""`
	Group          *group.Group `db:""`
	CreatedAt      *time.Time   `db:"created_at"`
	ModifiedAt     *time.Time   `db:"modified_at"`
}

// IsNode ...
func (Limit) IsNode() {}

func (r Limit) String() string {
	return fmt.Sprintf(
		"{ ID: %s, UID: %s, MaxCount: %s, CreatedAt: %s, ModifiedAt: %s, \nGroupID: %s, \nUserID: %s }",
		utils.MustToJSON(r.ID),
		utils.MustToJSON(r.UID),
		utils.MustToJSON(r.RecordMaxCount),
		utils.MustToJSON(r.CreatedAt),
		utils.MustToJSON(r.ModifiedAt),
		utils.MustToJSON(r.Group.String()),
		utils.MustToJSON(r.User),
	)
}
