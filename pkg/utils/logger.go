package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

func FormatKeyValueLogs(data [][2]string) string {
	var builder strings.Builder
	builder.Grow(len(data) * 10)

	for _, entry := range data {
		builder.WriteString(fmt.Sprintf("  %s: %s\n", entry[0], entry[1]))
	}

	return builder.String()
}

func LogInfo(title string, message string) {
	if logrus.GetLevel() < logrus.InfoLevel {
		return
	}

	logrus.Info(fmt.Sprintf(
		"\033[1m%s\033[0m:\n%s",
		title,
		message,
	))
}

func LogError(message string, errStr string) {
	logrus.Error(fmt.Sprintf(
		"%s: \033[38;5;197m%s\033[0m",
		message,
		errStr,
	))
}

func PrintStructFields(params interface{}) {
	val := reflect.ValueOf(params)

	// Ensure the value is a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Println("Expected a struct")
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Check if it's a nested struct
		if field.Kind() == reflect.Struct {
			fmt.Printf("\n%s:\n", fieldName)
			PrintStructFields(field.Interface()) // Recursively print nested struct fields
			fmt.Println()
		} else {
			fmt.Printf("\n%s: %v", fieldName, field.Interface()) // Print field value
		}
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("Error (Code: %d, Message: %s)", e.Code, e.Message)
}

func ErrMalformedRequest(message string) error {
	origin := GetOrigin()

	return Error{
		Code:    400,
		Message: "Malformed request",
		Details: message,
		Origin:  origin,
	}
}

func ErrInternal(message string) Error {
	origin := GetOrigin()

	return Error{
		Code:    500,
		Message: "Internal server error",
		Details: message,
		Origin:  origin,
	}
}
