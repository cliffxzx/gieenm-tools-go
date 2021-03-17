package seeder

import (
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/announcement"
	"github.com/openlyinc/pointy"
)

func AnnouncementSeeder() {
	noticeLevel := announcement.AnnounceLevel(announcement.NOTICE)
	infoLevel := announcement.AnnounceLevel(announcement.INFO)
	announcement.Adds(&[]announcement.Announcement{
		{
			Title:     pointy.String("First Announce Title"),
			Content:   pointy.String("# this is the **content** [I am Google](https://google.com)"),
			Announcer: pointy.String("網管工讀生"),
			Level:     &noticeLevel,
		},
		{
			Title:     pointy.String("Second Announce Title"),
			Content:   pointy.String("# this is the **content** [I am Google](https://google.com)"),
			Announcer: pointy.String("網管工讀生"),
			Level:     &infoLevel,
		},
	})
}
