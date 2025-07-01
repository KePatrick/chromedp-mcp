package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewPdfTool() mcp.Tool {
	return mcp.NewTool("generate_pdf",
		mcp.WithDescription("Generate PDF from HTML content or URL"),
		mcp.WithString("html",
			mcp.Description("HTML string to generate PDF from"),
			),
		mcp.WithString("url",
			mcp.Description("URL to generate PDF from"),
			),
		mcp.WithString("outputDir", 
			mcp.Description("Output directory path for the PDF file"),
			),
		)
}


func GenPdfHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	html := request.GetString("html", "")
	url := request.GetString("url", "")

	
	outputDir := request.GetString("outputDir", getDefaultDownloadPath())
	
	// Check if HTML or URL is provided
	if html == "" && url == "" {
		return nil, fmt.Errorf("must provide either html or url parameter")
	}
	
	// Expand ~ path
	if strings.HasPrefix(outputDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot get user home directory: %v", err)
		}
		outputDir = filepath.Join(homeDir, outputDir[2:])
	}
	
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create output directory %s: %v", outputDir, err)
	}
	
	// Create Chrome context
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, 
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)...)
	defer cancel()
	
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 30*time.Second)
	defer cancel()
	
	// Generate filename
	var filename string
	if url != "" {
		filename = generateFilenameFromURL(url)
	} else {
		filename = fmt.Sprintf("html_%d.pdf", time.Now().Unix())
	}
	
	outputPath := filepath.Join(outputDir, filename)
	
	// Generate PDF
	var buf []byte
	var err error
	
	if url != "" {
		// Generate PDF from URL
		err = chromedp.Run(timeoutCtx, printToPDF(url, &buf))
	} else {
		// Generate PDF from HTML content
		err = chromedp.Run(timeoutCtx, printToPDFFromHTML(html, &buf))
	}
	
	if err != nil {
		return nil, fmt.Errorf("PDF generation failed: %v", err)
	}

	if outputPath == "" {
		return nil, fmt.Errorf("outputPath is not available %s", outputPath)
	}
	
	// Write PDF file
	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return nil, fmt.Errorf("cannot write PDF file %s: %v", outputPath, err)
	}
	

	return mcp.NewToolResultText(fmt.Sprintf("PDF generation completed! Output path: %s", outputPath)), nil	
}

// Generate PDF from URL
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				// WithPaperWidth(8.27).   // A4 width (inches)
				// WithPaperHeight(11.7).  // A4 height (inches)
				// WithMarginTop(0.4).
				// WithMarginBottom(0.4).
				// WithMarginLeft(0.4).
				// WithMarginRight(0.4).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// Generate PDF from HTML content
func printToPDFFromHTML(html string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				// WithPaperWidth(8.27).   // A4 width (inches)
				// WithPaperHeight(11.7).  // A4 height (inches)
				// WithMarginTop(0.4).
				// WithMarginBottom(0.4).
				// WithMarginLeft(0.4).
				// WithMarginRight(0.4).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// Generate filename from URL
func generateFilenameFromURL(url string) string {
	// Remove protocol and special characters
	filename := strings.ReplaceAll(url, "https://", "")
	filename = strings.ReplaceAll(filename, "http://", "")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "&", "_")
	filename = strings.ReplaceAll(filename, "=", "_")
	
	// Limit filename length
	if len(filename) > 50 {
		filename = filename[:50]
	}
	
	// Add timestamp and extension
	return fmt.Sprintf("%s_%d.pdf", filename, time.Now().Unix())
}
