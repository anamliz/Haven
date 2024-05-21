package main

import (
    "context"
    "io"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/fsnotify/fsnotify"
    "github.com/anamliz/Haven/internal/services/pollData"
    "github.com/spf13/viper"
    "gopkg.in/natefinch/lumberjack.v2"
)

var defaultHTTPClient = &http.Client{
    Timeout: time.Second * 15,
    Transport: &http.Transport{
        Dial: (&net.Dialer{
            Timeout: time.Second * 15,
        }).Dial,
        TLSHandshakeTimeout: 10 * time.Second,
    },
}

var addConfigPathLive = `C:\apps\go\Haven\cmd\poll_data\`
var addConfigPathLocal = `C:\apps\go\Haven\cmd\poll_data\`

var inProgress bool

func main() {
    InitConfig()

    pd, err := pollData.NewPollDataService(
        pollData.WithMysqlPollDataRepository(viper.GetString("mySQL.live")),
    )
    if err != nil {
        log.Printf("Unable to start poll data service : %s", err)
        os.Exit(1)
    }

    pollDataEndPoint := viper.GetString("poll_data.url")
    timeouts := viper.GetDuration("poll_data.timeouts")

    ctx := context.Background()

    log.Printf("About to start")

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    go func() {
        for t := range ticker.C {
            if !inProgress {
                inProgress = true
                go getData(ctx, pd, pollDataEndPoint, timeouts)
            } else {
                log.Printf("Poll data still running: %v", t)
            }
        }
    }()

    sig := make(chan os.Signal, 1)
    defer close(sig)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

    s := <-sig
    log.Printf("Caught signal and exiting: %v", s)
}

func getData(ctx context.Context, pd *pollData.PollDataService, pollDataEndPoint string, timeouts time.Duration) {
    defer func() {
        inProgress = false
        log.Printf("Done processing poll data")
    }()

    err := pd.PollData(ctx, pollDataEndPoint, timeouts, defaultHTTPClient)
    if err != nil {
        log.Printf("Error: %v", err)
    }
}

func InitConfig() {
    configUtils(addConfigPathLive, addConfigPathLocal)
    logUtils(viper.GetString("poll_data.logs"), viper.GetInt("log_setting.MaxSize"),
        viper.GetInt("log_setting.MaxBackups"), viper.GetInt("log_setting.MaxAge"),
        viper.GetBool("log_setting.Compress"))
}

func logUtils(logDirectory string, maxSize int, maxBackups int, maxAge int, compress bool) {
    logFile := &lumberjack.Logger{
        Filename:   logDirectory,
        MaxSize:    maxSize,
        MaxBackups: maxBackups,
        MaxAge:     maxAge,
        Compress:   compress,
    }
    log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func configUtils(addConfigPathLive string, addConfigPathLocal string) {
    log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
    viper.SetDefault("host", "localhost")
    viper.SetConfigName("config")
    viper.AddConfigPath(addConfigPathLive)
    viper.AddConfigPath(addConfigPathLocal)
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        log.Printf("Error: %v", err)
    }

    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.Name)
    })
}
