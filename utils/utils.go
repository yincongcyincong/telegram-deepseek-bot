package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/volcengine/volc-sdk-golang/service/visual"
	"github.com/yincongcyincong/telegram-deepseek-bot/conf"
	"github.com/yincongcyincong/telegram-deepseek-bot/i18n"
	"github.com/yincongcyincong/telegram-deepseek-bot/logger"
)

func GetChatIdAndMsgIdAndUserID(update tgbotapi.Update) (int64, int, int64) {
	chatId := int64(0)
	msgId := 0
	userId := int64(0)
	if update.Message != nil {
		chatId = update.Message.Chat.ID
		userId = update.Message.From.ID
		msgId = update.Message.MessageID
	}
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
		userId = update.CallbackQuery.From.ID
		msgId = update.CallbackQuery.Message.MessageID
	}
	
	return chatId, msgId, userId
}

func GetChat(update tgbotapi.Update) *tgbotapi.Chat {
	if update.Message != nil {
		return update.Message.Chat
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat
	}
	return nil
}

func GetMessage(update tgbotapi.Update) *tgbotapi.Message {
	if update.Message != nil {
		return update.Message
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message
	}
	return nil
}

func GetChatType(update tgbotapi.Update) string {
	chat := GetChat(update)
	return chat.Type
}

func CheckMsgIsCallback(update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

// Utf16len calculates the length of a string in UTF-16 code units.
func Utf16len(s string) int {
	utf16Str := utf16.Encode([]rune(s))
	return len(utf16Str)
}

func ParseInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

func SendMsg(chatId int64, msgContent string, bot *tgbotapi.BotAPI, replyToMessageID int, parseMode string) {
	msg := tgbotapi.NewMessage(chatId, msgContent)
	msg.ParseMode = parseMode
	msg.ReplyToMessageID = replyToMessageID
	_, err := bot.Send(msg)
	if err != nil {
		logger.Warn("send clear message fail", "err", err)
	}
}

func ReplaceCommand(content string, command string, botName string) string {
	mention := "@" + botName
	
	content = strings.ReplaceAll(content, command, mention)
	content = strings.ReplaceAll(content, mention, "")
	prompt := strings.TrimSpace(content)
	
	return prompt
}

func ForceReply(chatId int64, msgId int, i18MsgId string, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(chatId, i18n.GetMessage(*conf.BaseConfInfo.Lang, i18MsgId, nil))
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}
	msg.ReplyToMessageID = msgId
	_, err := bot.Send(msg)
	return err
}

func GetAudioContent(update tgbotapi.Update, bot *tgbotapi.BotAPI) []byte {
	if update.Message == nil || update.Message.Voice == nil {
		return nil
	}
	
	fileID := update.Message.Voice.FileID
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		logger.Warn("get file fail", "err", err)
		return nil
	}
	
	// 构造下载 URL
	downloadURL := file.Link(bot.Token)
	
	transport := &http.Transport{}
	
	if *conf.BaseConfInfo.TelegramProxy != "" {
		proxy, err := url.Parse(*conf.BaseConfInfo.TelegramProxy)
		if err != nil {
			logger.Warn("parse proxy url fail", "err", err)
			return nil
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 设置超时
	}
	
	// 通过代理下载
	resp, err := client.Get(downloadURL)
	if err != nil {
		logger.Warn("download fail", "err", err)
		return nil
	}
	defer resp.Body.Close()
	voice, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn("read response fail", "err", err)
		return nil
	}
	return voice
}

func GetPhotoContent(update tgbotapi.Update, bot *tgbotapi.BotAPI) []byte {
	if update.Message == nil || update.Message.Photo == nil {
		return nil
	}
	
	var photo tgbotapi.PhotoSize
	for i := len(update.Message.Photo) - 1; i >= 0; i-- {
		if update.Message.Photo[i].FileSize < 8*1024*1024 {
			photo = update.Message.Photo[i]
			break
		}
	}
	
	fileID := photo.FileID
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		logger.Warn("get file fail", "err", err)
		return nil
	}
	
	downloadURL := file.Link(bot.Token)
	
	client := GetTelegramProxyClient()
	resp, err := client.Get(downloadURL)
	if err != nil {
		logger.Warn("download fail", "err", err)
		return nil
	}
	defer resp.Body.Close()
	photoContent, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn("read response fail", "err", err)
		return nil
	}
	
	return photoContent
}

func MD5(input string) string {
	// 计算 MD5
	hash := md5.Sum([]byte(input))
	
	// 转换为 16 进制字符串
	md5Str := hex.EncodeToString(hash[:])
	return md5Str
}

func GetTelegramProxyClient() *http.Client {
	transport := &http.Transport{}
	
	if *conf.BaseConfInfo.TelegramProxy != "" {
		proxy, err := url.Parse(*conf.BaseConfInfo.TelegramProxy)
		if err != nil {
			logger.Warn("parse proxy url fail", "err", err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	
	return &http.Client{
		Transport: transport,
	}
}

func GetDeepseekProxyClient() *http.Client {
	transport := &http.Transport{}
	
	if *conf.BaseConfInfo.DeepseekProxy != "" {
		proxy, err := url.Parse(*conf.BaseConfInfo.DeepseekProxy)
		if err != nil {
			logger.Warn("parse proxy url fail", "err", err)
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	
	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Minute, // 设置超时
	}
}

func CreateBot() *tgbotapi.BotAPI {
	// 配置自定义 HTTP Client 并设置代理
	client := GetTelegramProxyClient()
	
	var err error
	conf.BaseConfInfo.Bot, err = tgbotapi.NewBotAPIWithClient(*conf.BaseConfInfo.TelegramBotToken, tgbotapi.APIEndpoint, client)
	if err != nil {
		panic("Init bot fail" + err.Error())
	}
	
	if *logger.LogLevel == "debug" {
		conf.BaseConfInfo.Bot.Debug = true
	}
	
	// set command
	cmdCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "help",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.help.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "clear",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.clear.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "retry",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.retry.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "mode",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.mode.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "balance",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.balance.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "state",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.state.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "photo",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.photo.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "video",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.video.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "chat",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.chat.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "task",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.task.description", nil),
		},
		tgbotapi.BotCommand{
			Command:     "mcp",
			Description: i18n.GetMessage(*conf.BaseConfInfo.Lang, "commands.mcp.description", nil),
		},
	)
	conf.BaseConfInfo.Bot.Send(cmdCfg)
	
	return conf.BaseConfInfo.Bot
}

func GetContent(update tgbotapi.Update, bot *tgbotapi.BotAPI, content string) (string, error) {
	// check user chat exceed max count
	if CheckUserChatExceed(update, bot) {
		return "", errors.New("token exceed")
	}
	
	if content == "" && update.Message.Voice != nil && *conf.AudioConfInfo.AudioAppID != "" {
		audioContent := GetAudioContent(update, bot)
		if audioContent == nil {
			logger.Warn("audio url empty")
			return "", errors.New("audio url empty")
		}
		content = FileRecognize(audioContent)
	}
	
	if content == "" && update.Message.Photo != nil {
		imageContent, err := GetImageContent(GetPhotoContent(update, bot))
		if err != nil {
			logger.Warn("get image content err", "err", err)
			return "", err
		}
		content = imageContent
	}
	
	if content == "" {
		logger.Warn("content empty")
		return "", errors.New("content empty")
	}
	
	text := strings.ReplaceAll(content, "@"+bot.Self.UserName, "")
	return text, nil
}

func FileRecognize(audioContent []byte) string {
	
	client := BuildAsrClient()
	client.Appid = *conf.AudioConfInfo.AudioAppID
	client.Token = *conf.AudioConfInfo.AudioToken
	client.Cluster = *conf.AudioConfInfo.AudioCluster
	
	asrResponse, err := client.RequestAsr(audioContent)
	if err != nil {
		logger.Error("fail to request asr ", "err", err)
		return ""
	}
	
	if len(asrResponse.Results) == 0 {
		logger.Error("fail to request asr", "results", asrResponse.Results)
		return ""
	}
	
	return asrResponse.Results[0].Text
	
}

func GetImageContent(imageContent []byte) (string, error) {
	visual.DefaultInstance.Client.SetAccessKey(*conf.BaseConfInfo.VolcAK)
	visual.DefaultInstance.Client.SetSecretKey(*conf.BaseConfInfo.VolcSK)
	
	form := url.Values{}
	form.Add("image_base64", base64.StdEncoding.EncodeToString(imageContent))
	
	resp, _, err := visual.DefaultInstance.OCRNormal(form)
	if err != nil {
		logger.Error("request img api fail", "err", err)
		return "", err
	}
	
	if resp.Code != 10000 {
		logger.Error("request img api fail", "code", resp.Code, "msg", resp.Message)
		return "", errors.New("request img api fail")
	}
	
	return strings.Join(resp.Data.LineTexts, ","), nil
}

func FileToMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	hash := md5.New()
	
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	
	md5sum := fmt.Sprintf("%x", hash.Sum(nil))
	return md5sum, nil
}

func SetStructFieldByJSONTag(s interface{}, key string, value interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("want a struct，get a %T", s)
	}
	
	elem := val.Elem()
	typ := elem.Type()
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := elem.Field(i)
		
		jsonTag := field.Tag.Get("json")
		
		if jsonTag == key {
			if !fieldValue.CanSet() {
				return fmt.Errorf("field %s can't be set", field.Name)
			}
			
			setValue := reflect.ValueOf(value)
			if fieldValue.Kind() == reflect.Ptr {
				elemType := fieldValue.Type().Elem()
				
				// 检查传入的值是否可转换为元素类型
				if !setValue.Type().ConvertibleTo(elemType) {
					return fmt.Errorf("can't put %T change to field %s with %s type", value, field.Name, elemType)
				}
				
				newValue := reflect.New(elemType)
				newValue.Elem().Set(setValue.Convert(elemType))
				
				fieldValue.Set(newValue)
				return nil
			} else {
				if setValue.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(setValue.Convert(fieldValue.Type()))
					return nil
				} else {
					return fmt.Errorf("can't put %T change to field %s with %s type", value, field.Name, fieldValue.Type())
				}
			}
		}
	}
	
	return fmt.Errorf("can't find with key '%s' matched JSON field", key)
}

func ValueToString(value interface{}) string {
	if value == nil {
		return "nil" // Handle nil values
	}
	
	// Get the reflect.Value and reflect.Type of the input
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)
	
	// Handle pointer types: if it's a pointer, get the element it points to
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil"
		}
		v = v.Elem() // Get the value the pointer points to
		t = v.Type() // Update type to the actual type
	}
	
	switch t.Kind() {
	case reflect.String:
		return v.String() // Return string directly
	
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10) // Convert integer to decimal string
	
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10) // Convert unsigned integer to decimal string
	
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64) // Format float as string
		// 'f' means decimal format, -1 uses necessary precision, 64 for float64
	
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()) // Convert boolean to string
	
	case reflect.Array, reflect.Slice:
		// For arrays and slices, convert each element to a string and join with commas
		elements := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			elements[i] = ValueToString(v.Index(i).Interface()) // Recursive call for nested types
		}
		return strings.Join(elements, ",")
	
	case reflect.Map:
		// For maps, convert each key-value pair to string
		keys := v.MapKeys()
		parts := make([]string, len(keys))
		for i, key := range keys {
			val := v.MapIndex(key)
			parts[i] = fmt.Sprintf("%s: %s", ValueToString(key.Interface()), ValueToString(val.Interface()))
		}
		return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
	
	case reflect.Struct:
		// For structs, convert each field to a string
		parts := make([]string, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name
			parts[i] = fmt.Sprintf("%s: %s", fieldName, ValueToString(field.Interface()))
		}
		return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
	
	case reflect.Invalid:
		return "<invalid value>" // Handle invalid values such as uninitialized reflect.Value
	
	default:
		// For other types like Channel, Func, Interface, Ptr (already handled), UnsafePointer, etc.
		// Use fmt.Sprintf to get the default string representation
		return fmt.Sprintf("%v", value)
	}
}

func MapKeysToString(input interface{}) string {
	val := reflect.ValueOf(input)
	
	// Check if it's a map
	if val.Kind() != reflect.Map {
		return ""
	}
	
	keys := val.MapKeys()
	keyStrs := make([]string, len(keys))
	
	for i, key := range keys {
		// Convert key to string using fmt.Sprint (handles int, string, etc.)
		keyStrs[i] = fmt.Sprint(key.Interface())
	}
	
	return strings.Join(keyStrs, ",")
}
