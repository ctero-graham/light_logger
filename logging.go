//  light_logger is a very light library which supports rotation of log files. So far, it can only rotate log files
//  by day.
package light_logger

import (
    "os"
    "time"
    "fmt"
    "runtime"
    "log"
)

// author: heyd

var (
    //  Rotates decides whether log file will be rotated.
    Rotates = true

    //  PostfixFormat specifies the postfix of the log file which is to be rotated, it is supposed to conform to
    //  the layout format in golang's time package.
    PostfixFormat = "20060102"
)

const (
    LEVEL_DEBUG = iota
    LEVEL_INFO
    LEVEL_WARN
    LEVEL_ERROR
)

type InvalidLogLevel int

func (InvalidLogLevel) Error() string {
    return "invalid log level."
}

var (
    debugLogger *lightLogger
    infoLogger  *lightLogger
    warnLogger  *lightLogger
    errorLogger *lightLogger

    logsMap map[int]*lightLogger

    timer *time.Timer
)

//  init creates channels as well as loggers, and then start a goroutine to do the real logging.
func init() {
    //  log to stdout by default.
    debugLogger = newLogger("[DEBUG]")
    infoLogger = newLogger("[INFO]")
    warnLogger = newLogger("[WARN]")
    errorLogger = newLogger("[ERROR]")

    logsMap = map[int]*lightLogger{
        LEVEL_DEBUG: debugLogger,
        LEVEL_INFO:  infoLogger,
        LEVEL_WARN:  warnLogger,
        LEVEL_ERROR: errorLogger,
    }

    go startLogging()
}

//  startLogging starts to listening to the channels and do the logging, it is supposed to be called in a goroutine.
func startLogging() {
    fmt.Println("timer started")
    setTimer()
    for {
        select {
        case logs := <-debugLogger.C:
            debugLogger.Logger.Println(logs...)
        case logs := <-infoLogger.C:
            infoLogger.Logger.Println(logs...)
        case logs := <-warnLogger.C:
            warnLogger.Logger.Println(logs...)
        case logs := <-errorLogger.C:
            errorLogger.Logger.Println(logs...)
        case <-timer.C:
            fmt.Println("rotates started")
            rotate()
        }
    }
}

//  SetLogFile specifies a file path to log to, a prefix and a flag of the logger for a log level.
//  The available value of flag is same as that in golang's log package.
func SetLogFile(level int, filePath, prefix string, flag int) error {
    logger, ok := logsMap[level]
    if !ok {
        var e InvalidLogLevel
        return e
    } else {
        f, e := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
        if e != nil {
            fmt.Errorf("open log file error: %s", e.Error())
            return e
        }
        oldFile := logger.File
        if oldFile != nil && oldFile != os.Stdout {
            oldFile.Close()
        }
        resetLogger(logger, f, prefix, flag)
        return nil
    }
}

//  resetLogger specifies a file to log to, a prefix and a flag of the logger for a logger.
func resetLogger(logger *lightLogger, f *os.File, prefix string, flag int) {
    logger.File = f
    logger.Logger.SetOutput(f)
    if prefix != "" {
        logger.Logger.SetPrefix(prefix)
    }
    if flag != 0 {
        logger.Flag = flag
        if flag&log.Llongfile != 0 {
            flag = flag & (^log.Llongfile)
        }
        if flag&log.Lshortfile != 0 {
            flag = flag & (^log.Lshortfile)
        }
        logger.Logger.SetFlags(flag)
    }
}

//  Error logs error messages.
func Error(data ... interface{}) {
    logTo(errorLogger, data...)
}

//  Warn logs warn messages.
func Warn(data ... interface{}) {
    logTo(warnLogger, data...)
}

//  Debug logs debug messages.
func Debug(data ... interface{}) {
    logTo(debugLogger, data...)
}

//  Info logs info messages.
func Info(data ... interface{}) {
    logTo(infoLogger, data...)
}

//  logTo logs a given message to the specified logger.
func logTo(logger *lightLogger, data ... interface{}) {
    if logger.Flag&(log.Lshortfile|log.Llongfile) != 0 {
        _, file, line, _ := runtime.Caller(2)
        if logger.Flag&log.Lshortfile != 0 {
            short := file
            for i := len(file) - 1; i > 0; i-- {
                if file[i] == '/' {
                    short = file[i+1:]
                    break
                }
            }
            file = short
        }
        logger.C <- append([]interface{}{fmt.Sprintf("%s:%d", file, line)}, data...)
    } else {
        logger.C <- data
    }
}

//  rotate will do the log rotation work.
func rotate() {
    if !Rotates {
        return
    }
    now := time.Now().Format(PostfixFormat)
    for k, v := range logsMap {
        if v.File != os.Stdout && v.File != nil {
            filePath := v.File.Name()
            v.File.Close()
            v.File = nil
            os.Rename(filePath, filePath+"."+now)
            SetLogFile(k, filePath, v.Logger.Prefix(), v.Flag)
        }
    }
    setTimer() //  定下一次的闹钟
}

//  setTimer will call setNextTimer to set a timer which will trigger the next rotation work.
func setTimer() {
    now := time.Now()
    next := nextTimerDuration(now)
    duration := next.Sub(now)
    timer = setNextTimer(timer, duration)
}

//  setNextTimer sets or resets a timer to trigger for a specified duration.
func setNextTimer(timer *time.Timer, duration time.Duration) *time.Timer {
    if timer == nil {
        timer = time.NewTimer(duration)
    } else {
        timer.Reset(duration)
    }
    return timer
}

//  nextTimerDuration calculates the next time a rotation work will be performed given a specified time.
func nextTimerDuration(since time.Time) time.Time {
    tomorrow := since.Add(time.Hour * 24)
    next := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
    return next
}
