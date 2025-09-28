package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Usuario struct {
	State      int64
	ID         int64
	Name       string
	Username   string
	Parameters map[string]string
}

func StartKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Nova busca", "start"),
		),
	)
}

func handleCommandStart(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *Usuario) int64 {
	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		bot.Request(callback)
	}
	log.Println("Tratando comando start")

	msg := tgbotapi.NewMessage(user.ID, "Vamos lá, me diga o nome do Pokemon que você deseja pesquisar")
	bot.Send(msg)

	return LE_NOME
}

func handleLeNome(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *Usuario) int64 {
	log.Println("lendo nome do pokemon")
	pokemon := update.Message.Text
	log.Printf("Pokemon informado: %s\n", pokemon)
	user.Parameters["pokemonName"] = pokemon
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Você pesquisou por: %s", pokemon))
	bot.Send(msg)

	imgURL := search(user.Parameters["pokemonName"])
	msgPhoto := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FileURL(imgURL))
	bot.Send(msgPhoto)

	msgNovaBusca := tgbotapi.NewMessage(user.ID, "Use o botão abaixo para iniciar uma nova pesquisa:")
	msgNovaBusca.ReplyMarkup = StartKeyboard()
	bot.Send(msgNovaBusca)

	return END
}

func handleCommandCancelar(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *Usuario) int64 {
	log.Println("Tratando comando cancelar")
	return END
}
