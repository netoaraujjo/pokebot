package main

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TEXT int64 = iota
	PHOTO
)

type MatcherFunc func(update tgbotapi.Update) bool
type HandlerFunc func(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *Usuario) int64

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
	Users       map[int64]*Usuario
}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{
		States: make(map[int64][]EventHandler),
		Users:  make(map[int64]*Usuario),
	}
}

func (ch *ConversationHandler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Executar a ação correspondente
	// 1. Entrypoints
	// 2. Ações por estado
	// 3. Fallbacks
	userID := update.FromChat().ID
	user, exists := ch.Users[userID]
	log.Println(user)
	if !exists {
		log.Printf("nao existe")
		for _, h := range ch.EntryPoints {
			if h.Match(update) {
				user = &Usuario{
					State:      END,
					ID:         userID,
					Name:       update.FromChat().FirstName,
					Username:   update.FromChat().UserName,
					Parameters: make(map[string]string),
				}
				ch.Users[userID] = user
				user.State = h.Handler(bot, update, user)
				log.Printf("Estado atual: %d\n", user.State)
				return
			}
		}
		return
	}

	log.Println("Existe")

	for _, h := range ch.States[user.State] {
		log.Printf("Estado buscado: %d", user.State)
		log.Println(user)
		if h.Match(update) {
			user.State = h.Handler(bot, update, user)
			if user.State == END {
				delete(ch.Users, userID)
			}
			return
		}
	}

	for _, h := range ch.Fallbacks {
		if h.Match(update) {
			// user.State = END
			h.Handler(bot, update, user)
			delete(ch.Users, userID)
			return
		}
	}

}
