package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
	mcpcdp "patrickke/chromedp-mcp/chromedp"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewDownloadFileTool() mcp.Tool {
	return mcp.NewTool("download-file",
		mcp.WithDescription("The download-file tool in chromedp downloads files by clicking on download links or buttons. It can optionally specify a download directory."),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the download action on"),
		),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("The selector to identify the download element to click. Examples: '#download-link', '.download-btn', 'a[href*=\"download\"]', '//a[contains(@href, \"download\")]'"),
		),
		mcp.WithString("download_path",
			mcp.Description("Optional download directory path. If not specified, uses user's default Downloads directory"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Download timeout in seconds (default: 30)"),
		),
	)
}

// getDefaultDownloadPath returns the user's default Downloads directory
func getDefaultDownloadPath() string {
	var downloadPath string
	switch runtime.GOOS {
	case "windows":
		downloadPath = filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	case "darwin":
		downloadPath = filepath.Join(os.Getenv("HOME"), "Downloads")
	default: // linux and other unix-like systems
		downloadPath = filepath.Join(os.Getenv("HOME"), "Downloads")
	}
	return downloadPath
}

func DownloadFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := request.GetString("selector", "")
	id := request.GetString("id", "")
	downloadPath := request.GetString("download_path", getDefaultDownloadPath())
	timeout := request.GetInt("timeout", 30)

	if selector == "" {
		return mcp.NewToolResultError("selector is required"), nil
	}
	if id == "" {
		return mcp.NewToolResultError("Chrome instance id is required"), nil
	}

	var elements []*cdp.Node
	err := mcpcdp.Manager.Execute(id, 
		chromedp.Nodes(selector, &elements),
		)
	if err != nil || len(elements) < 1{
		return mcp.NewToolResultError(fmt.Sprintf("No elements found with selector: %s, please check element tree again or use screenshot for anlalyze", selector)), nil
	}

	// Set up download directory
	var downloadDir string
	// Ensure download directory exists
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create download directory: %v", err)), nil
	}
	downloadDir, _ = filepath.Abs(downloadPath)

	var downloadGUID string

	// Create execution tasks
	tasks := []chromedp.Action{
		chromedp.WaitVisible(selector),
		chromedp.Nodes(selector, &elements),
	}

	// Set download behavior with specified directory
	tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
		return browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath(downloadDir).
			WithEventsEnabled(true).
			Do(ctx)
	}))

	// Listen for download events
	downloadStarted := make(chan string, 1)
	downloadComplete := make(chan bool, 1)

	tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			switch ev := ev.(type) {
			case *browser.EventDownloadWillBegin:
				downloadGUID = ev.GUID
				downloadStarted <- ev.GUID
			case *browser.EventDownloadProgress:
				if ev.GUID == downloadGUID && ev.State == browser.DownloadProgressStateCompleted {
					downloadComplete <- true
				} else if ev.GUID == downloadGUID && ev.State == browser.DownloadProgressStateCanceled {
					downloadComplete <- false
				}
			}
		})
		return nil
	}))

	// Click download element
	tasks = append(tasks, chromedp.Click(selector))

	err = mcpcdp.Manager.Execute(id, tasks...)
	if err != nil {
		return nil, fmt.Errorf("error happen when execute download, err: %v", err)
	}

	if len(elements) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No elements found with selector: %s", selector)), nil
	}

	// Wait for download to start
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	select {
	case guid := <-downloadStarted:
		// Wait for download to complete
		select {
		case success := <-downloadComplete:
			if success {
				result := fmt.Sprintf("Download completed successfully. GUID: %s, Download directory: %s", guid, downloadDir)
				return mcp.NewToolResultText(result), nil
			} else {
				return mcp.NewToolResultError("Download was canceled"), nil
			}
		case <-timeoutCtx.Done():
			return mcp.NewToolResultError("Download timeout"), nil
		}
	case <-timeoutCtx.Done():
		return mcp.NewToolResultError("Download did not start within timeout"), nil
	}
}
