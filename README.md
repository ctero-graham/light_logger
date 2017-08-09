# light_logger
light_logger is a very light library which supports rotation of log files. So far, it can only rotate log files by day.

Usage
=====
Simple logging like this:
```golang
light_logger.Error(msg...)    //  log error level messages
light_logger.Debug(msg...)    //  log debug level messages
light_logger.Info(msg...)     //  log info level messages
light_logger.Warn(msg...)     //  log warn level messages
```

By default, the messages will be print to stdout, If you would like to specify a file to log to:
 ```golang
light_logger.SetLogFile(level, filePath, prefix, flag)
```

So far, light_logger only support log rotation in one day period. If you do not like the rotation:
```golang
light_logger.Rotates = false    //  rotation will be forbidden
```