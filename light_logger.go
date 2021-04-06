package light_logger

import (
    "log"
    "os"
)

// author: heyidong@bytedance.com

type lightLogger struct {
    C      chan []interface{}
    Logger *log.Logger
    File   *countedFile
    Flag   int
}

func newLogger(prefix string) *lightLogger {
    cntFile := &countedFile{
        File:  os.Stdout,
        Count: 0,
    }
    logger := &lightLogger{
        Flag:   log.Ldate | log.Ltime | log.Lshortfile,
        File:   cntFile,
        C:      make(chan []interface{}, 100),
        Logger: log.New(cntFile, prefix, log.Ldate|log.Ltime),
    }
    return logger
}
