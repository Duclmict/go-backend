package log_service

import (
    "fmt"
    "log"
    "time"
    "os"
    "runtime"

    "github.com/Duclmict/go-backend/config"
)

// private 
func writeLog(mode string, message string) {

	now := time.Now()
    cur_date := now.Format("2006-01-02")

    f, err := os.OpenFile(config.LOG_FOLDER + "/" + cur_date + ".log", 
        os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    
    log.SetOutput(f)

    // set 2 because run 2 time in package log service
    _, file, line, _ := runtime.Caller(2)
    
    log.Println(fmt.Sprintf("[%s] %s  {%s %d }", mode, message, file, line))
}

// public
func Debug(message string) {
    if(config.App_Debug == "True") {
        writeLog("DEBUG", message)
    }  
}

func Info(message string) {
    writeLog("INFO", message)
}

func Error(message string) {
    writeLog("ERROR", message)
}

func Warning(message string) {
    writeLog("WARNING", message)
}