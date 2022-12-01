package messenger

import (
	"errors"

	tlg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	tokenApi string
	bot      *tlg.BotAPI
	config   map[string]string
}

var (
	InvalidRoomidError      = errors.New("roomid is invalid type")
	InvalidTextMessageError = errors.New("Message text is invalid type")
	InvalidMsgIdError       = errors.New("Message Id is invalid type")
)

func NewTelegramBot(apiToken string) *TelegramBot {
	bot, err := tlg.NewBotAPI(apiToken)
	if err != nil {
		l.Error(err.Error())
		return nil
	}

	return &TelegramBot{
		tokenApi: apiToken,
		bot:      bot,
	}
}
func (t TelegramBot) GetTokenApi() string {
	return t.tokenApi
}

func (t TelegramBot) Send(msg Message) error {

	roomId, ok := msg["roomId"].(int64)
	if !ok {
		l.Error(InvalidRoomidError.Error())
		return InvalidRoomidError
	}
	text, ok := msg["text"].(string)
	if !ok {
		l.Error(InvalidTextMessageError.Error())
		return InvalidTextMessageError
	}

	tlgMessage := tlg.NewMessage(roomId, text)
	msgId, ok := msg["msgId"]
	if ok {
		repMsgId, ok := msgId.(int)
		if !ok {
			l.Error(InvalidMsgIdError.Error())
			return InvalidMsgIdError
		}
		tlgMessage.ReplyToMessageID = repMsgId
	}

	_, err := t.bot.Send(tlgMessage)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	return nil
}

func (t TelegramBot) Listen(done chan interface{}) (chan Message, chan interface{}) {
	msgStream := make(chan Message)
	terminated := make(chan interface{})

	updateConfig := tlg.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := t.bot.GetUpdatesChan(updateConfig)

	go func() {
		for {
			select {
			case update := <-updates:
				msg := make(Message)
				msg["roomId"] = update.Message.Chat.ID
				msg["text"] = update.Message.Text
				msg["msgId"] = update.Message.MessageID
				msgStream <- msg
				l.Debug(update.Message.Text)
			case <-done:
				l.Info("Terminated message bot")
				close(terminated)
			}
		}
	}()

	return msgStream, terminated
}
