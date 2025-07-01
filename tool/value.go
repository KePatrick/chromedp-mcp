package tool

import (
	"context"
	"fmt"

	mcpcdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewSetValueTool() mcp.Tool {
	return mcp.NewTool("set-value",
		mcp.WithDescription("The set-value tool in chromedp directly sets the value of a specified element by selector on a web page. This is more efficient than send-key for setting form field values."),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the set value action on"),
		),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("The selector to identify the element to set value. Examples: '#input-id', '.input-class', 'input[name=\"username\"]', '//input[@id=\"email\"]'"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The value to set for the element"),
		),
	)
}

func SetValueHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := request.GetString("selector", "")
	id := request.GetString("id", "")
	value := request.GetString("value", "")

	if selector == "" {
		return mcp.NewToolResultError("selector is required"), nil
	}
	if id == "" {
		return mcp.NewToolResultError("Chrome instance id is required"), nil
	}
	if value == "" {
		return mcp.NewToolResultError("value is required"), nil
	}
	var elements []*cdp.Node
	err := mcpcdp.Manager.Execute(id, 
		chromedp.Nodes(selector, &elements),
		)
	if err != nil || len(elements) < 1{
		return mcp.NewToolResultError(fmt.Sprintf("No elements found with selector: %s, please check element tree again or use screenshot for anlalyze", selector)), nil
	}

	err = mcpcdp.Manager.Execute(id,
		chromedp.SetValue(selector, value),
		chromedp.Nodes(selector, &elements),
	)
	if err != nil {
		return nil, fmt.Errorf("error happen when execute set value, err: %v", err)
	}

	if len(elements) == 0 {
		return mcp.NewToolResultError(fmt.Sprintf("No elements found with selector: %s, please check element tree again or use screenshot for anlalyze", selector)), nil
	}

	htmlResults, err := convertNodesToHTML(id, elements)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to convert nodes to HTML: %v", err)), nil
	}

	result := formatElementResults(selector, htmlResults)
	return mcp.NewToolResultText(fmt.Sprintf("Successfully set value '%s' to element. Element after set value: %s", value, result)), nil
}
