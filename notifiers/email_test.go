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

package notifiers

import (
	"fmt"
	"github.com/hunterlong/statup/core/notifier"
	"github.com/hunterlong/statup/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	EMAIL_HOST     = os.Getenv("EMAIL_HOST")
	EMAIL_USER     = os.Getenv("EMAIL_USER")
	EMAIL_PASS     = os.Getenv("EMAIL_PASS")
	EMAIL_OUTGOING = os.Getenv("EMAIL_OUTGOING")
	EMAIL_SEND_TO  = os.Getenv("EMAIL_SEND_TO")
	EMAIL_PORT     = utils.StringInt(os.Getenv("EMAIL_PORT"))
)

var testEmail *EmailOutgoing

func init() {
	EMAIL_HOST = os.Getenv("EMAIL_HOST")
	EMAIL_USER = os.Getenv("EMAIL_USER")
	EMAIL_PASS = os.Getenv("EMAIL_PASS")
	EMAIL_OUTGOING = os.Getenv("EMAIL_OUTGOING")
	EMAIL_SEND_TO = os.Getenv("EMAIL_SEND_TO")
	EMAIL_PORT = utils.StringInt(os.Getenv("EMAIL_PORT"))

	emailer.Host = EMAIL_HOST
	emailer.Username = EMAIL_USER
	emailer.Password = EMAIL_PASS
	emailer.Var1 = EMAIL_OUTGOING
	emailer.Var2 = EMAIL_SEND_TO
	emailer.Port = int(EMAIL_PORT)
}

func TestEmailNotifier(t *testing.T) {
	t.Parallel()
	if EMAIL_HOST == "" || EMAIL_USER == "" || EMAIL_PASS == "" {
		t.Log("Email notifier testing skipped, missing EMAIL_ environment variables")
		t.SkipNow()
	}
	currentCount = CountNotifiers()

	t.Run("New Emailer", func(t *testing.T) {
		emailer.Host = EMAIL_HOST
		emailer.Username = EMAIL_USER
		emailer.Password = EMAIL_PASS
		emailer.Var1 = EMAIL_OUTGOING
		emailer.Var2 = EMAIL_SEND_TO
		emailer.Port = int(EMAIL_PORT)
		emailer.Delay = time.Duration(100 * time.Millisecond)

		testEmail = &EmailOutgoing{
			To:       emailer.GetValue("var2"),
			Subject:  fmt.Sprintf("Service %v is Failing", TestService.Name),
			Template: TEMPLATE,
			Data:     TestService,
			From:     emailer.GetValue("var1"),
		}
	})

	t.Run("Add Email Notifier", func(t *testing.T) {
		err := notifier.AddNotifier(emailer)
		assert.Nil(t, err)
		assert.Equal(t, "Hunter Long", emailer.Author)
		assert.Equal(t, EMAIL_HOST, emailer.Host)
	})

	t.Run("Emailer Load", func(t *testing.T) {
		notifier.Load()
	})

	t.Run("Email Within Limits", func(t *testing.T) {
		ok, err := emailer.WithinLimits()
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Email Test Source", func(t *testing.T) {
		emailSource(testEmail)
		assert.NotEmpty(t, testEmail.Source)
	})

	t.Run("Email OnFailure", func(t *testing.T) {
		emailer.OnFailure(TestService, TestFailure)
		assert.Len(t, emailer.Queue, 1)
	})

	t.Run("Email Check Offline", func(t *testing.T) {
		assert.False(t, emailer.Online)
	})

	t.Run("Email OnSuccess", func(t *testing.T) {
		emailer.OnSuccess(TestService)
		assert.Len(t, emailer.Queue, 2)
	})

	t.Run("Email Check Back Online", func(t *testing.T) {
		assert.True(t, emailer.Online)
	})

	t.Run("Email OnSuccess Again", func(t *testing.T) {
		emailer.OnSuccess(TestService)
		assert.Len(t, emailer.Queue, 2)
	})

	t.Run("Email Send", func(t *testing.T) {
		err := emailer.Send(testEmail)
		assert.Nil(t, err)
	})

	t.Run("Email Run Queue", func(t *testing.T) {
		go notifier.Queue(emailer)
		time.Sleep(5 * time.Second)
		assert.Equal(t, EMAIL_HOST, emailer.Host)
		assert.Equal(t, 0, len(emailer.Queue))
	})

}
