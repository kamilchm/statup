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

package core

import (
	"fmt"
	"github.com/ararog/timeago"
	"github.com/hunterlong/statup/types"
	"github.com/hunterlong/statup/utils"
	"strings"
	"time"
)

type Failure struct {
	*types.Failure
}

// CreateFailure will create a new failure record for a service
func (s *Service) CreateFailure(f *types.Failure) (int64, error) {
	f.Service = s.Id
	s.Failures = append(s.Failures, f)
	row := failuresDB().Create(f)
	if row.Error != nil {
		utils.Log(3, row.Error)
		return 0, row.Error
	}
	return f.Id, row.Error
}

// AllFailures will return all failures attached to a service
func (s *Service) AllFailures() []*Failure {
	var fails []*Failure
	col := failuresDB().Where("service = ?", s.Id).Order("id desc")
	err := col.Find(&fails)
	if err.Error != nil {
		utils.Log(3, fmt.Sprintf("Issue getting failures for service %v, %v", s.Name, err))
		return nil
	}
	for _, f := range fails {
		s.Failures = append(s.Failures, f)
	}
	return fails
}

// DeleteFailures will delete all failures for a service
func (u *Service) DeleteFailures() {
	err := DbSession.Exec(`DELETE FROM failures WHERE service = ?`, u.Id)
	if err.Error != nil {
		utils.Log(3, fmt.Sprintf("failed to delete all failures: %v", err))
	}
	u.Failures = nil
}

// LimitedFailures will return the last 10 failures from a service
func (s *Service) LimitedFailures() []*Failure {
	var failArr []*Failure
	col := failuresDB().Where("service = ?", s.Id).Order("id desc").Limit(10)
	col.Find(&failArr)
	return failArr
}

// Ago returns a human readable timestamp for a failure
func (f *Failure) Ago() string {
	got, _ := timeago.TimeAgoWithTime(time.Now(), f.CreatedAt)
	return got
}

// Delete will remove a failure record from the database
func (f *Failure) Delete() error {
	db := failuresDB().Delete(f)
	return db.Error
}

// Count24HFailures returns the amount of failures for a service within the last 24 hours
func (c *Core) Count24HFailures() uint64 {
	var count uint64
	for _, s := range CoreApp.Services {
		service := s.(*Service)
		fails, _ := service.TotalFailures24()
		count += fails
	}
	return count
}

// CountFailures returns the total count of failures for all services
func CountFailures() uint64 {
	var count uint64
	err := failuresDB().Count(&count)
	if err.Error != nil {
		utils.Log(2, err.Error)
		return 0
	}
	return count
}

// TotalFailures24 returns the amount of failures for a service within the last 24 hours
func (s *Service) TotalFailures24() (uint64, error) {
	ago := time.Now().Add(-24 * time.Hour)
	return s.TotalFailuresSince(ago)
}

// TotalFailures returns the total amount of failures for a service
func (s *Service) TotalFailures() (uint64, error) {
	var count uint64
	rows := failuresDB().Where("service = ?", s.Id)
	err := rows.Count(&count)
	return count, err.Error
}

// TotalFailuresSince returns the total amount of failures for a service since a specific time/date
func (s *Service) TotalFailuresSince(ago time.Time) (uint64, error) {
	var count uint64
	rows := failuresDB().Where("service = ? AND created_at > ?", s.Id, ago.UTC().Format("2006-01-02 15:04:05"))
	err := rows.Count(&count)
	return count, err.Error
}

// ParseError returns a human readable error for a failure
func (f *Failure) ParseError() string {
	err := strings.Contains(f.Issue, "connection reset by peer")
	if err {
		return fmt.Sprintf("Connection Reset")
	}
	err = strings.Contains(f.Issue, "operation timed out")
	if err {
		return fmt.Sprintf("HTTP Request Timed Out")
	}
	err = strings.Contains(f.Issue, "x509: certificate is valid")
	if err {
		return fmt.Sprintf("SSL Certificate invalid")
	}
	err = strings.Contains(f.Issue, "Client.Timeout exceeded while awaiting headers")
	if err {
		return fmt.Sprintf("Connection Timed Out")
	}
	err = strings.Contains(f.Issue, "no such host")
	if err {
		return fmt.Sprintf("Domain is offline or not found")
	}
	err = strings.Contains(f.Issue, "HTTP Status Code")
	if err {
		return fmt.Sprintf("Incorrect HTTP Status Code")
	}
	err = strings.Contains(f.Issue, "connection refused")
	if err {
		return fmt.Sprintf("Connection Failed")
	}
	err = strings.Contains(f.Issue, "can't assign requested address")
	if err {
		return fmt.Sprintf("Unable to Request Address")
	}
	err = strings.Contains(f.Issue, "no route to host")
	if err {
		return fmt.Sprintf("Domain is offline or not found")
	}
	return f.Issue
}
