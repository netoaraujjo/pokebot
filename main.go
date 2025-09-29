package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	END   int64 = 0
	START int64 = iota
	LE_NOME
)

func NewConversation() *ConversationHandler {
	ch := NewConversationHandler()
	ch.EntryPoints = []EventHandler{
		{Match: MessageHandler(TEXT), Handler: handleCommandStart},
		{Match: CommandHandler("start"), Handler: handleCommandStart},
		{Match: CallbackQueryHandler("start"), Handler: handleCommandStart},
	}
	ch.States[LE_NOME] = []EventHandler{
		{Match: MessageHandler(TEXT), Handler: handleLeNome},
	}
	ch.Fallbacks = []EventHandler{
		{Match: CommandHandler("cancelar"), Handler: handleCommandCancelar},
		{Match: MessageHandler(TEXT), Handler: handleCommandCancelar},
	}
	return ch
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar variáveis de ambiente: %s", err)
	}
	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Não foi possível inicializar o bot: %s", err)
	}

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	log.Println("Bot inicializado")

	ch := NewConversation()

	for update := range updates {
		go ch.HandleUpdate(bot, update)
	}
}
