package bot

import "gopkg.in/telebot.v3"

func HandleStart(c telebot.Context) error {
	return c.Send("Incident War Room is running.")
}
