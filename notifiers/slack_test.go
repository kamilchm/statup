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
	"github.com/hunterlong/statup/core/notifier"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	SLACK_URL    string
	slackMessage = `{"text":"this is a test from the Slack notifier!"}`
)

func init() {
	SLACK_URL = os.Getenv("SLACK_URL")
	slacker.Host = SLACK_URL
}

func TestSlackNotifier(t *testing.T) {
	t.Parallel()
	SLACK_URL = os.Getenv("SLACK_URL")
	slacker.Host = SLACK_URL
	if SLACK_URL == "" {
		t.Log("Slack notifier testing skipped, missing SLACK_URL environment variable")
		t.SkipNow()
	}
	currentCount = CountNotifiers()

	t.Run("Load Slack", func(t *testing.T) {
		slacker.Host = SLACK_URL
		slacker.Delay = time.Duration(100 * time.Millisecond)
		err := notifier.AddNotifier(slacker)
		assert.Nil(t, err)
		assert.Equal(t, "Hunter Long", slacker.Author)
		assert.Equal(t, SLACK_URL, slacker.Host)
	})

	t.Run("Load Slack Notifier", func(t *testing.T) {
		notifier.Load()
	})

	t.Run("Slack Notifier Tester", func(t *testing.T) {
		assert.True(t, slacker.CanTest())
	})

	t.Run("Slack parse message", func(t *testing.T) {
		err := parseSlackMessage(SLACK_TEXT, "this is a test!")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(slacker.Queue))
	})

	t.Run("Slack Within Limits", func(t *testing.T) {
		ok, err := slacker.WithinLimits()
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Slack OnFailure", func(t *testing.T) {
		slacker.OnFailure(TestService, TestFailure)
		assert.Len(t, slacker.Queue, 2)
	})

	t.Run("Slack Check Offline", func(t *testing.T) {
		assert.False(t, slacker.Online)
	})

	t.Run("Slack OnSuccess", func(t *testing.T) {
		slacker.OnSuccess(TestService)
		assert.Len(t, slacker.Queue, 3)
	})

	t.Run("Slack Check Back Online", func(t *testing.T) {
		assert.True(t, slacker.Online)
	})

	t.Run("Slack OnSuccess Again", func(t *testing.T) {
		slacker.OnSuccess(TestService)
		assert.Len(t, slacker.Queue, 3)
	})

	t.Run("Slack Send", func(t *testing.T) {
		err := slacker.Send(slackMessage)
		assert.Nil(t, err)
		assert.Len(t, slacker.Queue, 3)
	})

	t.Run("Slack Queue", func(t *testing.T) {
		go notifier.Queue(slacker)
		time.Sleep(4 * time.Second)
		assert.Equal(t, SLACK_URL, slacker.Host)
		assert.Equal(t, 0, len(slacker.Queue))
	})

}
