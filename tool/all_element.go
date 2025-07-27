package tool

import (
	"context"
	"time"

	mcpcdp "github.com/patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewAllElementTool() mcp.Tool {
	return mcp.NewTool("get-all-elements",
		mcp.WithDescription("Get all elements of current page, return a clean DOM element tree structure without scripts/styles and textContent, if content is truncated, call the select-element tool to get a deeper DOM tree"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance id"),
			),
		mcp.WithNumber("depth", 
			 mcp.Description("Maximum DOM tree depth to traverse (default: 5)"),
			),
		)
}

func AllElementHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error){
	id := request.GetString("id", "")

	if id == "" {
		return mcp.NewToolResultError("Chrome instance ID is required"), nil
	}

	depth := max(request.GetInt("depth", 5), 0)

	var cleanHTML string
	
	err := mcpcdp.Manager.Execute(id,
		chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(cleanElement(depth), &cleanHTML),
		// chromedp.Evaluate(`
		// 	(() => {
		// 		function cleanElement(element) {
		// 			const newEl = document.createElement(element.tagName);
		// 			Array.from(element.attributes).forEach(attr => {
		// 				newEl.setAttribute(attr.name, attr.value);
		// 			});
		// 			Array.from(element.children).forEach(child => {
		// 				if (!['SCRIPT', 'STYLE', 'NOSCRIPT'].includes(child.tagName)) {
		// 					const cleanChild = cleanElement(child);
		// 					newEl.appendChild(cleanChild);
		// 				}
		// 			});
		// 			return newEl;
		// 		}
		//
		// 		const cleanBody = cleanElement(document.body);
		// 		return cleanBody.outerHTML;
		// 	})()
		// `, &cleanHTML),
	)

	if (err != nil) {
		return nil, err
	}

	return mcp.NewToolResultText(cleanHTML),err
}
