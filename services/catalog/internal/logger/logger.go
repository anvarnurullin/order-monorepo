package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Log *log.Logger
var serviceName string

func Init() {
	serviceName = "catalog-service"
	Log = log.New(os.Stdout, "", 0)
}

func timestamp() string {
    return time.Now().Format("02/01/2006 15:04:05")
}

func Info(msg string) {
    Log.Printf("%s - %s - [INFO] %s\n", serviceName, timestamp(), msg)
}

func Infof(format string, a ...any) {
    Log.Printf("%s - %s - [INFO] %s\n", serviceName, timestamp(), formatMessage(format, a...))
}

func Error(msg string, err error) {
    if err != nil {
        Log.Printf("%s - %s - [ERROR] %s: %s\n", serviceName, timestamp(), msg, err.Error())
    } else {
        Log.Printf("%s - %s - [ERROR] %s\n", serviceName, timestamp(), msg)
    }
}

func Debug(msg string) {
    Log.Println("[DEBUG] " + msg)
}

func formatMessage(format string, a ...any) string {
    return fmt.Sprintf(format, a...)
}
