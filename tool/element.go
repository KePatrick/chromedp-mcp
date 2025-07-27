package tool

import (
	"context"
	"fmt"
	"strings"

	mcpcdp "github.com/patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewGetElementTool() mcp.Tool {
	return mcp.NewTool("get-element-withtext",
		mcp.WithDescription("Get element specified by selector (CSS selector, XPath, ID, class, etc.) with text.  You should create a instance before operation"),
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

func NewSelectElementTool() mcp.Tool {
    return mcp.NewTool("select-element",
        mcp.WithDescription("Select element by CSS selector and return clean DOM structure at specified depth without text content"),
        mcp.WithString("selector",
            mcp.Required(),
            mcp.Description("The selector to identify the element to click. Examples: '#button-id', '.button-class', 'button[type=\"submit\"]', '//button[@id=\"submit\"]'"),
        ),
        mcp.WithString("id",
            mcp.Required(),
            mcp.Description("Chrome instance ID"),
        ),
        mcp.WithNumber("depth", 
            mcp.Description("Maximum depth to traverse from selected element (default: 3, max: 10)"),
        ),
        mcp.WithBoolean("all",
            mcp.Description("Select all matching elements instead of just the first one (default: false)"),
        ),
    )
}

func SelectElementHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    selector := request.GetString("selector", "")
    if selector == "" {
        return mcp.NewToolResultError("selector parameter is required"), nil
    }
    
    id := request.GetString("id", "")
    if id == "" {
        return mcp.NewToolResultError("Chrome instance ID is required"), nil
    }
    
    depth := max(request.GetInt("depth", 3), 0)
    
    selectAll := request.GetBool("all", false)
    
    var result string
    err := mcpcdp.Manager.Execute(id,
        chromedp.WaitReady("body"),
        chromedp.Evaluate(selectElement(selector, depth, selectAll), &result),
    )
    
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Element selection failed: %v", err)), nil
    }
    
    return mcp.NewToolResultText(result), nil
}

func selectElement(selector string, depth int, selectAll bool) string {
    return fmt.Sprintf(`
        (() => {
            const SELECTOR = %q;
            const MAX_DEPTH = %d;
            const SELECT_ALL = %t;
            
            function cleanElement(element, currentDepth = 0) {
                if (currentDepth > MAX_DEPTH) {
                    const placeholder = document.createElement('div');
                    placeholder.textContent = '[Depth limit reached]';
                    placeholder.setAttribute('data-truncated', 'true');
                    placeholder.setAttribute('data-depth', currentDepth.toString());
                    return placeholder;
                }
                
                if (!element || !element.tagName) {
                    return null;
                }
                
                const newEl = document.createElement(element.tagName);
                
                // Copy attributes but exclude dangerous content
                Array.from(element.attributes || []).forEach(attr => {
                    if (!attr.name.startsWith('on') && 
                        !['javascript:', 'data:', 'vbscript:'].some(proto => 
                            (attr.value || '').toLowerCase().includes(proto))) {
                        try {
                            newEl.setAttribute(attr.name, attr.value || '');
                        } catch (e) {
                            console.warn('Failed to set attribute:', attr.name, e);
                        }
                    }
                });
                
                // Add depth and selector information
                newEl.setAttribute('data-depth', currentDepth.toString());
                if (currentDepth === 0) {
                    newEl.setAttribute('data-selected-by', SELECTOR);
                }
                
                // Recursively process child elements
                Array.from(element.children || []).forEach(child => {
                    if (!['SCRIPT', 'STYLE', 'NOSCRIPT', 'META', 'LINK', 'TITLE'].includes(child.tagName)) {
                        try {
                            const cleanChild = cleanElement(child, currentDepth + 1);
                            if (cleanChild) {
                                newEl.appendChild(cleanChild);
                            }
                        } catch (e) {
                            console.warn('Failed to process child element:', e);
                        }
                    }
                });
                
                return newEl;
            }
            
            function processElements() {
                try {
                    const elements = SELECT_ALL ? 
                        document.querySelectorAll(SELECTOR) : 
                        [document.querySelector(SELECTOR)].filter(el => el !== null);
                    
                    if (elements.length === 0) {
                        return JSON.stringify({
                            error: 'No elements found for selector: ' + SELECTOR,
                            selector: SELECTOR,
                            selectAll: SELECT_ALL,
                            totalFound: 0
                        });
                    }
                    
                    const results = [];
                    
                    elements.forEach((element, index) => {
                        try {
                            const cleanedElement = cleanElement(element, 0);
                            if (cleanedElement) {
                                results.push({
                                    index: index,
                                    tagName: element.tagName,
                                    selector: SELECTOR,
                                    depth: MAX_DEPTH,
                                    html: cleanedElement.outerHTML
                                });
                            }
                        } catch (e) {
                            console.error('Failed to process element at index', index, ':', e);
                            results.push({
                                index: index,
                                error: 'Processing failed: ' + e.message,
                                selector: SELECTOR
                            });
                        }
                    });
                    
                    if (SELECT_ALL) {
                        return JSON.stringify({
                            selector: SELECTOR,
                            totalFound: elements.length,
                            depth: MAX_DEPTH,
                            elements: results
                        }, null, 2);
                    } else {
                        return results[0] ? results[0].html : '<div>No valid element found</div>';
                    }
                    
                } catch (error) {
                    console.error('Selection failed:', error);
                    return JSON.stringify({
                        error: 'Selection failed: ' + error.message,
                        selector: SELECTOR,
                        selectAll: SELECT_ALL
                    });
                }
            }
            
            return processElements();
        })()
    `, selector, depth, selectAll)
}
