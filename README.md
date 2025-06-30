
# chromedp-mcp

A Model Context Protocol (MCP) server that provides browser automation capabilities using ChromeDP. This server enables AI assistants to interact with web pages, manage browser instances, and perform various web automation tasks.

## ⚠️ Development Status

**This project is currently under development. Some features may be incomplete, unstable, or contain bugs. Use with caution and please report any issues you encounter.**


## Features

- **Multi-instance Chrome Management**:  Create and manage multiple Chrome browser instances with extensive configuration options (headless mode, security settings, performance optimization)
- **Complete Page Navigation**: Navigate to URLs, browser history navigation (back/forward)
- **Element Interaction**: Find, click, and interact with web page elements using various selectors
- **Form Operations**: Input text, set values directly
- **Cookie Management**: Set and manage browser cookies for session handling
- **File Operations**: Download files and images from web pages
- **PDF Generation**: Convert web pages or HTML content to PDF documents
- **Keyboard Simulation**: Send keyboard events, key combinations, and modifier keys with comprehensive key support
- **DOM Tree Extraction**: Get clean DOM structures without scripts/styles for analysis

## Available Tools

For detailed tool specifications,and parameters, see [Tools Specification](doc/tools-spec.md)

### Instance Management
- **`create-chrome-instance`** - Create a new Chrome browser instance with extensive configuration options (headless mode, security settings, performance options)
- **`close`** - Close a specific Chrome browser instance

### Page Navigation  
- **`navigate`** - Navigate to a specified URL and return clean DOM structure
- **`navigate-back`** - Navigate to previous page in browser history
- **`navigate-forward`** - Navigate to next page in browser history

### Element Operations
- **`get-element-withtext`** - Find and retrieve information about a specific element with text content
- **`get-all-elements`** - Get all elements of current page as clean DOM tree structure
- **`select-element`** - Select element by CSS selector and return clean DOM structure at specified depth
- **`click-element`** - Click on an element with support for different click types (left, right, double)

### Input Operations
- **`send-key`** - Send keyboard input to specified elements
- **`set-value`** - Directly set the value of form elements (more efficient than send-key)
- **`key-event`** - Send specific keyboard events and combinations with modifier keys

### Cookie Management
- **`set-cookie`** - Set HTTP cookies with full control over domain, path, security, and expiration settings

### File Operations
- **`download-file`** - Download files by clicking download links or buttons
- **`download_image`** - Download images from URL or by selector

### Document Generation
- **`generate_pdf`** - Generate PDF from HTML content or URL

### Other
- **`tips`** - Important usage tips and best practices for Chrome automation tools (recommended to check first)



## Requirements

- **Chrome/Chromium Browser**: Must be installed and accessible in system PATH
- **Go**: Version 1.19 or higher for building from source

## Environment Configuration

The following environment variables can be configured for Chrome management:

| Variable | Description | Default |
|----------|-------------|---------|
| `CHROME_MAXIMUM_INSTANCE` | Maximum concurrent Chrome instances | `5` |
| `CHROME_TTL` | Instance time-to-live in minutes | `15` |
| `CHROME_EXE_TIMEOUT` | Operation timeout in seconds | `300` |

These can be configured through your MCP client configuration or environment variables.

## Setup

### 1. Clone the repository:
   ```bash
   git clone https://github.com/KePatrick/chromedp-mcp.git

   cd chromedp-mcp
   ```

### 2. Build the project:
   **Windows**
   ```bash
   go build -o chromedp-mcp.exe .\cmd\chromedp-mcp\main.go
   ```
   **Linux/MacOs**
   ```bash
   go build -o chromedp-mcp ./cmd/chromedp-mcp/main.go
   ```


### 3. MCP Client Configuration

Configure your MCP-compatible client to use this server:

**Basic Configuration:**
```json
{
  "mcpServers": {
    "chromedp-mcp": {
      "command": "/path/to/chromedp-mcp",
      "args": []
    }
  }
}
```

**With Environment Variables:**
```json
{
  "mcpServers": {
    "chromedp-mcp": {
      "command": "/path/to/chromedp-mcp",
      "args": [],
      "env": {
        "CHROME_MAXIMUM_INSTANCE": "10",
        "CHROME_TTL": "30",
        "CHROME_EXE_TIMEOUT": "600"
      }
    }
  }
}
```

## Important Notes

- Always close instances when done to free resources
- Respect website terms of service and robots.txt
- Some websites may have anti-automation measures



