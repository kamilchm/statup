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

type Checkin struct {
	Id               int64     `gorm:"primary_key;column:id"`
	Service          int64     `gorm:"index;column:service"`
	Interval         int64     `gorm:"column:check_interval"`
	Api              string    `gorm:"column:api"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	Hits             int64     `json:"hits"`
	Last             time.Time `json:"last"`
	CheckinInterface `json:"-"`
}

type CheckinInterface interface {
	// Database functions
	Create() (int64, error)
	//Update() error
	//Delete() error
	Ago() string
	Receivehit()
}
