package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/SolidShake/wetherboy-tg-bot/iternal/config"
	"github.com/SolidShake/wetherboy-tg-bot/iternal/connections"
	"github.com/SolidShake/wetherboy-tg-bot/iternal/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	dbConnection := connections.MongoConnection{}
	dbConnection.ConnectMongo()
	fmt.Println("Connected to MongoDB!")
	//fmt.Println("Connected to database:"+m.GetVersion())
	fmt.Printf("Connected to database:%s", dbConnection.GetDbName())
	defer dbConnection.Disconnect()

	bot, err := tgbotapi.NewBotAPI(config.GetConfig().Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		//Ограничить ввод
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		btn := tgbotapi.KeyboardButton{
			RequestLocation: true,
			Text:            "Обновить свою геолокацию",
		}

		subButton := tgbotapi.KeyboardButton{
			RequestLocation: true,
			Text:            "Подписаться на прогноз",
		}

		//unsubButton := tgbotapi.KeyboardButton{
		//	Text: "Подписаться на прогноз",
		//}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			case "Подписаться на прогноз":
				dbConnection.AddSub(update.Message.Chat.ID, *update.Message.Location)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для получения актуальной погоды нажмите кнопку ниже")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
				bot.Send(msg)
			}
		}

		if update.Message.Location != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, getWeatherInfoByCoord(update.Message.Location.Latitude, update.Message.Location.Longitude))
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn, subButton})
			bot.Send(msg)
			dbConnection.AddSub(update.Message.Chat.ID, *update.Message.Location)
		}
	}
}

func getWeatherInfoByCoord(latitude, longitude float64) string {
	requestUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=709f4fa20539d7cecc587597bd8e417c&units=metric&&lang=ru", latitude, longitude)
	r, err := http.Get(requestUrl)
	// Добавить обработку
	if err != nil {
		log.Println("Request failed")
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	var response types.RequestStruct
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}

	result := fmt.Sprintf(
		"Доброго времени суток!\n\nВы находитесь в городе %s\nТемпература: %.2f ℃\nОщущается как %.2f ℃\nОжидаемая погода: %s\nСкорость ветра: %.2f м/с",
		response.Name,
		response.Main.Temp,
		response.Main.FeelsLike,
		response.Weather[0].Description,
		response.Wind.Speed,
	)

	return result
}
