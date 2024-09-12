package message

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func WithReplyToMessageID(id int) sendOption {
	return func(msg tgbotapi.Chattable) tgbotapi.Chattable {
		switch m := msg.(type) {
		case tgbotapi.MessageConfig:
			m.ReplyToMessageID = id
			return m
		}
		return msg
	}
}

func WithForceReply() sendOption {
	return func(msg tgbotapi.Chattable) tgbotapi.Chattable {
		switch m := msg.(type) {
		case tgbotapi.MessageConfig:
			m.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true}
			return m
		}
		return msg
	}
}

type KeyboardButton = [2]string
type KeyboardRow = []KeyboardButton

func WithMessageButtons(buttons ...KeyboardRow) sendOption {
	return func(msg tgbotapi.Chattable) tgbotapi.Chattable {
		var rows [][]tgbotapi.InlineKeyboardButton
		for _, row := range buttons {
			var buttons []tgbotapi.InlineKeyboardButton
			for _, button := range row {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button[0], button[1]))
			}
			rows = append(rows, buttons)
		}

		markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		return WithReplyMarkup(markup)(msg)
	}
}

func WithWebappButton(text, url string) sendOption {
	return func(msg tgbotapi.Chattable) tgbotapi.Chattable {
		markup := tgbotapi.NewReplyKeyboard(
			[]tgbotapi.KeyboardButton{
				tgbotapi.NewKeyboardButtonWebApp(text, tgbotapi.WebAppInfo{URL: url}),
			},
		)
		return WithReplyMarkup(markup)(msg)
	}
}

func WithReplyMarkup(markup any) sendOption {
	return func(msg tgbotapi.Chattable) tgbotapi.Chattable {
		switch m := msg.(type) {
		case tgbotapi.EditMessageTextConfig:
			if markup == nil {
				m.ReplyMarkup = nil
			}
			if c, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok {
				m.ReplyMarkup = &c
			}
			return m
		case tgbotapi.MessageConfig:
			m.ReplyMarkup = markup
			return m
		}
		return msg
	}
}
