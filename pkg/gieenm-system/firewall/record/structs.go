package record

import (
	"fmt"
	"time"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// Record ...
type Record struct {
	ID         *int             `db:"id"`
	UID        *string          `db:"uid"          gqlgen:"id"`
	Name       *string          `db:"name"         gqlgen:"name"`
	NusoftID   *common.NusoftID `db:"nusoft_id"`
	IPAddr     *scalars.IPAddr  `db:"ip_addr"`
	MacAddr    *scalars.MacAddr `db:"mac_addr"     gqlgen:"macAddr"`
	CreatedAt  *time.Time       `db:"created_at"`
	ModifiedAt *time.Time       `db:"modified_at"`
	User       *user.User       `db:""`
	Group      *group.Group     `db:""             gqlgen:"group"`
}

// IsNode ...
func (Record) IsNode() {}

func (r Record) String() string {
	return fmt.Sprintf(
		"{ ID: %s, UID: %s, Name: %s, IPAddr: %s, MacAddr: %s, NusoftID: %s, CreatedAt: %s, ModifiedAt: %s, \n\t\tGroup: %s, \n\t\tUser: %s }",
		utils.MustToJSON(r.ID),
		utils.MustToJSON(r.UID),
		utils.MustToJSON(r.Name),
		utils.IPAddrToJSON(r.IPAddr.ToStd()),
		utils.MacAddrToJSON(r.MacAddr.ToStd()),
		utils.MustToJSON(r.NusoftID.String()),
		utils.MustToJSON(r.CreatedAt),
		utils.MustToJSON(r.ModifiedAt),
		utils.MustToJSON(r.User),
		utils.MustToJSON(r.Group),
	)
}
