package sync

import (
	"errors"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/limit"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/graphql/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/user"
	"github.com/cliffxzx/gieenm-tools/pkg/nusoft-firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/openlyinc/pointy"
)

func admin() (*user.User, error) {
	user := &user.User{
		Role:      &user.MANAGER,
		Name:      pointy.String(utils.MustGetEnv("ADMIN_NAME")),
		StudentID: pointy.String(utils.MustGetEnv("ADMIN_STUDENT_ID")),
		Email:     pointy.String(utils.MustGetEnv("ADMIN_EMAIL")),
		Password:  pointy.String(utils.MustGetEnv("ADMIN_PASSWORD")),
	}

	err := user.AddOrGet()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func addUnlimit(gs *[]group.Group, u *user.User) error {
	for _, g := range *gs {
		l := &[]limit.Limit{{User: u, Group: &g, RecordMaxCount: pointy.Int(-1)}}

		err := limit.Adds(l)
		if err != nil {
			return err
		}
	}

	return nil
}

// SyncNusoftToDatabase
func SyncNusoftToDatabase(defaultSubnet map[string]*scalars.IPAddr) error {
	groups := []group.Group{}
	includes := [][]nusoft.Record{}
	records := map[int64]*nusoft.Record{}

	for _, fw := range *firewall.GetFirewalls() {
		infos, err := fw.Nusoft.GetRecordGroupInfos()
		if err != nil {
			return err
		}

		for _, info := range infos {
			nsGroup, err := fw.Nusoft.GetRecordGroup(info)
			if err != nil {
				return err
			}

			includes = append(includes, *nsGroup.Includes)
			nusoftID := (*common.NusoftID)(info.ID)

			groups = append(groups, group.Group{
				Name:       info.Name,
				NusoftID:   nusoftID,
				FirewallID: fw.ID,
				Subnet:     defaultSubnet[fmt.Sprintf("%d,%d", *nusoftID.Time, *nusoftID.Serial)],
			})
		}

		nsRecords, err := fw.Nusoft.GetRecords()
		if err != nil {
			return err
		}

		for idx, r := range nsRecords {
			if *r.ID.Time != 19700102030000 && *r.ID.Serial != 1 && *r.Name != "Inside Any" {
				records[*r.ID.Time] = &nsRecords[idx]
			}
		}
	}

	if err := group.Adds(&groups); err != nil {
		return err
	}

	user, err := admin()
	if err != nil {
		return err
	}

	if err = addUnlimit(&groups, user); err != nil {
		return err
	}

	if len(groups) != len(includes) {
		return errors.New("unknown error: groups and info count not same")
	}

	result := []record.Record{}
	for idx := range groups {
		for _, include := range includes[idx] {
			r := records[*include.ID.Time]
			result = append(result, record.Record{
				Name:     r.Name,
				NusoftID: (*common.NusoftID)(r.ID),
				IPAddr:   (*scalars.IPAddr)(r.IPAddr),
				MacAddr:  (*scalars.MacAddr)(r.MacAddr),
				User:     user,
				Group:    &groups[idx],
			})
		}
	}

	if err = record.Adds(&result); err != nil {
		return err
	}

	return nil
}
