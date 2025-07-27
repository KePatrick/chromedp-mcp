package tool

import (
	"context"
	"fmt"
	"strings"
	"time"

	mcpcdp "github.com/patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewSetCookieTool() mcp.Tool {
	return mcp.NewTool("set-cookie",
		mcp.WithDescription("Set a HTTP cookie on requests. Cookies will be automatically sent with subsequent requests to matching domains and paths."),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the set cookie action on"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the cookie to set"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The value of the cookie to set"),
		),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain to set the cookie for (e.g., 'example.com', '.example.com' for subdomains)"),
		),
		mcp.WithString("path", 
			mcp.Required(),
			mcp.Description("The path to set the cookie for. Use '/' for all paths under the domain"),
		),
		mcp.WithBoolean("httpOnly",
			mcp.Description("If true, the cookie is only accessible via HTTP requests and not JavaScript. Recommended for security-sensitive cookies like session tokens. Default: false"),
		),
		mcp.WithBoolean("secure",
			mcp.Description("If true, the cookie is only sent over HTTPS connections. Recommended for production environments. Default: false"),
		),
		mcp.WithString("sameSite",
			mcp.Description("Controls when the cookie is sent with cross-site requests. Options: 'Strict' (only same-site), 'Lax' (some cross-site), 'None' (all cross-site, requires Secure=true). Default: 'Lax'"),
		),
		mcp.WithNumber("expires",
			mcp.Description("Cookie expiration time as Unix timestamp (seconds since epoch). If not set, cookie expires when browser session ends. Example: 1735689600 for 2025-01-01"),
		),
		mcp.WithNumber("maxAge",
			mcp.Description("Cookie lifetime in seconds from when it's set. Takes precedence over 'expires' if both are set. Example: 3600 for 1 hour, 86400 for 1 day"),
		),
	)
}

func SetCookieHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Required parameters
	id := request.GetString("id", "")
	if id == "" {
		return mcp.NewToolResultError("Chrome instance ID is required"), nil
	}
	
	name := request.GetString("name", "")
	if name == "" {
		return mcp.NewToolResultError("Cookie name is required"), nil
	}
	
	value := request.GetString("value", "")
	if value == "" {
		return mcp.NewToolResultError("Cookie value is required"), nil
	}
	
	domain := validateDomain(request.GetString("domain", ""))
	if domain == "" {
		return mcp.NewToolResultError("Cookie domain is required"), nil
	}
	
	path := request.GetString("path", "/")
	if path == "" {
		path = "/" // Ensure at least root path
	}
	
	// Optional parameters - Security settings
	httpOnly := request.GetBool("httpOnly", false)
	secure := request.GetBool("secure", false)
	
	// SameSite setting
	sameSite := request.GetString("sameSite", "Lax")
	// Validate sameSite value
	validSameSiteValues := map[string]bool{
		"Strict": true,
		"Lax":    true,
		"None":   true,
	}
	if !validSameSiteValues[sameSite] {
		sameSite = "Lax" // Default value
	}
	
	// Expiration time settings
	expires := request.GetInt("expires", 0)
	maxAge := request.GetInt("maxAge", 0)
	
	// Validate SameSite=None requires Secure=true
	if sameSite == "None" && !secure {
		return mcp.NewToolResultError("SameSite=None requires Secure=true"), nil
	}

	instance, err := mcpcdp.Manager.GetInstance(id)

	if err != nil {
		return nil, fmt.Errorf("fail to get instance: %v", err)
	}

	chromeCtx := instance.Context

	err = chromedp.Run(chromeCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Build basic SetCookie command
			cookieCmd := network.SetCookie(name, value).
				WithDomain(domain).
				WithPath(path).
				WithHTTPOnly(httpOnly).
				WithSecure(secure)
			
			// Set SameSite attribute
			switch sameSite {
			case "Strict":
				cookieCmd = cookieCmd.WithSameSite(network.CookieSameSiteStrict)
			case "Lax":
				cookieCmd = cookieCmd.WithSameSite(network.CookieSameSiteLax)
			case "None":
				cookieCmd = cookieCmd.WithSameSite(network.CookieSameSiteNone)
			}
			
			// Set expiration time (maxAge takes precedence over expires)
			if maxAge > 0 {
				// maxAge is relative time, convert to absolute time
				expirationTime := time.Now().Add(time.Duration(maxAge) * time.Second)
				expireTimestamp := cdp.TimeSinceEpoch(expirationTime)
				cookieCmd = cookieCmd.WithExpires(&expireTimestamp)
			} else if expires > 0 {
				expireTime := time.Unix(int64(expires), 0)
				expireTimestamp := cdp.TimeSinceEpoch(expireTime)
				cookieCmd = cookieCmd.WithExpires(&expireTimestamp)
			}
			
			// Execute the setting command
			return cookieCmd.Do(ctx)
		}),
	)
	
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set cookie: %v", err)), nil
	}
	
	// Build success response
	resultData := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Cookie '%s' set successfully", name),
		"cookie": map[string]interface{}{
			"name":     name,
			"value":    value,
			"domain":   domain,
			"path":     path,
			"httpOnly": httpOnly,
			"secure":   secure,
			"sameSite": sameSite,
		},
	}

	// Include expiration info in response if set
	if maxAge > 0 {
		resultData["cookie"].(map[string]interface{})["maxAge"] = maxAge
		resultData["cookie"].(map[string]interface{})["expiresAt"] = formatExpirationTime(float64(time.Now().Add(time.Duration(maxAge) * time.Second).Unix()))
	} else if expires > 0 {
		resultData["cookie"].(map[string]interface{})["expires"] = expires
		resultData["cookie"].(map[string]interface{})["expiresAt"] = formatExpirationTime(float64(expires))
	} else {
		resultData["cookie"].(map[string]interface{})["expiresAt"] = "Session (expires when browser closes)"
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("Cookie set successfully: %+v", resultData)), nil

}

// Helper function: Validate and clean domain
func validateDomain(domain string) string {
	// Remove protocol prefix
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	
	// Remove path
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}
	
	// Remove port number
	if idx := strings.Index(domain, ":"); idx != -1 {
		domain = domain[:idx]
	}
	
	return strings.ToLower(domain)
}

// Helper function: Format expiration time to human-readable format
func formatExpirationTime(expires float64) string {
	if expires <= 0 {
		return "Session (expires when browser closes)"
	}
	
	expirationTime := time.Unix(int64(expires), 0)
	return expirationTime.Format("2006-01-02 15:04:05 UTC")
}


