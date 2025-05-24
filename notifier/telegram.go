package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func escapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, `\`+char)
	}
	return text
}

func SendTelegramMessage(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("не указан TELEGRAM_BOT_TOKEN или TELEGRAM_CHAT_ID")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", escapeMarkdown(message))
	data.Set("parse_mode", "MarkdownV2")

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось отправить сообщение, статус: %s", resp.Status)
	}

	return nil
}
