package tool

import (
	"context"
	"errors"

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
		return nil,errors.New("url should be provide")
	}

	id := request.GetString("id", "")

	if id == "" {
		return nil, errors.New("id should be provide")
	}


	var cleanHTML string
	err := cdp.Manager.Execute(id,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		 chromedp.Evaluate(`
		          function cleanElement(element) {
		              // 建立新的元素副本
		              const newEl = document.createElement(element.tagName);

		              // 複製屬性
		              Array.from(element.attributes).forEach(attr => {
		                  newEl.setAttribute(attr.name, attr.value);
		              });

		              // 遞迴處理子元素
		              Array.from(element.children).forEach(child => {
		                  // 跳過 script, style 等
		                  if (!['SCRIPT', 'STYLE', 'NOSCRIPT'].includes(child.tagName)) {
		                      const cleanChild = cleanElement(child);
		                      newEl.appendChild(cleanChild);
		                  }
		              });

		              return newEl;
		          }

		          // 清理整個 body
		          const cleanBody = cleanElement(document.body);
		          cleanBody.outerHTML;
		      `, &cleanHTML),
	)

	if (err != nil) {
		return nil, err
	}

	return mcp.NewToolResultText(cleanHTML),err
}
