package bot

import (
	"regexp"
	"testing"

	"gopkg.in/telebot.v3"
)

// encodeForSend mimics telebot's processButtons (options.go): before every
// send it rewrites each button's callback_data in place to "\f<unique>|<data>".
// Because it mutates in place, sending the same *ReplyMarkup twice prefixes it
// twice — which is the bug this package guards against.
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

// cbackRx mirrors telebot's callback parser (bot.go). The payload telebot hands
// to Context.Data is group 3 — everything after the first "|".
var cbackRx = regexp.MustCompile(`^\f([-\w]+)(\|(.+))?$`)

func decodePayload(data string) string {
	m := cbackRx.FindStringSubmatch(data)
	if m == nil {
		return ""
	}
	return m[3]
}

// TestSeverityPayloadSurvivesRepeatedSends is the regression test for the
// corrupted-severity bug: each send must start from a clean markup so the
// callback payload decodes back to the exact severity level.
func TestSeverityPayloadSurvivesRepeatedSends(t *testing.T) {
	want := map[int]string{0: "LOW", 1: "MEDIUM", 2: "HIGH"}

	for send := 1; send <= 3; send++ {
		m := severityMenu() // fresh per send, as the handlers now do
		encodeForSend(m)
		for col, level := range want {
			got := decodePayload(m.InlineKeyboard[0][col].Data)
			if got != level {
				t.Fatalf("send %d, button %d: payload = %q, want %q", send, col, got, level)
			}
		}
	}
}

// TestMenusAreFreshPerCall proves the builders do not hand out a shared markup;
// mutating one returned value must not be visible in the next.
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

// TestSharedMenuCorruptsPayload documents why the builders return fresh markups:
// reusing one markup across two sends double-prefixes the data and breaks it.
func TestSharedMenuCorruptsPayload(t *testing.T) {
	shared := severityMenu()
	encodeForSend(shared)
	encodeForSend(shared) // second send of the SAME markup

	if got := decodePayload(shared.InlineKeyboard[0][2].Data); got == "HIGH" {
		t.Fatal("expected a reused markup to corrupt the payload across sends")
	}
}
