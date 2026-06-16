package bot

import (
	"regexp"
	"testing"

	"gopkg.in/telebot.v3"
)

func encodeForSend(m *telebot.ReplyMarkup) {
	for i := range m.InlineKeyboard {
		for j := range m.InlineKeyboard[i] {
			k := &m.InlineKeyboard[i][j]
			if k.Unique == "" {
				continue
			}
			if k.Data == "" {
				k.Data = "\f" + k.Unique
			} else {
				k.Data = "\f" + k.Unique + "|" + k.Data
			}
		}
	}
}

var cbackRx = regexp.MustCompile(`^\f([-\w]+)(\|(.+))?$`)

func decodePayload(data string) string {
	m := cbackRx.FindStringSubmatch(data)
	if m == nil {
		return ""
	}
	return m[3]
}

func TestSeverityPayloadSurvivesRepeatedSends(t *testing.T) {
	want := map[int]string{0: "LOW", 1: "MEDIUM", 2: "HIGH"}

	for send := 1; send <= 3; send++ {
		m := severityMenu()
		encodeForSend(m)
		for col, level := range want {
			got := decodePayload(m.InlineKeyboard[0][col].Data)
			if got != level {
				t.Fatalf("send %d, button %d: payload = %q, want %q", send, col, got, level)
			}
		}
	}
}

func TestMenusAreFreshPerCall(t *testing.T) {
	cases := map[string]func() *telebot.ReplyMarkup{
		"severityMenu": severityMenu,
		"incidentMenu": incidentMenu,
	}
	for name, build := range cases {
		a, b := build(), build()
		a.InlineKeyboard[0][0].Data = "MUTATED"
		if b.InlineKeyboard[0][0].Data == "MUTATED" {
			t.Fatalf("%s returns a shared markup; telebot would corrupt callback_data across sends", name)
		}
	}
}

func TestSharedMenuCorruptsPayload(t *testing.T) {
	shared := severityMenu()
	encodeForSend(shared)
	encodeForSend(shared)

	if got := decodePayload(shared.InlineKeyboard[0][2].Data); got == "HIGH" {
		t.Fatal("expected a reused markup to corrupt the payload across sends")
	}
}
