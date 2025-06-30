package tool

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	cdp "patrickke/chromedp-mcp/chromedp"
	"runtime"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewDownloadImageTool() mcp.Tool {
    return mcp.NewTool("download_image",
        mcp.WithDescription("Download image from URL or selector"),
        mcp.WithString("id", mcp.Required(), mcp.Description("Chrome instance ID")),
        mcp.WithString("url", mcp.Description("Image URL to download")),
        mcp.WithString("selector", 
            mcp.Description("The selector to identify the element to click. Examples: '#button-id', '.button-class', 'button[type=\"submit\"]', '//button[@id=\"submit\"]'"),
        ),
        mcp.WithString("output_path", mcp.Description("Output directory path (default: user downloads directory)")),
    )
}

func DownloadImageHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    id := request.GetString("id", "")
    if id == "" {
        return mcp.NewToolResultError("Chrome instance ID required"), nil
    }
    
    imageURL := request.GetString("url", "")
    selector := request.GetString("selector", "")
    outputPath := request.GetString("output_path", "")
    
    if imageURL == "" && selector == "" {
        return mcp.NewToolResultError("Either url or selector required"), nil
    }
    
    var imageData string
    var err error
    
    if imageURL != "" {
        err = downloadFromURL(id, imageURL, &imageData)
    } else {
        err = downloadFromSelector(id, selector, &imageData)
    }
    
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Download failed: %v", err)), nil
    }
    
    // Set output directory
    outputDir := outputPath
    if outputDir == "" {
        var err error
        outputDir, err = getDownloadDirectory()
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("Failed to get download directory: %v", err)), nil
        }
    }
    
    // Generate filename
    filename := fmt.Sprintf("image_%d.png", time.Now().Unix())
    
    // Combine directory and filename
    fullPath := filepath.Join(outputDir, filename)
    
    // Ensure absolute path
    absPath, err := filepath.Abs(fullPath)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Failed to get absolute path: %v", err)), nil
    }
    
    size, err := saveImage(imageData, absPath)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Save failed: %v", err)), nil
    }
    
    result := fmt.Sprintf("Downloaded successfully!\nFilename: %s\nLocation: %s\nSize: %d bytes", 
        filename, absPath, size)
    return mcp.NewToolResultText(result), nil
}

func getDownloadDirectory() (string, error) {
    var downloadDir string
    
    switch runtime.GOOS {
    case "windows":
        userProfile := os.Getenv("USERPROFILE")
        if userProfile == "" {
            return "", fmt.Errorf("USERPROFILE environment variable not set")
        }
        downloadDir = filepath.Join(userProfile, "Downloads")
    case "darwin": // macOS
        homeDir := os.Getenv("HOME")
        if homeDir == "" {
            return "", fmt.Errorf("HOME environment variable not set")
        }
        downloadDir = filepath.Join(homeDir, "Downloads")
    case "linux":
        homeDir := os.Getenv("HOME")
        if homeDir == "" {
            return "", fmt.Errorf("HOME environment variable not set")
        }
        // Try XDG user directory first
        xdgDownloadDir := os.Getenv("XDG_DOWNLOAD_DIR")
        if xdgDownloadDir != "" {
            downloadDir = xdgDownloadDir
        } else {
            downloadDir = filepath.Join(homeDir, "Downloads")
        }
    default:
        return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
    }
    
    return downloadDir, nil
}

func downloadFromURL(instanceID, url string, imageData *string) error {
    return cdp.Manager.Execute(instanceID,
        chromedp.Navigate(url),
        chromedp.WaitVisible("img"),
        chromedp.Evaluate(captureImageScript(), imageData),
		chromedp.NavigateBack(),
    )
}

func downloadFromSelector(instanceID, selector string, imageData *string) error {
    return cdp.Manager.Execute(instanceID,
        chromedp.WaitVisible(selector),
        chromedp.Evaluate(captureElementScript(selector), imageData),
    )
}

func captureImageScript() string {
    return `
        (() => {
            const img = document.querySelector('img');
            if (!img) return '';
            
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            
            canvas.width = img.naturalWidth || img.width;
            canvas.height = img.naturalHeight || img.height;
            
            ctx.drawImage(img, 0, 0);
            
            return canvas.toDataURL('image/png').split(',')[1];
        })()
    `
}

func captureElementScript(selector string) string {
    return fmt.Sprintf(`
        (() => {
            const element = document.querySelector('%s');
            if (!element) return '';
            
            const img = element.tagName === 'IMG' ? element : element.querySelector('img');
            if (!img) return '';
            
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            
            canvas.width = img.naturalWidth || img.width;
            canvas.height = img.naturalHeight || img.height;
            
            ctx.drawImage(img, 0, 0);
            
            return canvas.toDataURL('image/png').split(',')[1];
        })()
    `, selector)
}

func saveImage(base64Data, filename string) (int64, error) {
    if base64Data == "" {
        return 0, fmt.Errorf("no image data")
    }
    
    data, err := base64.StdEncoding.DecodeString(base64Data)
    if err != nil {
        return 0, err
    }
    
    dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return 0, err
    }
    
    err = os.WriteFile(filename, data, 0644)
    return int64(len(data)), err
}
