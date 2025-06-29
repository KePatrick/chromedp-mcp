package tool

import (
	"context"
	"fmt"
	"strings"

	mcpcdp "patrickke/chromedp-mcp/chromedp"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewKeyEventTool() mcp.Tool {
	return mcp.NewTool("key-event",
		mcp.WithDescription("The key-event tool simulates keyboard events in chromedp, supporting single keys and key combinations with modifiers."),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Chrome instance ID to perform the key event on"),
		),
		mcp.WithString("key", 
			mcp.Required(),
			mcp.Description("The target key to press. Supported keys include:\n" +
				"• Alphanumeric: a-z, A-Z, 0-9\n" +
				"• Special characters: space, !, @, #, $, %, ^, &, *, (, ), -, _, =, +, [, ], {, }, |, \\, :, ;, \", ', <, >, ,, ., ?, /\n" +
				"• Navigation: ArrowUp, ArrowDown, ArrowLeft, ArrowRight, Home, End, PageUp, PageDown\n" +
				"• Function keys: F1-F24\n" +
				"• Control keys: Tab, Enter, Escape, Backspace, Delete, Insert\n" +
				"• Modifier keys: Alt, Control, Meta, Shift, CapsLock, NumLock, ScrollLock\n" +
				"• Media keys: MediaPlayPause, MediaStop, MediaTrackNext, MediaTrackPrevious, AudioVolumeUp, AudioVolumeDown, AudioVolumeMute\n" +
				"• System keys: PrintScreen, Pause, ContextMenu, Copy, Cut, Paste, Undo, Redo, Find, Help\n" +
				"• Browser keys: BrowserBack, BrowserForward, BrowserRefresh, BrowserHome, BrowserSearch\n" +
				"Examples: 'a', 'Enter', 'F5', 'ArrowUp', 'MediaPlayPause'"),
		),
		mcp.WithString("modifier", 
			mcp.Description("Optional modifier keys separated by semicolons. Available modifiers: 'ctrl', 'shift', 'alt', 'meta'. Examples: 'ctrl' for Ctrl+key, 'ctrl;shift' for Ctrl+Shift+key, 'ctrl;shift;alt' for three-key combinations."),
		),
	)
}


func KeyEventHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("id", "")
	key := request.GetString("key", "")
	modifier := request.GetString("modifier", "")
	
	// Validate required parameters
	if id == "" {
		return mcp.NewToolResultError("Chrome instance id is required"), nil
	}
	if key == "" {
		return mcp.NewToolResultError("key is required"), nil
	}
	
	// Map key string to kb constant
	mappedKey, err := KeyMapper(key)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("key [%s] is not available, err: %v", key, err)), nil
	}
	
	// Handle case without modifiers
	if modifier == "" {
		err := mcpcdp.Manager.Execute(id, 
			chromedp.KeyEvent(mappedKey),
		)
		if err != nil {
			return nil, fmt.Errorf("fail to execute key event, err: %v", err)
		}
		
		return mcp.NewToolResultText(fmt.Sprintf("key event is triggered: %s", key)), nil
	}
	
	// Handle case with modifiers
	modifiers := strings.Split(modifier, ";")
	keyModifiers := make([]input.Modifier, 0, len(modifiers)) 
	
	for _, mod := range modifiers {
		trimmedMod := strings.TrimSpace(mod)
		if trimmedMod == "" {
			continue // Skip empty modifiers
		}
		
		keyMod, err := modifierMapper(trimmedMod) // Fix: use trimmedMod, not modifier
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("modifier [%s] is not available, err: %v", trimmedMod, err)), nil
		}
		keyModifiers = append(keyModifiers, keyMod)
	}
	
	// Execute key event with modifiers
	err = mcpcdp.Manager.Execute(id, 
		chromedp.KeyEvent(mappedKey, chromedp.KeyModifiers(keyModifiers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to execute key event, err: %v", err)
	}
	
	// Format response message
	modifierDisplay := strings.ReplaceAll(modifier, ";", " + ")
	return mcp.NewToolResultText(fmt.Sprintf("key event is triggered: %s + %s", modifierDisplay, key)), nil
}

func modifierMapper(modifier string) (input.Modifier, error) {
	modifier = strings.TrimSpace(strings.ToLower(modifier))
	switch modifier {
	case "ctrl":
		return input.ModifierCtrl, nil
	case "shift":
		return input.ModifierShift, nil
	case "alt": 
		return input.ModifierAlt, nil
	case "meta":
		return input.ModifierMeta, nil
	}
	return 0, fmt.Errorf("unknown key: %s", modifier)
}

func KeyMapper(keyName string) (string, error) {
	// Normalize input - remove spaces and convert to lowercase
	keyName = strings.TrimSpace(strings.ToLower(keyName))
	
	// Key mapping table based on tool description
	keyMap := map[string]string{
		// Navigation keys
		"arrowup":    kb.ArrowUp,
		"arrowdown":  kb.ArrowDown,
		"arrowleft":  kb.ArrowLeft,
		"arrowright": kb.ArrowRight,
		"home":       kb.Home,
		"end":        kb.End,
		"pageup":     kb.PageUp,
		"pagedown":   kb.PageDown,
		
		// Function keys F1-F24
		"f1":  kb.F1,  "f2":  kb.F2,  "f3":  kb.F3,  "f4":  kb.F4,
		"f5":  kb.F5,  "f6":  kb.F6,  "f7":  kb.F7,  "f8":  kb.F8,
		"f9":  kb.F9,  "f10": kb.F10, "f11": kb.F11, "f12": kb.F12,
		"f13": kb.F13, "f14": kb.F14, "f15": kb.F15, "f16": kb.F16,
		"f17": kb.F17, "f18": kb.F18, "f19": kb.F19, "f20": kb.F20,
		"f21": kb.F21, "f22": kb.F22, "f23": kb.F23, "f24": kb.F24,
		
		// Control keys
		"tab":       kb.Tab,
		"enter":     kb.Enter,
		"escape":    kb.Escape,
		"backspace": kb.Backspace,
		"delete":    kb.Delete,
		"insert":    kb.Insert,
		
		// Modifier keys
		"alt":       kb.Alt,
		"control":   kb.Control,
		"meta":      kb.Meta,
		"shift":     kb.Shift,
		"capslock":  kb.CapsLock,
		"numlock":   kb.NumLock,
		"scrolllock": kb.ScrollLock,
		
		// Media keys
		"mediaplaypause":      kb.MediaPlayPause,
		"mediastop":           kb.MediaStop,
		"mediatracknext":      kb.MediaTrackNext,
		"mediatrackprevious":  kb.MediaTrackPrevious,
		"audiovolumeup":       kb.AudioVolumeUp,
		"audiovolumedown":     kb.AudioVolumeDown,
		"audiovolumemute":     kb.AudioVolumeMute,
		
		// System keys
		"printscreen":  kb.PrintScreen,
		"pause":        kb.Pause,
		"contextmenu":  kb.ContextMenu,
		"copy":         kb.Copy,
		"cut":          kb.Cut,
		"paste":        kb.Paste,
		"undo":         kb.Undo,
		"redo":         kb.Redo,
		"find":         kb.Find,
		"help":         kb.Help,
		
		// Browser keys
		"browserback":    kb.BrowserBack,
		"browserforward": kb.BrowserForward,
		"browserrefresh": kb.BrowserRefresh,
		"browserhome":    kb.BrowserHome,
		"browsersearch":  kb.BrowserSearch,
		
		// Special case for space
		"space": " ",
	}
	
	// Check if key exists in mapping
	if value, exists := keyMap[keyName]; exists {
		return value, nil
	}
	
	// Handle single characters (a-z, A-Z, 0-9, special characters)
	originalKey := strings.TrimSpace(keyName)
	if len(originalKey) == 1 {
		return originalKey, nil
	}
	
	// Return error for unknown keys
	return keyName, fmt.Errorf("unknown key: %s", keyName)
}

// GetSupportedKeys returns a list of all supported key names
func GetSupportedKeys() []string {
	return []string{
		// Alphanumeric (handled as single characters)
		"a-z", "A-Z", "0-9", "space",
		
		// Navigation
		"ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", 
		"Home", "End", "PageUp", "PageDown",
		
		// Function keys
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10",
		"F11", "F12", "F13", "F14", "F15", "F16", "F17", "F18", "F19", 
		"F20", "F21", "F22", "F23", "F24",
		
		// Control keys
		"Tab", "Enter", "Escape", "Backspace", "Delete", "Insert",
		
		// Modifier keys
		"Alt", "Control", "Meta", "Shift", "CapsLock", "NumLock", "ScrollLock",
		
		// Media keys
		"MediaPlayPause", "MediaStop", "MediaTrackNext", "MediaTrackPrevious",
		"AudioVolumeUp", "AudioVolumeDown", "AudioVolumeMute",
		
		// System keys
		"PrintScreen", "Pause", "ContextMenu", "Copy", "Cut", "Paste", 
		"Undo", "Redo", "Find", "Help",
		
		// Browser keys
		"BrowserBack", "BrowserForward", "BrowserRefresh", 
		"BrowserHome", "BrowserSearch",
	}
}
