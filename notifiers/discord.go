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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterlong/statup/core/notifier"
	"github.com/hunterlong/statup/types"
	"io/ioutil"
	"net/http"
	"time"
)

type Discord struct {
	*notifier.Notification
}

var discorder = &Discord{&notifier.Notification{
	Method:      "discord",
	Title:       "Discord",
	Description: "Send notifications to your discord channel using discord webhooks. Insert your Discord channel webhook URL to receive notifications. Based on the <a href=\"https://discordapp.com/developers/docs/resources/webhook\">Discord Webhook API</a>.",
	Author:      "Hunter Long",
	AuthorUrl:   "https://github.com/hunterlong",
	Delay:       time.Duration(5 * time.Second),
	Host:        "https://discordapp.com/api/webhooks/****/*****",
	Form: []notifier.NotificationForm{{
		Type:        "text",
		Title:       "Discord Webhook URL",
		Placeholder: "Insert your webhook URL here",
		DbField:     "host",
	}}},
}

// init the Discord notifier
func init() {
	err := notifier.AddNotifier(discorder)
	if err != nil {
		panic(err)
	}
}

// Send will send a HTTP Post to the Discord API. It accepts type: []byte
func (u *Discord) Send(msg interface{}) error {
	message := msg.(string)
	req, _ := http.NewRequest("POST", discorder.GetValue("host"), bytes.NewBuffer([]byte(message)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func (u *Discord) Select() *notifier.Notification {
	return u.Notification
}

// OnFailure will trigger failing service
func (u *Discord) OnFailure(s *types.Service, f *types.Failure) {
	msg := fmt.Sprintf(`{"content": "Your service '%v' is currently failing! Reason: %v"}`, s.Name, f.Issue)
	u.AddQueue(msg)
	u.Online = false
}

// OnSuccess will trigger successful service
func (u *Discord) OnSuccess(s *types.Service) {
	if !u.Online {
		msg := fmt.Sprintf(`{"content": "Your service '%v' is back online!"}`, s.Name)
		u.AddQueue(msg)
	}
	u.Online = true
}

// OnSave triggers when this notifier has been saved
func (u *Discord) OnSave() error {
	msg := fmt.Sprintf(`{"content": "The Discord notifier on Statup was just updated."}`)
	u.AddQueue(msg)
	return nil
}

// OnSave triggers when this notifier has been saved
func (u *Discord) OnTest() error {
	outError := errors.New("Incorrect Discord URL, please confirm URL is correct")
	message := `{"content": "Testing the Discord notifier"}`
	req, _ := http.NewRequest("POST", discorder.Host, bytes.NewBuffer([]byte(message)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	if string(contents) == "" {
		return nil
	}
	var d discordTestJson
	err = json.Unmarshal(contents, &d)
	if err != nil {
		return outError
	}
	if d.Code == 0 {
		return outError
	}
	fmt.Println("discord: ", string(contents))
	return nil
}

type discordTestJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
