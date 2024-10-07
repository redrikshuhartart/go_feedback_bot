package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const adminId int64 = "YOUR ADMIN ID"

func checkUserInGroup(bot *tgbotapi.BotAPI, chatId int64, userId int64) bool {

	config := tgbotapi.ChatConfigWithUser{
		ChatID: chatId,
		UserID: userId,
	}

	configNew := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: config,
	}

	result, err := bot.GetChatMember(configNew)
	if err != nil {
		log.Fatal(err)
	}

	if result.Status == "member" || result.Status == "creator" {
		return true
	}

	return false
}

func main() {
	// Айди чата администраторов
	chatID := int64("YOUR CHAT ADMIN ID")
	// Мапа для хранения данных для сопоставления [Имя пользователя - айди пользователя]
	userDataMap := make(map[string]int64)
	var userToSendMessage string

	bot, err := tgbotapi.NewBotAPI("YOUR TELEGRAM TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msgToSend := tgbotapi.NewForward(chatID, update.Message.From.ID, update.Message.MessageID)

		// Складываем в мапу данные - Имя пользователя - Айди пользователя
		if update.Message.From != nil && update.Message.From.LastName != "" {
			userDataMap[update.Message.From.FirstName+" "+update.Message.From.LastName] = update.Message.From.ID
		} else if update.Message.From.FirstName == "" && update.Message.From.LastName != "" {
			userDataMap[update.Message.From.LastName] = update.Message.From.ID
		} else {
			userDataMap[update.Message.From.FirstName] = update.Message.From.ID
		}

		fmt.Println(userDataMap)

		if update.Message.Chat.Type == "supergroup" {
			if update.Message.ReplyToMessage != nil {
				if update.Message.ReplyToMessage.ForwardFrom != nil {
					if update.Message.ReplyToMessage.ForwardFrom.LastName != "" {
						userToSendMessage = update.Message.ReplyToMessage.ForwardFrom.FirstName + " " + update.Message.ReplyToMessage.ForwardFrom.LastName
					} else if update.Message.ReplyToMessage.ForwardFrom.FirstName == "" && update.Message.ReplyToMessage.ForwardFrom.LastName != "" {
						userToSendMessage = update.Message.ReplyToMessage.ForwardFrom.LastName
					} else {
						userToSendMessage = update.Message.ReplyToMessage.ForwardFrom.FirstName
					}
				} else {
					userToSendMessage = update.Message.ReplyToMessage.ForwardSenderName
				}

				fmt.Println(userToSendMessage)
				userID := userDataMap[userToSendMessage]
				fmt.Println(userID)
				msgToSend := tgbotapi.NewMessage(userID, update.Message.Text)
				_, err := bot.Send(msgToSend)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					startMessage := tgbotapi.NewMessage(update.Message.From.ID, "Добро пожаловать в бот обратной связи. Напишите нам и мы обязательно с Вами свяжемся.")
					bot.Send(startMessage)
				}
			} else {
				// Отправка сообщения в другой чат
				msgToUserSend := tgbotapi.NewMessage(update.Message.From.ID, "Спасибо, ваше сообщение отправлено администратору!")
				_, err := bot.Send(msgToUserSend)
				if err != nil {
					log.Println(err)
				}
				_, err = bot.Send(msgToSend)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	defer func() {
		msgToAdmin := tgbotapi.NewMessage(adminId, "Я упал. Пытаюсь встать.")
		bot.Send(msgToAdmin)
		fmt.Println("[Произошло аварийное завершение работы бота.]")
	}()
}
