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
	return mcp.NewTool("create-chrome-instance",
		mcp.WithDescription("Create Chrome Instance, every session should start by create_chrome_instance and end by end_chrome_instance"),
		mcp.WithBoolean("headless",
			mcp.Description("Headless mode flag for create chrome instance (default: true)"),
			),
		mcp.WithBoolean("disable-gpu",
			mcp.Description("Disable gpu for chrome instance (default: true)"),
			),
		mcp.WithBoolean("disable-popup-blocking",
			mcp.Description("Disable popup blocking to allow popups (default: false, meaning popups are blocked)"),
			),
		mcp.WithBoolean("block-new-tab",
			mcp.Description("Block opening new tabs/windows and redirect to current tab (default: false)"),
			),
		mcp.WithBoolean("disable-extensions",
			mcp.Description("Disable browser extensions (default: true)"),
			),
		mcp.WithBoolean("disable-plugins",
			mcp.Description("Disable browser plugins (default: true)"),
			),
		mcp.WithBoolean("disable-web-security",
			mcp.Description("Disable web security (CORS) for testing purposes (default: false)"),
			),
		mcp.WithBoolean("no-sandbox",
			mcp.Description("Disable Chrome sandbox for running in containers (default: false)"),
			),
		mcp.WithBoolean("disable-dev-shm-usage",
			mcp.Description("Disable /dev/shm usage to avoid memory issues in containers (default: false)"),
			),
		mcp.WithBoolean("disable-background-timer-throttling",
			mcp.Description("Disable background timer throttling for better performance (default: false)"),
			),
		mcp.WithBoolean("disable-backgrounding-occluded-windows",
			mcp.Description("Disable backgrounding occluded windows (default: false)"),
			),
		mcp.WithBoolean("disable-renderer-backgrounding",
			mcp.Description("Disable renderer backgrounding (default: false)"),
			),
		)
}

func CreateInstanceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	headless := request.GetBool("headless", true)
	disableGpu := request.GetBool("disable-gpu", true)
	disablePopupBlocking := request.GetBool("disable-popup-blocking", false)
	// blockNewTab := request.GetBool("block-new-tab", false)
	disableExtensions := request.GetBool("disable-extensions", true)
	disablePlugins := request.GetBool("disable-plugins", true)
	disableWebSecurity := request.GetBool("disable-web-security", false)
	noSandbox := request.GetBool("no-sandbox", false)
	disableDevShmUsage := request.GetBool("disable-dev-shm-usage", false)
	disableBackgroundTimerThrottling := request.GetBool("disable-background-timer-throttling", false)
	disableBackgroundingOccludedWindows := request.GetBool("disable-backgrounding-occluded-windows", false)
	disableRendererBackgrounding := request.GetBool("disable-renderer-backgrounding", false)

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.Flag("headless", headless),                                                    // Enable/disable headless mode
	chromedp.Flag("disable-gpu", disableGpu),                                              // Enable/disable GPU
	chromedp.Flag("disable-popup-blocking", disablePopupBlocking),                        // Control popup blocking
	chromedp.Flag("disable-extensions", disableExtensions),                               // Disable browser extensions
	chromedp.Flag("disable-plugins", disablePlugins),                                     // Disable browser plugins
	chromedp.Flag("disable-web-security", disableWebSecurity),                           // Disable web security (CORS)
	chromedp.Flag("no-sandbox", noSandbox),                                               // Disable sandbox for containers
	chromedp.Flag("disable-dev-shm-usage", disableDevShmUsage),                          // Disable /dev/shm usage
	chromedp.Flag("disable-background-timer-throttling", disableBackgroundTimerThrottling), // Disable background timer throttling
	chromedp.Flag("disable-backgrounding-occluded-windows", disableBackgroundingOccludedWindows), // Disable backgrounding occluded windows
	chromedp.Flag("disable-renderer-backgrounding", disableRendererBackgrounding),        // Disable renderer backgrounding
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
