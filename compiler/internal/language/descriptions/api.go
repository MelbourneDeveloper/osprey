// Package descriptions provides documentation for built-in language elements.
package descriptions

import (
	"strings"
)

// LanguageElement represents any element in the language that can have documentation.
type LanguageElement struct {
	Type        string // "function", "type", "operator", "keyword"
	Name        string
	Description string
	Example     string
	Signature   string // For functions only
}

// GetAllLanguageElements returns all documented language elements.
func GetAllLanguageElements() map[string]*LanguageElement {
	elements := make(map[string]*LanguageElement)

	// Add functions
	for name, desc := range GetBuiltinFunctionDescriptions() {
		elements[name] = &LanguageElement{
			Type:        "function",
			Name:        desc.Name,
			Description: desc.Description,
			Example:     desc.Example,
			Signature:   desc.Signature,
		}
	}

	// Add types
	for name, desc := range GetBuiltinTypeDescriptions() {
		elements[name] = &LanguageElement{
			Type:        "type",
			Name:        desc.Name,
			Description: desc.Description,
			Example:     desc.Example,
		}
	}

	// Add operators
	for symbol, desc := range GetOperatorDescriptions() {
		elements[symbol] = &LanguageElement{
			Type:        "operator",
			Name:        desc.Name,
			Description: desc.Description,
			Example:     desc.Example,
		}
	}

	// Add keywords
	for keyword, desc := range GetKeywordDescriptions() {
		elements[keyword] = &LanguageElement{
			Type:        "keyword",
			Name:        desc.Keyword,
			Description: desc.Description,
			Example:     desc.Example,
		}
	}

	return elements
}

// GetLanguageElementDescription returns description for any language element.
func GetLanguageElementDescription(name string) *LanguageElement {
	elements := GetAllLanguageElements()
	if element, exists := elements[name]; exists {
		return element
	}

	return nil
}

// GetHoverDocumentation returns hover documentation for any language element.
func GetHoverDocumentation(name string) string {
	element := GetLanguageElementDescription(name)
	if element == nil {
		return ""
	}

	var parts []string

	switch element.Type {
	case "function":
		if element.Signature != "" {
			parts = append(parts, "```osprey\n"+element.Signature+"\n```")
		}

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "\n**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	case "type":
		parts = append(parts, "```osprey\ntype "+element.Name+"\n```")

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "\n**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	case "operator":
		parts = append(parts, "**Operator:** `"+name+"`")

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "\n**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	case "keyword":
		parts = append(parts, "**Keyword:** `"+element.Name+"`")

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "\n**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	}

	return strings.Join(parts, "\n\n")
}
