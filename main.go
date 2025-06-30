package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/Jaezmien/wlr/jsonl"
)

type WlanEntry struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
}

var (
	entries []WlanEntry
	force   bool
	debug   bool
)

func init() {
	flag.BoolVar(&force, "force", false, "Force rescanning")
	flag.BoolVar(&debug, "debug", false, "Show debug logs")
	flag.Parse()
}

func init() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if _, err := exec.LookPath("nmcli"); err != nil {
		slog.Error("could not find 'nmcli'")
		os.Exit(1)
	}

	if _, err := os.Stat(("./wlan.jsonl")); err != nil {
		slog.Error("'./wlan.jsonl' file missing")
		os.Exit(1)
	}

	data, err := os.ReadFile("./wlan.jsonl")
	if err != nil {
		slog.Error("error while trying to read wlan file", slog.Any("error", err))
		os.Exit(1)
	}

	if err := jsonl.Unmarshal(data, &entries); err != nil {
		slog.Error("error while trying to unmarshal wlan file", slog.Any("error", err))
		os.Exit(1)
	}

	if len(entries) == 0 {
		slog.Error("wlan entries are empty")
		os.Exit(1)
	}
}

func ChangeWlan(ssid string, password string) (bool, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
	defer c()

	cmd := exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", ssid, "password", password)
	_, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error while trying to change wlan: %v", err)
	}
	return true, nil
}

func TestConnection() bool {
	req, err := http.NewRequest("GET", "https://www.gstatic.com/generate_204", nil)
	if err != nil {
		return false
	}

	slog.Debug("testing connection...")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("error while testing connection!", slog.Any("error", err))
		return false
	}
	defer res.Body.Close()

	slog.Debug("response get!")

	return res.StatusCode == 204
}

func main() {
	if !force {
		if TestConnection() {
			slog.Info("connection is healthy!")
			os.Exit(0)
		}
	}

	for _, e := range entries {
		slog.Debug("changing wlan connection...", slog.Any("ssid", e.SSID))

		if _, err := ChangeWlan(e.SSID, e.Password); err != nil {
			slog.Debug("error while trying to change wlan", slog.Any("error", err))
			continue
		}

		slog.Debug("changed wlan", slog.String("ssid", e.SSID))

		if !TestConnection() {
			slog.Debug("connection is unhealthy...")
			continue
		}

		slog.Info("connection is healthy!")
		os.Exit(0)
	}

	slog.Info("could not find a good connection...")
	os.Exit(1)
}
