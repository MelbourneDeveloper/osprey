// Package descriptions provides documentation for built-in language elements.
package descriptions

import (
	"strings"
)

// Element type constants for language documentation
const (
	ElementTypeFunction = "function"
	ElementTypeType     = "type"
	ElementTypeOperator = "operator"
	ElementTypeKeyword  = "keyword"
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
			Type:        ElementTypeFunction,
			Name:        desc.Name,
			Description: desc.Description,
			Example:     desc.Example,
			Signature:   desc.Signature,
		}
	}

	// Add types
	for name, desc := range GetBuiltinTypeDescriptions() {
		elements[name] = &LanguageElement{
			Type:        ElementTypeType,
			Name:        desc.Name,
			Description: desc.Description,
			Example:     desc.Example,
		}
	}

	// Add operators
	for symbol, desc := range GetOperatorDescriptions() {
		elements[symbol] = &LanguageElement{
			Type:        ElementTypeOperator,
			Name:        desc.Name,
			Description: desc.Description,
			Example:     desc.Example,
		}
	}

	// Add keywords
	for keyword, desc := range GetKeywordDescriptions() {
		elements[keyword] = &LanguageElement{
			Type:        ElementTypeKeyword,
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
	case ElementTypeFunction:
		if element.Signature != "" {
			parts = append(parts, "```osprey\n"+element.Signature+"\n```")
		}

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	case ElementTypeType:
		parts = append(parts, "```osprey\ntype "+element.Name+"\n```")

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	case ElementTypeOperator:
		parts = append(parts, "**Operator:** `"+name+"`")

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	case ElementTypeKeyword:
		parts = append(parts, "**Keyword:** `"+element.Name+"`")

		parts = append(parts, element.Description)
		if element.Example != "" {
			parts = append(parts, "**Example:**\n```osprey\n"+element.Example+"\n```")
		}
	}

	return strings.Join(parts, "\n\n")
}
