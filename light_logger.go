package light_logger

import (
    "log"
    "os"
)

// author: heyidong@bytedance.com

type lightLogger struct {
    C      chan []interface{}
    Logger *log.Logger
    File   *os.File
    Flag   int
}

func NewLogger(prefix string) *lightLogger {
    logger := &lightLogger{
        Flag:   log.Ldate | log.Ltime | log.Lshortfile,
        File:   os.Stdout,
        C:      make(chan []interface{}, 100),
        Logger: log.New(os.Stdout, prefix, log.Ldate|log.Ltime),
    }
    return logger
}
