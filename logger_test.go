package light_logger

import (
    "testing"
    "time"
    "fmt"
    "log"
)

// author: heyd

func TestNextTimerDuration(t *testing.T) {
    tm := time.Date(2017, 1, 1, 23, 59, 59, 0, time.Local)
    fmt.Println(tm)
    next := nextTimerDuration(tm)
    fmt.Println(next)
    if next.Year() != 2017 || next.Month() != 1 || next.Day() != 2 ||
            next.Hour()+next.Minute()+next.Second()+next.Nanosecond() != 0 {
        t.Fatal()
    }

    tm = time.Date(2017, 12, 31, 0, 0, 0, 0, time.Local)
    fmt.Println(tm)
    next = nextTimerDuration(tm)
    fmt.Println(next)
    if next.Year() != 2018 || next.Month() != 1 || next.Day() != 1 ||
            next.Hour()+next.Minute()+next.Second()+next.Nanosecond() != 0 {
        t.Fatal()
    }

    tm = time.Date(2000, 2, 28, 3, 13, 0, 0, time.Local)
    fmt.Println(tm)
    next = nextTimerDuration(tm)
    fmt.Println(next)
    if next.Year() != 2000 || next.Month() != 2 || next.Day() != 29 ||
            next.Hour()+next.Minute()+next.Second()+next.Nanosecond() != 0 {
        t.Fatal()
    }

    tm = time.Date(2001, 2, 28, 0, 1, 0, 0, time.Local)
    fmt.Println(tm)
    next = nextTimerDuration(tm)
    fmt.Println(next)
    if next.Year() != 2001 || next.Month() != 3 || next.Day() != 1 ||
            next.Hour()+next.Minute()+next.Second()+next.Nanosecond() != 0 {
        t.Fatal()
    }
}

func TestRotate(t *testing.T) {
    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    if SetLogFile(8, "debug.log", "", 0) == nil {
        t.Fatal()
    }
    time.Sleep(time.Second)

    if SetLogFile(Debug, "debug.log", "", log.Llongfile|log.Ltime|log.Ldate) != nil {
        t.Fatal()
    }

    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    time.Sleep(time.Second)

    if SetLogFile(Info, "info.log", "", 0) != nil {
        t.Fatal()
    }

    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    time.Sleep(time.Second)

    rotate()

    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    Error("TAG", "this", "is", 10086)
    Warn("TAG", "this", "is", 10086)
    Info("TAG", "this", "is", 10086)
    Debug("TAG", "this", "is", 10086)

    time.Sleep(time.Second)
}
