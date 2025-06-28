package tool

import (
	"context"
	"fmt"
	"strings"

	mcpcdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetElementTool() mcp.Tool {
	return mcp.NewTool("get-element",
		mcp.WithDescription("Get element specified by selector (CSS selector, XPath, ID, class, etc.).  You should create a instance before operation"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("The selector to identify the element to click. Examples: '#button-id', '.button-class', 'button[type=\"submit\"]', '//button[@id=\"submit\"]'"),
			),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the click action on"),
			),
		)
}

func GetElementHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    selector := request.GetString("selector", "")
    id := request.GetString("id", "")
    
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
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to get elements: %v", err)), nil
    }
    
    if len(elements) == 0 {
        return mcp.NewToolResultText(fmt.Sprintf("No elements found with selector: %s", selector)), nil
    }
    
    htmlResults, err := convertNodesToHTML(id, elements)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to convert nodes to HTML: %v", err)), nil
    }
    
    result := formatElementResults(selector, htmlResults)
    return mcp.NewToolResultText(result), nil
}


func convertNodesToHTML(instanceID string, nodes []*cdp.Node) ([]string, error) {
    var htmlResults []string
    
    for i, node := range nodes {
        var html string
        err := mcpcdp.Manager.Execute(instanceID,
            chromedp.OuterHTML([]cdp.NodeID{node.NodeID}, &html, chromedp.ByNodeID),
        )
        if err != nil {
            fmt.Printf("Warning: Failed to get HTML for node %d: %v\n", i, err)
            continue
        }
        htmlResults = append(htmlResults, html)
    }
    
    return htmlResults, nil
}

func formatElementResults(selector string, htmlResults []string) string {
    var result strings.Builder
    
    result.WriteString(fmt.Sprintf("Found %d element(s) with selector '%s':\n\n", len(htmlResults), selector))
    
    for i, html := range htmlResults {
        result.WriteString(fmt.Sprintf("=== Element %d ===\n", i+1))
        result.WriteString(html)
        result.WriteString("\n\n")
    }
    
    return result.String()
}
