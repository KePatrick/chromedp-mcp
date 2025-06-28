package tool

import (
	"context"
	"errors"
	"fmt"

	cdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewCreateInstanceTool() mcp.Tool {
	return mcp.NewTool("create_chrome_instance",
		mcp.WithDescription("Create Chrome Instance, every session should start by create_chrome_instance and end by end_chrome_instance"),
		mcp.WithBoolean("headless",
			mcp.Description("Headless mode flag for create chrome instance, default: true"),
			),
		mcp.WithBoolean("disable-gpu",
			mcp.Description("Disable gpu for chrome instance, default: true"),
			),
		)
}

func CreateInstanceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	headless := request.GetBool("headless", true)
	disableGpu := request.GetBool("disable-gpu", true)

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),        // Enable visible browser
		chromedp.Flag("disable-gpu", disableGpu),     // Enable GPU
	)
	
	id,_,err := cdp.Manager.CreateVisibleInstance(allocOpts)

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Create new Chrome completed! instance ID: %s", id)), nil
	
}

func NewCloseInstanceTool() mcp.Tool {
	return mcp.NewTool("close",
		mcp.WithDescription("close Chrome Instance, every session should start by create_chrome_instance and end by end_chrome_instance"),
		mcp.WithString("id",
			mcp.Description("The ID of the Chrome instance to close"),
			),
		)
}

func CloseInstanceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	if id == "" {
		return nil, errors.New("id should be provide")
	}

	cdp.Manager.CloseInstance(id)

	return mcp.NewToolResultText("Successfully closed the Chrome instance"), nil
}
