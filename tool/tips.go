package tool

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewTipsTool() mcp.Tool {
	return mcp.NewTool("tips",
		mcp.WithDescription("Get important usage tips and best practices for Chrome automation tools, see this before you start"),
	)
}

func TipsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tips := `
Chrome Automation Best Practices:

1. INSTANCE MANAGEMENT:
   - Always call 'create-instance' first before any Chrome operations
   - Remember to call 'close-instance' when done to prevent memory leaks
   - Set headless=false if you want to see the browser window

2. ELEMENT SELECTION:
   - If selectors fail, use 'get-all-element' or 'select-element' to examine the DOM
   - Start with simple selectors before trying complex ones
   - Common selectors: #id, .class, [attribute], tagname

3. PERFORMANCE TIPS:
   - High depth values can cause context limit issues

4. TROUBLESHOOTING:
   - Check if Chrome instance exists before operations
   - Verify page load completion before element operations
   - Use 'get-element' to confirm element exists before clicking

5. NAVIGATION:
   - Use 'navigate' for initial page loads
   - Use 'navigate-back'/'navigate-forward' for browser history navigation
   - Wait for page load after navigation before element operations
	`
	
	return mcp.NewToolResultText(tips), nil
}
