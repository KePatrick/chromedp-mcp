package tool

import (
	"context"
	"fmt"
	mcpcdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewSendKeyTool() mcp.Tool {
	return mcp.NewTool("send-key",
		mcp.WithDescription("The send-key tool in chromedp simulates keyboard input by sending keystrokes to a specified element by selector on a web page."),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the click action on"),
			),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("The selector to identify the element to click. Examples: '#button-id', '.button-class', 'button[type=\"submit\"]', '//button[@id=\"submit\"]'"),
			),
		mcp.WithString("key", 
			mcp.Description("The input value"),
			),
		)
}

func SendKeyHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    selector := request.GetString("selector", "")
    id := request.GetString("id", "")
	key := request.GetString("key", "")
    
    if selector == "" {
        return mcp.NewToolResultError("selector is required"), nil
    }
    if id == "" {
        return mcp.NewToolResultError("Chrome instance id is required"), nil
    }
	var elements []*cdp.Node

	err := mcpcdp.Manager.Execute(id,
		chromedp.SendKeys(selector, key),
		chromedp.Nodes(selector, &elements),
        )
	if err != nil {
		return nil, fmt.Errorf("error happen when excute send key, err: %v", err)
	}
	
    if len(elements) == 0 {
        return mcp.NewToolResultText(fmt.Sprintf("No elements found with selector: %s", selector)), nil
    }
    
    htmlResults, err := convertNodesToHTML(id, elements)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to convert nodes to HTML: %v", err)), nil
    }
    
    result := formatElementResults(selector, htmlResults)


	return mcp.NewToolResultText(fmt.Sprintf("the element after send key : %s", result)), nil
}
