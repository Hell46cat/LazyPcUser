package main

import (
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const serverURL = "http://localhost:8088"

func main() {
	bot, err := tgbotapi.NewBotAPI(token) // Создаём бота с токеном
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Бот %s включен", bot.Self.UserName)  
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {  // Пропуск всех типов обновлений, кроме сообщений
			continue
		}
		// Обработка команд
		command := update.Message.Command() // Получаем команду
		switch command {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Для управления громкостью используйте команду /volume (и цифру которая будет являться процентами от общей громкости) Например /volume 50")
			bot.Send(msg)
		case "volume":
			err := volume(update, bot)
			if err != nil {
				continue
			}
		case "whatvolume":
			err := findVolume(update, bot)
			if err != nil {
				continue
			}
		case "sleep":
		case "unsleep":
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда, используйте /volume, /whatvolume, /sleep, /unsleep")
			bot.Send(msg)
		}
	}
}

// Функция работы с громкостью
func volume(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	args := update.Message.CommandArguments() // Получаем аргументы команды
	resp, err := http.Get(serverURL + "/volume?level=" + args) /// Отправка GET-запроса на сервер с параметрами
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка сервера")
		bot.Send(msg)
		return err
	}
	defer resp.Body.Close() // Освобождение ресурсов тела ответа
	response, err := io.ReadAll(resp.Body) // Чтение всего ответа
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(response))
	bot.Send(msg)

	return nil
}

// // Функция поиска громкости
func findVolume(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	resp, err := http.Get(serverURL + "/whatvolume")
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка сервера") // Отправка GET-запроса на сервер для получения текущей громкости
		bot.Send(msg)
		return err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(response))
	bot.Send(msg)

	return nil
}
