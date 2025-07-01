package tool

import (
	"context"
	"fmt"

	cdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewNavigateTool() mcp.Tool {
	return mcp.NewTool("navigate",
		mcp.WithDescription("Navigate to provided URL and return a clean DOM element tree structure without scripts/styles and textContent, if content is truncated, call the select-element tool to get a deeper DOM tree"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to navigate"),
			),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance id"),
			),
		mcp.WithNumber("depth", 
			 mcp.Description("Maximum DOM tree depth to traverse (default: 5)"),
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

	depth := max(request.GetInt("depth", 5), 0)



	var cleanHTML string
	err := cdp.Manager.Execute(id,
		chromedp.Navigate(url),
		chromedp.Evaluate(cleanElement(depth), &cleanHTML),
		)

	if (err != nil) {
		return nil, err
	}

	return mcp.NewToolResultText(cleanHTML), nil
}

func NewNavigateBackTool() mcp.Tool {
	return mcp.NewTool("navigate-back",
		mcp.WithDescription("Navigate to previous page"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance id"),
			),
		)
}

func NavigateBackHandler (ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error){
	id := request.GetString("id", "")

	if id == "" {
		return mcp.NewToolResultError("Chrome instance ID is required"), nil
	}
	
	var url string

	err := cdp.Manager.Execute(id,
		chromedp.NavigateBack(),
		chromedp.Location(&url),
		)
	if (err != nil) {
		return nil, err
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("Navigate back success, current url : %s", url)), nil
}

func NewNavigateForwardTool() mcp.Tool {
	return mcp.NewTool("navigate-forward",
		mcp.WithDescription("Navigate to next page"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance id"),
		),
	)
}

func NavigateForwardHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	if id == "" {
		return mcp.NewToolResultError("Chrome instance ID is required"), nil
	}
	
	var url string
	err := cdp.Manager.Execute(id,
		chromedp.NavigateForward(),
		chromedp.Location(&url),
	)
	if err != nil {
		return nil, err
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("Navigate forward success, current url: %s", url)), nil
}

func cleanElement(depth int) string {
    return fmt.Sprintf(`
        (() => {
            const MAX_DEPTH = %d;
            
            // Dangerous protocols to filter
            const DANGEROUS_PROTOCOLS = ['javascript:', 'data:', 'vbscript:', 'file:', 'about:', 'blob:'];
            
            function cleanForLLM(element, depth = 0) {
                if (depth > MAX_DEPTH) return '[Content truncated - max depth reached]';
                
                if (['SCRIPT', 'STYLE', 'NOSCRIPT', 'META', 'LINK', 'IFRAME', 'OBJECT', 'EMBED'].includes(element.tagName)) {
                    return '';
                }
                
                const newEl = document.createElement(element.tagName);
                
                // Process attributes with special handling for href and src
                Array.from(element.attributes).forEach(attr => {
                    const name = attr.name.toLowerCase();
                    const value = attr.value || '';
                    
                    // Block event handlers
                    if (name.startsWith('on')) {
                        // Optionally preserve for LLM analysis
                        newEl.setAttribute('data-event-' + name.substring(2), value);
                        return;
                    }
                    
                    // Block style attribute
                    if (name === 'style') return;
                    
                    // Special handling for href and src
                    if (['href', 'src'].includes(name)) {
                        const lowerValue = value.toLowerCase().trim();
                        
                        // Check for dangerous protocols
                        const isDangerous = DANGEROUS_PROTOCOLS.some(protocol => 
                            lowerValue.startsWith(protocol)
                        );
                        
                        if (!isDangerous) {
                            // Safe to include
                            newEl.setAttribute(attr.name, value);
                        } else {
                            // Preserve for LLM analysis but mark as filtered
                            newEl.setAttribute('data-filtered-' + name, value);
                            newEl.setAttribute('data-filter-reason', 'Dangerous protocol: ' + lowerValue.split(':')[0]);
                        }
                        return;
                    }
                    
                    // Keep all other attributes
                    try {
                        newEl.setAttribute(attr.name, value);
                    } catch (e) {
                        // Skip invalid attributes
                    }
                });
                
                Array.from(element.children).forEach(child => {
                    const cleaned = cleanForLLM(child, depth + 1);
                    if (cleaned) {
                        newEl.innerHTML += cleaned;
                    }
                });
                
                // Preserve text content for leaf nodes
                if (element.children.length === 0) {
                    const text = element.textContent?.trim();
                    if (text) newEl.textContent = text;
                }
                
                return newEl.outerHTML;
            }
            
            return cleanForLLM(document.body, 0);
        })()
    `, depth)
}

