package firewall

import (
	"bytes"
	"errors"
	"net"
	"sort"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/base/scalars"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/common"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/group"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/limit"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/firewall/record"
	"github.com/cliffxzx/gieenm-tools/pkg/nusoft-firewall"
	"github.com/openlyinc/pointy"
	"github.com/thoas/go-funk"
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

func GetNextIP(g group.Group) (*scalars.IPAddr, error) {
	records, err := record.GetsByGroupID(*g.ID)
	if err != nil {
		return nil, err
	}

	sort.Slice(*records, func(i, j int) bool { return bytes.Compare((*records)[i].IPAddr.IP, (*records)[j].IPAddr.IP) < 0 })

	result := (*records)[0].IPAddr
	result.IP = result.IP.To4()
	result.IP[3] = 1

	for _, r := range *records {
		if result.IP[3] < 1 && result.IP[3] >= 200 {
			return nil, ErrorRecordCountExceed
		}

		if net.IP.Equal(r.IPAddr.IP, result.IP) {
			result.IP[3]++
		} else {
			return result, nil
		}
	}

	if result.IP[3] < 1 && result.IP[3] >= 200 {
		return nil, ErrorRecordCountExceed
	}

	return result, nil
}

// AddRecordWithNusoft record to database and return to parameter
//
// Require Name, MacAddr, User, Group.UID fields.
func AddRecordWithNusoft(r *record.Record) error {
	if r.Name == nil || r.MacAddr == nil || r.User.ID == nil || r.Group.UID == nil {
		return errors.New("requried Name, MacAddr, User.ID, Group.UID fields")
	}

	g, err := group.GetByUID(*r.Group.UID)
	if err != nil {
		return err
	}

	ip, err := GetNextIP(*g)
	if err != nil {
		return err
	}

	fw := (*GetFirewalls())[*g.FirewallID]

	nuRecord := nusoft.Record{
		Name:    r.Name,
		MacAddr: (*net.HardwareAddr)(r.MacAddr),
		IPAddr:  (*net.IPNet)(ip),
	}

	nuRecord, err = fw.Nusoft.AddRecord(nuRecord)
	if err != nil {
		return err
	}

	r.NusoftID = (*common.NusoftID)(nuRecord.ID)
	r.IPAddr = ip
	r.Group = g

	destSource := &[]record.Record{*r}
	err = record.Adds(destSource)
	if err != nil {
		return err
	}

	r.ID = (*destSource)[0].ID
	r.UID = (*destSource)[0].UID
	r.CreatedAt = (*destSource)[0].CreatedAt
	r.ModifiedAt = (*destSource)[0].ModifiedAt

	infos, err := fw.Nusoft.GetRecordGroupInfos()
	if err != nil {
		return err
	}

	info := funk.Find(infos, func(i nusoft.RecordGroupInfo) bool {
		return *i.ID.Time == *g.NusoftID.Time && *i.ID.Serial == *g.NusoftID.Serial
	}).(nusoft.RecordGroupInfo)
	nuGroup, err := fw.Nusoft.GetRecordGroup(info)
	if err != nil {
		return err
	}

	includes := *nuGroup.Includes
	includes = append(includes, nuRecord)
	err = fw.Nusoft.SetRecordGroup(info, includes)
	if err != nil {
		return err
	}

	return nil
}

// DelRecordWithNusoft record to database and return to parameter
//
// Require UID User fields.
func DelRecordWithNusoft(r *record.Record) error {
	if r.UID == nil {
		return errors.New("requried UID fields")
	}

	rtmp, err := record.GetsByUID(*r.UID)
	if err != nil {
		return err
	}

	*r = *rtmp

	g, err := group.GetByID(*r.Group.ID)
	if err != nil {
		return err
	}

	fw := (*GetFirewalls())[*g.FirewallID]

	err = fw.Nusoft.DelRecord(nusoft.NusoftID(*r.NusoftID))
	if err != nil {
		return err
	}

	destSource := &[]record.Record{*r}
	err = record.Dels(destSource)
	if err != nil {
		return err
	}

	r.Group = g

	return nil
}
