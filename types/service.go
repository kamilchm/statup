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

type Service struct {
	Id             int64         `gorm:"primary_key;column:id" json:"id"`
	Name           string        `gorm:"column:name" json:"name"`
	Domain         string        `gorm:"column:domain" json:"domain"`
	Expected       string        `gorm:"not null;column:expected" json:"expected"`
	ExpectedStatus int           `gorm:"default:200;column:expected_status" json:"expected_status"`
	Interval       int           `gorm:"default:30;column:check_interval" json:"check_interval"`
	Type           string        `gorm:"column:check_type" json:"type"`
	Method         string        `gorm:"column:method" json:"method"`
	PostData       string        `gorm:"not null;column:post_data" json:"post_data"`
	Port           int           `gorm:"not null;column:port" json:"port"`
	Timeout        int           `gorm:"default:30;column:timeout" json:"timeout"`
	Order          int           `gorm:"default:0;column:order_id" json:"order_id"`
	CreatedAt      time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at" json:"updated_at"`
	Online         bool          `gorm:"-" json:"online"`
	Latency        float64       `gorm:"-" json:"latency"`
	Online24Hours  float32       `gorm:"-" json:"24_hours_online"`
	AvgResponse    string        `gorm:"-" json:"avg_response"`
	Running        chan bool     `gorm:"-" json:"-"`
	Checkpoint     time.Time     `gorm:"-" json:"-"`
	SleepDuration  time.Duration `gorm:"-" json:"-"`
	LastResponse   string        `gorm:"-" json:"-"`
	LastStatusCode int           `gorm:"-" json:"status_code"`
	LastOnline     time.Time     `gorm:"-" json:"last_online"`
	DnsLookup      float64       `gorm:"-" json:"dns_lookup_time"`
	Failures       []interface{} `gorm:"-" json:"failures,omitempty"`
	Checkins       []*Checkin    `gorm:"-" json:"checkins,omitempty"`
}

type ServiceInterface interface {
	Select() *Service
	CheckQueue(bool)
	Check(bool)
	Create(bool) (int64, error)
	Update(bool) error
	Delete() error
}

// Start will create a channel for the service checking go routine
func (s *Service) Start() {
	s.Running = make(chan bool)
}

// Close will stop the go routine that is checking if service is online or not
func (s *Service) Close() {
	if s.IsRunning() {
		close(s.Running)
	}
}

// IsRunning returns true if the service go routine is running
func (s *Service) IsRunning() bool {
	if s.Running == nil {
		return false
	}
	select {
	case <-s.Running:
		return false
	default:
		return true
	}
}
