package main

import (
	"encoding/json"
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
	checkInterval            = 10 * time.Second
	lastChecked        time.Time
	isRunning          bool = false
	statusMutex        sync.Mutex
	logsMutex          sync.Mutex
	logFont            rl.Font
	status             = Status{}
	isActionInProgress bool
	botIp              string
	logs               []string
	backgroundTexture  rl.Texture2D
)

type Status struct {
	text  string
	color rl.Color
}

func main() {
	botIpEnv, exists := os.LookupEnv("BOT_IP")
	if !exists {
		fmt.Printf("err: no bot ip")
	}
	botIp = botIpEnv

	rl.InitWindow(width, height, "BOTLORD")
	defer rl.CloseWindow()

	logFont = rl.LoadFontEx("fonts/font.ttf", 20, nil, 0)
	defer rl.UnloadFont(logFont)

	rl.SetTargetFPS(60)

	xButtonRec := float32((width - 200) / 2)
	yButtonRec := float32(150)
	buttonRec := rl.NewRectangle(xButtonRec, yButtonRec, 200, 50)

	for !rl.WindowShouldClose() {
		if time.Since(lastChecked) > checkInterval && !isActionInProgress {
			go update()
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Color{R: 30, G: 30, B: 30, A: 255})

		handleButtonInteractions(buttonRec)
		drawStatusText()
		drawLogs()

		rl.EndDrawing()
	}
}

func handleButtonInteractions(buttonRec rl.Rectangle) {
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
	drawCenteredText("Change Status", buttonRec)
}

func drawCenteredText(text string, rec rl.Rectangle) {
	fontSize := int32(20)
	textWidth := rl.MeasureText(text, fontSize)
	x := int32(rec.X + (rec.Width-float32(textWidth))/2)
	y := int32(rec.Y + (rec.Height-float32(fontSize))/2)
	rl.DrawText(text, x, y, fontSize, rl.Lime)
}

func drawStatusText() {
	statusMutex.Lock()
	localStatus := status
	statusMutex.Unlock()

	fontSize := int32(20)
	textWidthCounter := rl.MeasureText(localStatus.text, fontSize)
	x := (width - int32(textWidthCounter)) / 2
	y := int32(350)
	rl.DrawText(localStatus.text, x, y, fontSize, localStatus.color)
}

func drawLogs() {
	logTextSize := int32(16)
	maxLogs := 9

	startIndex := len(logs) - maxLogs
	if startIndex < 0 {
		startIndex = 0
	}
	visibleLogs := logs[startIndex:]

	cardX := int32(20)
	cardY := int32(400)
	cardWidth := int32(960)
	cardHeight := int32(160)

	rl.DrawRectangle(int32(cardX), cardY, cardWidth, cardHeight, rl.NewColor(0, 0, 0, 150))

	padding := float32(10)
	for i, msg := range visibleLogs {
		pos := rl.NewVector2(float32(cardX)+padding, float32(cardY)+padding+float32(i)*float32(logTextSize))
		rl.DrawTextEx(logFont, msg, pos, float32(logTextSize), 1.0, rl.RayWhite)
	}
}

func update() {
	statusMutex.Lock()
	status = Status{
		text:  "CHECKING...",
		color: rl.Yellow,
	}
	isRunning = checkBotRun()
	if isRunning {
		status = Status{
			text:  "BOT IS ONLINE",
			color: rl.Green,
		}
	} else {
		status = Status{
			text:  "BOT IS OFFLINE",
			color: rl.Red,
		}
	}
	statusMutex.Unlock()
	newLogs, err := fetchLogs(botIp)
	if err != nil {
		fmt.Println("err fetching logs:", err)
		return
	}
	logsMutex.Lock()
	logs = newLogs
	logsMutex.Unlock()
	lastChecked = time.Now()
}

func changeBotStatus() {
	isActionInProgress = true
	statusMutex.Lock()
	currentlyRunning := isRunning
	statusMutex.Unlock()

	statusMutex.Lock()
	if currentlyRunning {
		status = Status{
			text:  "STOPPING...",
			color: rl.Yellow,
		}
		stopBot()
	} else {
		status = Status{
			text:  "STARTING...",
			color: rl.Yellow,
		}
		startBot()
	}
	statusMutex.Unlock()
	time.Sleep(500 * time.Millisecond)
	update()
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

func fetchLogs(botIp string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/logs", botIp))
	if err != nil {
		return nil, fmt.Errorf("err fetching logs: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("err reading body: %v", err)
	}

	var fetchedLogs []string
	err = json.Unmarshal(body, &fetchedLogs)
	if err != nil {
		return nil, fmt.Errorf("err parsing json: %v", err)
	}

	return fetchedLogs, nil
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
