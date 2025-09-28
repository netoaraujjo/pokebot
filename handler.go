package main

import (
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TEXT int64 = iota
	PHOTO
)

type MatcherFunc func(update tgbotapi.Update) bool
type HandlerFunc func(bot *tgbotapi.BotAPI, update tgbotapi.Update) int64

// Define o formato das funções que tratarão os eventos: mensagem, callback query, comandos
// Retorna o próximo estado
type EventHandler struct {
	Match   MatcherFunc
	Handler HandlerFunc
}

func MessageHandler(msgType int64) MatcherFunc {
	return func(update tgbotapi.Update) bool {
		if update.Message != nil && !update.Message.IsCommand() {
			if (update.Message.Text != "" && msgType == TEXT) || (len(update.Message.Photo) > 0 && msgType == PHOTO) {
				return true
			}
		}
		return false
	}
}

func CommandHandler(command string) MatcherFunc {
	return func(update tgbotapi.Update) bool {
		return update.Message != nil && update.Message.IsCommand() && command == update.Message.Command()
	}
}

func CallbackQueryHandler(data string) MatcherFunc {
	return func(update tgbotapi.Update) bool {
		return update.CallbackQuery != nil && update.CallbackData() == data
	}
}

func PatternHandler(pattern string) MatcherFunc {
	re := regexp.MustCompile(pattern)
	return func(update tgbotapi.Update) bool {
		return update.Message != nil && re.MatchString(update.Message.Text)
	}
}

type ConversationHandler struct {
	EntryPoints []EventHandler
	States      map[int64][]EventHandler
	Fallbacks   []EventHandler
	State       int64
}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{
		States: make(map[int64][]EventHandler),
		State:  -1,
	}
}

func (ch *ConversationHandler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Executar a ação correspondente
	// 1. Entrypoints
	// 2. Ações por estado
	// 3. Fallbacks

	if ch.State == -1 {
		for _, h := range ch.EntryPoints {
			if h.Match(update) {
				ch.State = h.Handler(bot, update)
				return
			}
		}
	}

	if ch.State > -1 {
		for _, h := range ch.States[ch.State] {
			if h.Match(update) {
				ch.State = h.Handler(bot, update)
				return
			}
		}
	}

	for _, h := range ch.Fallbacks {
		if h.Match(update) {
			ch.State = -1
			h.Handler(bot, update)
			return
		}
	}

}
