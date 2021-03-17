package viewer

import (
	"context"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/announcement"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/authentication"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

func Controller(ctx context.Context) (*Viewer, error) {
	gCtx, err := utils.GetGinContext(ctx)
	if err != nil {
		return nil, err
	}

	token, err := authentication.VerifyToken(gCtx)
	if err != nil {
		return nil, err
	}

	v := Viewer{}

	u, err := user.Get(*token.Claims.UserID)
	if err != nil {
		return nil, err
	}

	v.User = &u

	_ = firewall.GetFirewalls()
	recordGroupTmp, err := firewall.GetGroupsByUserID(*u.ID)
	if err != nil {
		return nil, err
	}

	recordGroupMap := map[int]*group.Group{}
	for idx, g := range *recordGroupTmp {
		recordGroupMap[*g.ID] = &(*recordGroupTmp)[idx]
	}

	v.Groups = *recordGroupTmp

	recordsTmp, err := record.GetsByUserID(*u.ID)
	if err != nil {
		return nil, err
	}

	for idx, r := range *recordsTmp {
		(*recordsTmp)[idx].Group = recordGroupMap[*r.Group.ID]
	}

	v.Records = *recordsTmp

	announcementsTmp, err := announcement.Gets()
	if err != nil {
		return nil, err
	}

	v.Announcements = *announcementsTmp

	switch *v.User.Role {
	default:
	case user.STUDENT:
	case user.MANAGER:
	case user.TEACHER:
	case user.GUEST:
	}

	return &v, nil
}
