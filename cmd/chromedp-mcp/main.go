package main

import (
	"fmt"
	"patrickke/chromedp-mcp/tool"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
    s := server.NewMCPServer(
        "chromedp-mcp",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

	// Add tool
	pdfTool := tool.NewPdfTool()
	s.AddTool(pdfTool, tool.GenPdfHandler)

	// Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }

}




