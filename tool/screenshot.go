package tool

import (
	"context"
	"encoding/base64"
	"fmt"
	mcpcdp "github.com/patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewScreenshotTool() mcp.Tool {
	return mcp.NewTool("screenshot",
		mcp.WithDescription("Take a screenshot of a specific element or the entire browser viewport"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the screenshot action on"),
		),
		mcp.WithString("url", 
			mcp.Description("Navigate to the provided URL before taking screenshot. If URL is not provided, will stay at current page"),	
		),	
		mcp.WithString("selector", 
			mcp.Description("The selector to identify the element for screenshot. Examples: '#input-id', '.input-class', 'input[name=\"username\"]', '//input[@id=\"email\"]'. If not provided, will take screenshot of entire browser viewport"),
		),
	)
}

func ScreenshotHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	url := request.GetString("url", "")
	selector := request.GetString("selector", "")
	
	if id == "" {
		return mcp.NewToolResultError("Chrome instance ID is required"), nil
	}

	// Navigate to URL if provided
	if url != "" {
		err := mcpcdp.Manager.Execute(id, 
			chromedp.Navigate(url),
			chromedp.WaitReady("body"),
		)
		if err != nil {
			return nil, fmt.Errorf("error happened when navigating to URL: %v", err)
		}
	}

	var screenshotBuf []byte
	var screenshotDescription string

	// Take screenshot based on whether selector is provided
	if selector != "" {
		// Check if element exists before taking screenshot
		var elements []*cdp.Node
		err := mcpcdp.Manager.Execute(id, 
			chromedp.Nodes(selector, &elements),
		)
		if err != nil || len(elements) < 1 {
			return mcp.NewToolResultError(fmt.Sprintf("No elements found with selector: %s, please check element tree again", selector)), nil
		}

		// Take screenshot of specific element
		err = mcpcdp.Manager.Execute(id,
			chromedp.Screenshot(selector, &screenshotBuf),
		)
		if err != nil {
			return nil, fmt.Errorf("error happened when taking element screenshot: %v", err)
		}
		screenshotDescription = fmt.Sprintf("Screenshot of element with selector: %s", selector)
	} else {
		// Take screenshot of entire viewport
		err := mcpcdp.Manager.Execute(id,
			chromedp.FullScreenshot(&screenshotBuf, 100),
		)
		if err != nil {
			return nil, fmt.Errorf("error happened when taking full screenshot: %v", err)
		}
		screenshotDescription = "Screenshot of entire browser viewport"
	}

	// Convert screenshot to base64
	imageData := base64.StdEncoding.EncodeToString(screenshotBuf)
	
	// Return image result
	return mcp.NewToolResultImage(screenshotDescription, imageData, "image/png"), nil
}
