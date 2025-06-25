# chromedp-mcp

## ⚠️ Sample Status

This project is currently in **SAMPLE** status and only implements basic PDF generation functionality.

## Features

- Generate PDF from HTML content
- Generate PDF from URLs
- Specify output path for PDF file
- Basic MCP server integration

## Tools Provided

This MCP server currently provides the following tool:

### `generate_pdf`
**Description:** Generate PDF from HTML content or URL

**Parameters:**
- `html` (string, optional): HTML string to generate PDF from
- `url` (string, optional): URL to generate PDF from  
- `outputDir` (string, optional): Output directory path for the PDF file (defaults to `~/pdf_output`)

**Requirements:** Either `html` or `url` parameter must be provided

**Output:** Returns the file path of the generated PDF document

## Requirements

- **Chrome Engine**: Chrome/Chromium browser installed on system
- **Go**: Go 1.19+ required for building

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/KePatrick/chromedp-mcp.git

   cd chromedp-mcp
   ```

2. Build the project:
   **Windows**
   ```bash
   go build -o chromedp-mcp.exe .\cmd\chromedp-mcp\main.go
   ```
   **Linux/MacOs**
   ```bash
   go build -o chromedp-mcp ./cmd/chromedp-mcp/main.go
   ```


3. Configure MCP settings:
   ```json
   {
	"mcpServers": {
		   "chromedp-mcp": {
			   "command": "path-to-your-mcp/chromedp-mcp",
			   "args": []
	
		}
	}
   }
   ```

## Future Planning

1. **Reorganize project structure** - Improve code organization and modularity

2. **Implement more chromedp tools** - Add additional browser automation capabilities

3. **Docker packaging** - Package environment with Chromium and native executable for easier deployment
