package main

import (
	"fmt"
	"log"
	"strings"

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
	pokemonName := update.Message.Text
	user.Parameters["pokemonName"] = pokemonName
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Você pesquisou por: %s", strings.ToUpper(pokemonName)))
	bot.Send(msg)

	fmt.Println(user)

	pokemon := search(user.Parameters["pokemonName"])
	msgPhoto := tgbotapi.NewPhoto(update.FromChat().ID, tgbotapi.FileURL(pokemon.Image))
	msgPhoto.Caption = fmt.Sprintf("📕 Nome: %s", strings.ToUpper(pokemonName))
	bot.Send(msgPhoto)

	msgAbilities := tgbotapi.NewMessage(update.FromChat().ID, pokemon.FormatAbilities())
	msgAbilities.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msgAbilities)

	msgNovaBusca := tgbotapi.NewMessage(user.ID, "Use o botão abaixo para iniciar uma nova pesquisa:")
	msgNovaBusca.ReplyMarkup = StartKeyboard()
	bot.Send(msgNovaBusca)

	return END
}

func handleCommandCancelar(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *Usuario) int64 {
	log.Println("Tratando comando cancelar")
	return END
}
