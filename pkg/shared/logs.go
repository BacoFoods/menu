package shared

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

func LogError(message, directory, function string, error error, args ...any) {
	logrus.WithFields(logrus.Fields{
		"directory": directory,
		"function":  function,
		"args":      ToStringArgs(args),
		"error":     error,
	}).Error(message)
}

func LogWarn(message, directory, function string, error error, args ...any) {
	logrus.WithFields(logrus.Fields{
		"directory": directory,
		"function":  function,
		"args":      ToStringArgs(args),
		"error":     error,
	}).Warning(message)
}

func LogInfo(message, directory, function string, error error, args ...any) {
	logrus.WithFields(logrus.Fields{
		"directory": directory,
		"function":  function,
		"args":      ToStringArgs(args),
		"error":     error,
	}).Info(message)
}

func LogRequest(uuid, method, url, body string, request, headers any) {
	logrus.WithFields(logrus.Fields{
		"method":  method,
		"url":     url,
		"body":    body,
		"headers": headers,
	}).Infof("HTTP_REQUEST=%s", uuid)
}

func LogResponse(uuid, status, body, method, url string, headers any) {
	logrus.WithFields(logrus.Fields{
		"method":  method,
		"url":     url,
		"status":  status,
		"body":    body,
		"headers": headers,
	}).Infof("HTTP_RESPONSE=%s", uuid)
}

func ToStringArgs(args ...any) string {
	var stringArgs []string
	for _, arg := range args {
		fmt.Printf("%T", arg)
		switch v := arg.(type) {
		case string:
			stringArgs = append(stringArgs, v)
		case []any:
			jsonMap, err := json.Marshal(arg)
			if err != nil {
				stringArgs = append(stringArgs, fmt.Sprintf("%+v", arg))
			} else {
				stringArgs = append(stringArgs, string(jsonMap))
			}
		default:
			stringArgs = append(stringArgs, fmt.Sprintf("%+v", arg))
		}
	}
	return fmt.Sprintf("%+v", stringArgs)
}
