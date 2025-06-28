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

	maximumInstance := os.Getenv("MAXIMUM_INSTANCE")
	ttlStr := os.Getenv("TTL")
	if maximumInstance == "" {
		maximumInstance = "5"
	}

	if ttlStr == "" {
		ttlStr = "15"
	}

	maximum,err := strconv.Atoi(maximumInstance)
	if err != nil {
		maximum = 5
	}

	ttl,err := strconv.Atoi(ttlStr)
	if err != nil {
		ttl = 15
	}

	chromedp.InitManager(maximum, time.Duration(ttl)*time.Minute)

	// Add tool
	pdfTool := tool.NewPdfTool()
	s.AddTool(pdfTool, tool.GenPdfHandler)
	s.AddTool(tool.NewCreateInstanceTool(), tool.CreateInstanceHandler)
	s.AddTool(tool.NewCloseInstanceTool(), tool.CloseInstanceHandler)
	s.AddTool(tool.NewNavigateTool(), tool.NavigateHandler)
	s.AddTool(tool.NewGetElementTool(), tool.GetElementHandler)
	s.AddTool(tool.NewClickElementTool(), tool.ClickElementHandler)
	s.AddTool(tool.NewAllElementTool(), tool.AllElementHandler)
	// Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }

}




