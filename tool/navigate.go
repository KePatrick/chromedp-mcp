package tool

import (
	"context"

	cdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewNavigateTool() mcp.Tool {
	return mcp.NewTool("navigate",
		mcp.WithDescription("Navigate to provide url, you should create a instance before operation"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to navigate"),
			),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance id"),
			),
		)
}

func NavigateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error){
	url := request.GetString("url", "")

	if url == "" {
		return mcp.NewToolResultError("url parameter is required"), nil
	}

	id := request.GetString("id", "")

	if id == "" {
		return mcp.NewToolResultError("Chrome instance ID is required"), nil
	}


	var cleanHTML string
	err := cdp.Manager.Execute(id,
		chromedp.Navigate(url),
		chromedp.Evaluate(`
			(() => {
			const MAX_DEPTH = 10;

			function cleanElement(element, depth = 0) {
			if (depth > MAX_DEPTH) {
			const placeholder = document.createElement('div');
			placeholder.textContent = '[Content truncated - too deep]';
			return placeholder;
			}

			const newEl = document.createElement(element.tagName);

			Array.from(element.attributes).forEach(attr => {
			try {
			newEl.setAttribute(attr.name, attr.value);
			} catch (e) {
			}
			});

			Array.from(element.children).forEach(child => {
			if (!['SCRIPT', 'STYLE', 'NOSCRIPT'].includes(child.tagName)) {
			try {
			const cleanChild = cleanElement(child, depth + 1);
			newEl.appendChild(cleanChild);
			} catch (e) {
			}
			}
			});

			return newEl;
			}

			try {
			const cleanBody = cleanElement(document.body);
			return cleanBody.outerHTML;
			} catch (error) {
			console.error('DOM cleaning failed:', error);
			return "<div>Error: " + error.message + "</div>";
			}
			})()
			`, &cleanHTML),
		)

	if (err != nil) {
		return nil, err
	}

	return mcp.NewToolResultText(cleanHTML),err
}
