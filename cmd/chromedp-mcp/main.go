package main

import (
	"fmt"
	"os"
	"patrickke/chromedp-mcp/chromedp"
	"patrickke/chromedp-mcp/tool"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
    s := server.NewMCPServer(
        "chromedp-mcp",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

	maximumInstance := os.Getenv("CHROME_MAXIMUM_INSTANCE")
	ttlStr := os.Getenv("CHROME_TTL")
	executeTimeout := os.Getenv("CHROME_EXE_TIMEOUT")
	if maximumInstance == "" {
		maximumInstance = "5"
	}

	

	if ttlStr == "" {
		ttlStr = "15"
	}

	if executeTimeout == "" {
		executeTimeout = "300"
	}

	maximum,err := strconv.Atoi(maximumInstance)
	if err != nil {
		maximum = 5
	}

	ttl,err := strconv.Atoi(ttlStr)
	if err != nil {
		ttl = 15
	}

	timeout, err := strconv.Atoi(executeTimeout) 
	if err != nil {
		timeout = 300
	}

	chromedp.InitManager(maximum, time.Duration(ttl)*time.Minute, time.Duration(timeout)*time.Second)

	// Add tool
	pdfTool := tool.NewPdfTool()
	s.AddTool(pdfTool, tool.GenPdfHandler)
	s.AddTool(tool.NewCreateInstanceTool(), tool.CreateInstanceHandler)
	s.AddTool(tool.NewCloseInstanceTool(), tool.CloseInstanceHandler)
	s.AddTool(tool.NewNavigateTool(), tool.NavigateHandler)
	s.AddTool(tool.NewGetElementTool(), tool.GetElementHandler)
	s.AddTool(tool.NewClickElementTool(), tool.ClickElementHandler)
	s.AddTool(tool.NewAllElementTool(), tool.AllElementHandler)
	s.AddTool(tool.NewSetCookieTool(), tool.SetCookieHandler)
	s.AddTool(tool.NewSendKeyTool(), tool.SendKeyHandler)
	s.AddTool(tool.NewSetValueTool(), tool.SetValueHandler)
	s.AddTool(tool.NewKeyEventTool(), tool.KeyEventHandler)
	s.AddTool(tool.NewDownloadFileTool(), tool.DownloadFileHandler)
	// Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }

}




