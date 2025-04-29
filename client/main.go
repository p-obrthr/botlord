package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	width              int32 = 1000
	height             int32 = 600
	checkInterval            = 30 * time.Second
	lastChecked        time.Time
	isRunning          bool = false
	statusMutex        sync.Mutex
	statusText         string
	statusColor        rl.Color = rl.Yellow
	isActionInProgress bool     = false
	botIp              string
)

func main() {

	botIpEnv, exists := os.LookupEnv("BOT_IP")
	if !exists {
		fmt.Printf("err: no bot ip")
	}
	botIp = botIpEnv

	rl.InitWindow(width, height, "BOTLORD")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	xButtonRec := float32((width - 200) / 2)
	yButtonRec := float32(150)
	buttonRec := rl.NewRectangle(xButtonRec, yButtonRec, 200, 50)

	go updateBotStatus()

	for !rl.WindowShouldClose() {
		if time.Since(lastChecked) > checkInterval && !isActionInProgress {
			go updateBotStatus()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Color{R: 30, G: 30, B: 30, A: 255})

		mousePoint := rl.GetMousePosition()
		color := rl.DarkGray

		if rl.CheckCollisionPointRec(mousePoint, buttonRec) {
			color = rl.Gray
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				color = rl.LightGray
			}
			if rl.IsMouseButtonReleased(rl.MouseLeftButton) && !isActionInProgress {
				go changeBotStatus()
			}
		}

		rl.DrawRectangleRec(buttonRec, color)
		text := "Change Status"
		fontSize := int32(20)
		textWidth := rl.MeasureText(text, fontSize)
		x := int32(buttonRec.X + (buttonRec.Width-float32(textWidth))/2)
		y := int32(buttonRec.Y + (buttonRec.Height-float32(fontSize))/2)
		rl.DrawText(text, x, y, fontSize, rl.Lime)

		statusMutex.Lock()
		localStatusText := statusText
		localStatusColor := statusColor
		statusMutex.Unlock()

		textWidthCounter := rl.MeasureText(localStatusText, fontSize)
		xButtonClicked := (width - int32(textWidthCounter)) / 2
		yButtonClicked := int32(350)
		rl.DrawText(localStatusText, xButtonClicked, yButtonClicked, fontSize, localStatusColor)

		rl.EndDrawing()
	}
}

func updateBotStatus() {
	statusMutex.Lock()
	statusText = "CHECKING..."
	statusColor = rl.Yellow
	statusMutex.Unlock()

	lastChecked = time.Now()

	running := checkBotRun()

	statusMutex.Lock()
	isRunning = running
	if running {
		statusText = "BOT IS ONLINE"
		statusColor = rl.Green
	} else {
		statusText = "BOT IS OFFLINE"
		statusColor = rl.Red
	}
	statusMutex.Unlock()
}

func changeBotStatus() {
	isActionInProgress = true

	statusMutex.Lock()
	currentlyRunning := isRunning
	statusMutex.Unlock()

	if currentlyRunning {
		statusMutex.Lock()
		statusText = "STOPPING..."
		statusColor = rl.Yellow
		statusMutex.Unlock()
		stopBot()
	} else {
		statusMutex.Lock()
		statusText = "STARTING..."
		statusColor = rl.Yellow
		statusMutex.Unlock()
		startBot()
	}

	time.Sleep(500 * time.Millisecond)
	updateBotStatus()
	isActionInProgress = false
}

func checkBotRun() bool {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/status", botIp))
	if err != nil {
		fmt.Println("error checking status:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading status response:", err)
		return false
	}

	return string(body) == "running"
}

func startBot() {
	resp, err := http.Post(fmt.Sprintf("http://%s:8080/start", botIp), "application/json", nil)
	if err != nil {
		fmt.Println("err starting container:", err)
		return
	}
	defer resp.Body.Close()
}

func stopBot() {
	resp, err := http.Post(fmt.Sprintf("http://%s:8080/stop", botIp), "application/json", nil)
	if err != nil {
		fmt.Println("err stop container:", err)
		return
	}
	defer resp.Body.Close()
}
