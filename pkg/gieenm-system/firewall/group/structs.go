package group

import (
	"fmt"
	"time"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Group ...
type Group struct {
	ID         *int             `db:"id"`
	UID        *string          `db:"uid"          gqlgen:"id"`
	Name       *string          `db:"name"         gqlgen:"name"`
	NusoftID   *common.NusoftID `db:"nusoft_id"`
	FirewallID *int             `db:"firewall_id"`
	MaxCount   *int             `db:""             gqlgen:"maxCount"`
	Subnet     *scalars.IPAddr  `db:"subnet"`
	CreatedAt  *time.Time       `db:"created_at"`
	ModifiedAt *time.Time       `db:"modified_at"`
}

// IsNode ...
func (Group) IsNode() {}

func (r Group) String() string {
	return fmt.Sprintf(
		"{ ID: %s, UID: %s, Name: %s, Subnet: %s, NusoftID: %s, FirewallID: %s, CreatedAt: %s, ModifiedAt: %s }",
		utils.MustToJSON(r.ID),
		utils.MustToJSON(r.UID),
		utils.MustToJSON(r.Name),
		utils.MustToJSON(r.Subnet),
		utils.MustToJSON(r.NusoftID),
		utils.MustToJSON(r.FirewallID),
		utils.MustToJSON(r.CreatedAt),
		utils.MustToJSON(r.ModifiedAt),
	)
}
