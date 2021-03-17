package viewer

import (
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/announcement"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
)

// Viewer ...
type Viewer struct {
	User          *user.User                  `gqlgen:"user"`
	Records       []record.Record             `gqlgen:"records"`
	Groups        []group.Group               `gqlgen:"groups"`
	Announcements []announcement.Announcement `gqlgen:"announcements"`
}
