package firewall

import (
	"database/sql/driver"
	"fmt"
	"net/url"
	"time"

	"github.com/cliffxzx/gieenm-tools/pkg/nusoft-firewall"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
)

// URL ...
type URL url.URL

// Scan ...
func (u *URL) Scan(value interface{}) error {
	url, err := url.Parse(value.(string))
	if err != nil {
		return err
	}

	*u = URL(*url)

	return nil
}

// Value ...
func (u URL) Value() (driver.Value, error) {
	url := url.URL(u)
	return (&url).String(), nil
}

// ToStdURL ...
func (u *URL) ToStdURL() *url.URL {
	if u == nil {
		return nil
	}

	result := url.URL(*u)

	return &result
}

// Firewall ...
type Firewall struct {
	ID         *int       `db:"id"`
	UID        *string    `db:"uid"`
	Name       *string    `db:"name"`
	Host       *URL       `db:"host"`
	Username   *string    `db:"username"`
	Password   *string    `db:"password"`
	CreatedAt  *time.Time `db:"created_at"`
	ModifiedAt *time.Time `db:"modified_at"`
	Nusoft     *nusoft.Firewall
}

// String ...
func (f Firewall) String() string {
	return fmt.Sprintf(
		`{ Nusoft: %s, ID: %s, UID: %s, Name: %s, Host: %s, Username: %s, Password: %s, CreatedAt: %s, ModifiedAt: %s }`,
		f.Nusoft.String(),
		utils.MustToJSON(f.ID),
		utils.MustToJSON(f.UID),
		utils.MustToJSON(f.Name),
		utils.URLToJSON(f.Host.ToStdURL()),
		utils.MustToJSON(f.Username),
		utils.MustToJSON(f.Password),
		utils.MustToJSON(f.CreatedAt),
		utils.MustToJSON(f.ModifiedAt),
	)
}

// Firewalls ...
type Firewalls map[int]Firewall

// String ...
func (f Firewalls) String() string {
	result := "[\n"
	for _, firewall := range f {
		result += fmt.Sprintf("\t%s", firewall.String())
		result += ",\n"
	}
	result += "]"
	return result
}
