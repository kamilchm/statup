// Statup
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statup
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"time"
)

type AllNotifiers interface{}

// Core struct contains all the required fields for Statup. All application settings
// will be saved into 1 row in the 'core' table. You can use the core.CoreApp
// global variable to interact with the attributes to the application, such as services.
type Core struct {
	Name          string             `gorm:"not null;column:name" json:"name"`
	Description   string             `gorm:"not null;column:description" json:"description,omitempty"`
	Config        string             `gorm:"column:config" json:"-"`
	ApiKey        string             `gorm:"column:api_key" json:"-"`
	ApiSecret     string             `gorm:"column:api_secret" json:"-"`
	Style         string             `gorm:"not null;column:style" json:"style,omitempty"`
	Footer        string             `gorm:"not null;column:footer" json:"footer,omitempty"`
	Domain        string             `gorm:"not null;column:domain" json:"domain,omitempty"`
	Version       string             `gorm:"column:version" json:"version"`
	MigrationId   int64              `gorm:"column:migration_id" json:"migration_id,omitempty"`
	UseCdn        bool               `gorm:"column:use_cdn;default:false" json:"using_cdn,omitempty"`
	Timezone      float32            `gorm:"column:timezone;default:-8.0" json:"timezone,omitempty"`
	CreatedAt     time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time          `gorm:"column:updated_at" json:"updated_at"`
	DbConnection  string             `gorm:"-" json:"database"`
	Started       time.Time          `gorm:"-" json:"started_on"`
	Services      []ServiceInterface `gorm:"-" json:"services,omitempty"`
	Plugins       []Info             `gorm:"-" json:"-"`
	Repos         []PluginJSON       `gorm:"-" json:"-"`
	AllPlugins    []PluginActions    `gorm:"-" json:"-"`
	Notifications []AllNotifiers     `gorm:"-" json:"-"`
	CoreInterface `gorm:"-" json:"-"`
}

type CoreInterface interface {
	SelectAllServices() ([]*Service, error)
	Count24HFailures() uint64
	ServicesCount() int
	CountOnline() int
}
