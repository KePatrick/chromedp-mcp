# ChromeDP MCP - Tools Specification

This document provides detailed specifications for all available tools in the ChromeDP MCP server.

## Instance Management

### `create_chrome_instance`

**Description**: Create Chrome Instance, every session should start by create_chrome_instance and end by close

**Parameters**:
- `headless` (boolean, optional): Headless mode flag for create chrome instance (default: true)
- `disable-gpu` (boolean, optional): Disable gpu for chrome instance (default: true)
- `disable-popup-blocking` (boolean, optional): Disable popup blocking to allow popups (default: false, meaning popups are blocked)
- `block-new-tab` (boolean, optional): Block opening new tabs/windows and redirect to current tab (default: false)
- `disable-extensions` (boolean, optional): Disable browser extensions (default: true)
- `disable-plugins` (boolean, optional): Disable browser plugins (default: true)
- `disable-web-security` (boolean, optional): Disable web security (CORS) for testing purposes (default: false)
- `no-sandbox` (boolean, optional): Disable Chrome sandbox for running in containers (default: false)
- `disable-dev-shm-usage` (boolean, optional): Disable /dev/shm usage to avoid memory issues in containers (default: false)
- `disable-background-timer-throttling` (boolean, optional): Disable background timer throttling for better performance (default: false)
- `disable-backgrounding-occluded-windows` (boolean, optional): Disable backgrounding occluded windows (default: false)
- `disable-renderer-backgrounding` (boolean, optional): Disable renderer backgrounding (default: false)

**Returns**: Chrome instance ID

---

### `close`

**Description**: Close Chrome Instance, every session should start by create_chrome_instance and end by close

**Parameters**:
- `id` (string, optional): The ID of the Chrome instance to close

**Returns**: Success message

---

## Page Navigation

### `navigate`

**Description**: Navigate to provided URL and return a clean DOM element tree structure without scripts/styles and textContent

**Parameters**:
- `url` (string, required): The URL to navigate
- `id` (string, required): Chrome instance id
- `depth` (number, optional): Maximum DOM tree depth to traverse (default: 5)

**Returns**: Clean DOM tree structure

---

### `navigate-back`

**Description**: Navigate to previous page

**Parameters**:
- `id` (string, required): Chrome instance id

**Returns**: Current URL

---

### `navigate-forward`

**Description**: Navigate to next page

**Parameters**:
- `id` (string, required): Chrome instance id

**Returns**: Current URL

---

## Element Operations

### `get-element-withtext`

**Description**: Get element specified by selector (CSS selector, XPath, ID, class, etc.) with text. You should create a instance before operation

**Parameters**:
- `selector` (string, required): The selector to identify the element. Examples: '#button-id', '.button-class', 'button[type="submit"]', '//button[@id="submit"]'
- `id` (string, required): Chrome instance ID to perform the operation on

**Returns**: Element information with text content

---

### `get-all-elements`

**Description**: Get all elements of current page, return a clean DOM element tree structure without scripts/styles and textContent

**Parameters**:
- `id` (string, required): Chrome instance id
- `depth` (number, optional): Maximum DOM tree depth to traverse (default: 5)

**Returns**: Complete DOM tree structure

---

### `select-element`

**Description**: Select element by CSS selector and return clean DOM structure at specified depth without text content

**Parameters**:
- `selector` (string, required): The selector to identify the element. Examples: '#button-id', '.button-class', 'button[type="submit"]', '//button[@id="submit"]'
- `id` (string, required): Chrome instance ID
- `depth` (number, optional): Maximum depth to traverse from selected element (default: 3, max: 10)
- `all` (boolean, optional): Select all matching elements instead of just the first one (default: false)

**Returns**: Clean DOM structure at specified depth

---

### `click-element`

**Description**: Click on an element specified by selector (CSS selector, XPath, ID, class, etc.). Returns Success message and any errors.

**Parameters**:
- `selector` (string, required): The selector to identify the element to click. Examples: '#button-id', '.button-class', 'button[type="submit"]', '//button[@id="submit"]'
- `id` (string, required): Chrome instance ID to perform the click action on
- `timeout` (number, optional): Timeout in seconds to wait for element (default: 10)
- `wait_visible` (boolean, optional): Whether to wait for element to be visible before clicking (default: true)
- `click_type` (string, optional): Type of click: 'left' (default), 'right', 'double'

**Returns**: Success message and error information

---

## Input Operations

### `send-key`

**Description**: The send-key tool in chromedp simulates keyboard input by sending keystrokes to a specified element by selector on a web page.

**Parameters**:
- `id` (string, required): Chrome instance ID to perform the operation on
- `selector` (string, required): The selector to identify the element. Examples: '#button-id', '.button-class', 'button[type="submit"]', '//button[@id="submit"]'
- `key` (string, optional): The input value

**Returns**: Success message and new element status

---

### `set-value`

**Description**: The set-value tool in chromedp directly sets the value of a specified element by selector on a web page. This is more efficient than send-key for setting form field values.

**Parameters**:
- `id` (string, required): Chrome instance ID to perform the set value action on
- `selector` (string, required): The selector to identify the element to set value. Examples: '#input-id', '.input-class', 'input[name="username"]', '//input[@id="email"]'
- `value` (string, required): The value to set for the element

**Returns**: Success message and new element status

---

### `key-event`

**Description**: The key-event tool simulates keyboard events in chromedp, supporting single keys and key combinations with modifiers.

**Parameters**:
- `id` (string, required): Chrome instance ID to perform the key event on
- `key` (string, required): The target key to press. Supported keys include:
  - **Alphanumeric**: a-z, A-Z, 0-9
  - **Special characters**: space, !, @, #, $, %, ^, &, *, (, ), -, _, =, +, [, ], {, }, |, \, :, ;, ", ', <, >, ,, ., ?, /
  - **Navigation**: ArrowUp, ArrowDown, ArrowLeft, ArrowRight, Home, End, PageUp, PageDown
  - **Function keys**: F1-F24
  - **Control keys**: Tab, Enter, Escape, Backspace, Delete, Insert
  - **Modifier keys**: Alt, Control, Meta, Shift, CapsLock, NumLock, ScrollLock
  - **Media keys**: MediaPlayPause, MediaStop, MediaTrackNext, MediaTrackPrevious, AudioVolumeUp, AudioVolumeDown, AudioVolumeMute
  - **System keys**: PrintScreen, Pause, ContextMenu, Copy, Cut, Paste, Undo, Redo, Find, Help
  - **Browser keys**: BrowserBack, BrowserForward, BrowserRefresh, BrowserHome, BrowserSearch
  - Examples: 'a', 'Enter', 'F5', 'ArrowUp', 'MediaPlayPause'
- `modifier` (string, optional): Optional modifier keys separated by semicolons. Available modifiers: 'ctrl', 'shift', 'alt', 'meta'. Examples: 'ctrl' for Ctrl+key, 'ctrl;shift' for Ctrl+Shift+key, 'ctrl;shift;alt' for three-key combinations.

**Returns**: Success message

---

## Cookie Management

### `set-cookie`

**Description**: Set a HTTP cookie on requests. Cookies will be automatically sent with subsequent requests to matching domains and paths.

**Parameters**:
- `id` (string, required): Chrome instance ID to perform the set cookie action on
- `name` (string, required): The name of the cookie to set
- `value` (string, required): The value of the cookie to set
- `domain` (string, required): The domain to set the cookie for (e.g., 'example.com', '.example.com' for subdomains)
- `path` (string, required): The path to set the cookie for. Use '/' for all paths under the domain
- `httpOnly` (boolean, optional): If true, the cookie is only accessible via HTTP requests and not JavaScript. Recommended for security-sensitive cookies like session tokens. Default: false
- `secure` (boolean, optional): If true, the cookie is only sent over HTTPS connections. Recommended for production environments. Default: false
- `sameSite` (string, optional): Controls when the cookie is sent with cross-site requests. Options: 'Strict' (only same-site), 'Lax' (some cross-site), 'None' (all cross-site, requires Secure=true). Default: 'Lax'
- `expires` (number, optional): Cookie expiration time as Unix timestamp (seconds since epoch). If not set, cookie expires when browser session ends. Example: 1735689600 for 2025-01-01
- `maxAge` (number, optional): Cookie lifetime in seconds from when it's set. Takes precedence over 'expires' if both are set. Example: 3600 for 1 hour, 86400 for 1 day

**Returns**: Success message

---

## File Operations

### `download-file`

**Description**: The download-file tool in chromedp downloads files by clicking on download links or buttons. It can optionally specify a download directory.

**Parameters**:
- `id` (string, required): Chrome instance ID to perform the download action on
- `selector` (string, required): The selector to identify the download element to click. Examples: '#download-link', '.download-btn', 'a[href*="download"]', '//a[contains(@href, "download")]'
- `download_path` (string, optional): Optional download directory path. If not specified, uses user's default Downloads directory
- `timeout` (number, optional): Download timeout in seconds (default: 30)

**Returns**: Download status and file path

---

### `download_image`

**Description**: Download image from URL or selector

**Parameters**:
- `id` (string, required): Chrome instance ID
- `url` (string, optional): Image URL to download
- `selector` (string, optional): The selector to identify the element. Examples: '#button-id', '.button-class', 'button[type="submit"]', '//button[@id="submit"]'
- `output_path` (string, optional): Output directory path (default: user downloads directory)

**Note**: Either `url` or `selector` parameter must be provided

**Returns**: Download status and file path

---

## Document Generation

### `generate_pdf`

**Description**: Generate PDF from HTML content or URL

**Parameters**:
- `html` (string, optional): HTML string to generate PDF from
- `url` (string, optional): URL to generate PDF from
- `outputDir` (string, optional): Output directory path for the PDF file

**Note**: Either `html` or `url` parameter must be provided

**Returns**: Generated PDF file path

---

## Other

### `tips`

**Description**: Get important usage tips and best practices for Chrome automation tools, see this before you start

**Parameters**: None

**Returns**: Usage tips and best practices

