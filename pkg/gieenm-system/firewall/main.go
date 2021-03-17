package firewall

import (
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/nusoft-firewall"
)

var firewalls Firewalls

// GetFirewalls ...
func GetFirewalls() *Firewalls {
	return &firewalls
}

// InitFirewalls ...
func InitFirewalls() error {
	fws := []Firewall{}
	sql := Sqls{}.GetFirewalls()
	err := database.GetDB().Select(&fws, sql)
	if err != nil {
		return err
	}

	firewalls = *New(fws)

	return nil
}

// New ...
func New(fws []Firewall) *Firewalls {
	firewalls := Firewalls{}

	for _, fw := range fws {
		fw.Nusoft = nusoft.New(*fw.Name, *fw.Host.ToStdURL(), *fw.Username, *fw.Password)

		firewalls[*fw.ID] = fw
	}

	return &firewalls
}
