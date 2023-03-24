package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"telegram-tz/models"
	"telegram-tz/repository"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getWeather(update tgbotapi.Update) tgbotapi.MessageConfig {
	URL := "http://api.openweathermap.org/data/2.5/weather?q=" + update.Message.Text + "&units=metric&APPID=" + os.Getenv("APPID")
	myClient := &http.Client{Timeout: 10 * time.Second}

	logrus.Info("Try to get '" + update.Message.Text + "' Weather")

	response, err := myClient.Get(URL)
	if err != nil {
		logrus.Fatalf("Myclient get %s", err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logrus.Fatalf("ReadAll %s", err.Error())
		}

		var result models.Response
		if err := json.Unmarshal(body, &result); err != nil {
			logrus.Fatalf("Unmarshal %s", err.Error())
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("City: %s\nTemp: %d\nFeels like: %d\nWind speed: %d \nTemp min: %d \nTemp max: %d \n", result.Name, int(result.Main.Temp), int(result.Main.FeelsLike), int(result.Wind.Speed), int(result.Main.TempMin), int(result.Main.TempMax)))
		return msg
	} else {
		logrus.Error("Status code is not 200")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error to get weather")
	return msg
}

func intiConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	if err := intiConfig(); err != nil {
		logrus.Fatalf("error config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("dotenv %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "db",
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Connection to db %s", err.Error())
	}

	repo := repository.NewRepository(db)

	if os.Getenv("CREATE_TABLE") == "yes" {
		if err := repository.CreateTable(repository.Config{
			Host:     "db",
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Password: os.Getenv("DB_PASSWORD"),
		}); err != nil {
			logrus.Fatalf("Create table %s", err.Error())
		}
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logrus.Fatalf("NewBotAPI %s", err.Error())
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			case "/start":
				user_id := int(update.Message.From.ID)

				var try int
				user, err := repo.Authorization.GetUser(user_id)
				if err != nil {
					try, err = repo.Authorization.CreateUser(user_id)
					if err != nil {
						logrus.Fatalf("Create user %s", err.Error())
					}
				}

				if user.Id != 0 {
					try = user.Id
				}

				_, err = repo.Request.CreateRequest(try, update.Message.Text)
				if err != nil {
					logrus.Fatalf("Create request %s", err.Error())
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a weather bot, i can show weather in your city.\nJust enter the name of your city \n\nCommands:\n/start\n/stats")
				bot.Send(msg)
			case "/stats":
				user_id := int(update.Message.From.ID)

				user, err := repo.Authorization.GetUser(user_id)
				if err != nil {
					logrus.Fatalf("Get user %s", err.Error())
				}

				_, err = repo.Request.CreateRequest(user.Id, update.Message.Text)
				if err != nil {
					logrus.Fatalf("Create request %s", err.Error())
				}

				len, err := repo.Request.GetRequests(user.Id)
				if err != nil {
					logrus.Fatalf("GetRequests %s", err.Error())
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("count of res: %d", len))
				bot.Send(msg)
			default:
				user_id := int(update.Message.From.ID)

				user, err := repo.Authorization.GetUser(user_id)
				if err != nil {
					logrus.Fatalf("Get user %s", err.Error())
				}

				_, err = repo.Request.CreateRequest(user.Id, update.Message.Text)
				if err != nil {
					logrus.Fatalf("Create request %s", err.Error())
				}

				msg := getWeather(update)
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}
