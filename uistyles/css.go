package uistyles

import (
	"github.com/gorilla/css/scanner"
	"strings"
)

type CSSLine struct {
	elementName  string
	elementType  string
	class        string
	subclass     string
	propertyName string
	value        string
}

type CSS struct {
	lines []CSSLine
}

type CSSProperty struct {
	name  string
	value string
}

func (c *CSS) Parse(cssText string) {
	s := scanner.New(cssText)

	selectors := make([]string, 0)
	var inSelectors bool
	inSelectors = true

	var currentSelector string

	var inPropName bool
	inPropName = false
	var currentPropName string
	var currentPropValue string

	properties := make([]CSSProperty, 0)

	for {
		token := s.Next()
		if token.Type == scanner.TokenEOF || token.Type == scanner.TokenError {
			break
		}
		//fmt.Println(token)

		if inSelectors {
			if token.Type == scanner.TokenIdent || token.Type == scanner.TokenHash || (token.Type == scanner.TokenChar && token.Value == ":") {
				currentSelector += token.Value
			}

			if token.Type == scanner.TokenChar && (token.Value == "{" || token.Value == ",") {
				selectors = append(selectors, currentSelector)
				currentSelector = ""
			}
			if token.Type == scanner.TokenChar && token.Value == "{" {
				inPropName = true
				inSelectors = false
			}
		} else {
			if inPropName {
				if token.Type != scanner.TokenChar && token.Value != "}" {
					if token.Type == scanner.TokenIdent {
						currentPropName = token.Value
						inPropName = false
						currentPropValue = ""
					} else {
						if token.Type != scanner.TokenS {
							println("Error")
							return
						}
					}
				}
			} else {

				if token.Type == scanner.TokenNumber || token.Type == scanner.TokenDimension || token.Type == scanner.TokenString || token.Type == scanner.TokenIdent || token.Type == scanner.TokenHash {
					currentPropValue += token.Value
				}

				if token.Type == scanner.TokenChar && (token.Value == ";" || token.Value == "}") {

					properties = append(properties, CSSProperty{currentPropName, currentPropValue})
					currentPropName = ""
					currentPropValue = ""
					inPropName = true
				}
			}

		}

		if token.Type == scanner.TokenChar && token.Value == "}" {
			for _, selector := range selectors {
				for _, property := range properties {
					var line CSSLine
					line.elementName = c.getSelectorElementName(selector)
					line.elementType = c.getSelectorElementType(selector)
					line.subclass = c.getSelectorSubclass(selector)
					line.propertyName = property.name
					line.value = property.value
					c.lines = append(c.lines, line)
				}
			}

			properties = make([]CSSProperty, 0)
			selectors = make([]string, 0)
			inSelectors = true
		}

	}

	//println(c.lines)
}

// tr#id:hover
func (c *CSS) getSelectorElementName(selector string) string {
	fieldsWithoutSubclass := strings.FieldsFunc(selector, func(c rune) bool {
		return c == ':'
	})
	if len(fieldsWithoutSubclass) > 0 {

		fields := strings.FieldsFunc(fieldsWithoutSubclass[0], func(c rune) bool {
			return c == '#'
		})

		if len(fields) == 2 {
			return fields[1]
		}
	}
	return ""
}

func (c *CSS) getSelectorSubclass(selector string) string {
	f := func(c rune) bool {
		return c == ':'
	}
	fields := strings.FieldsFunc(selector, f)
	if len(fields) == 2 {
		return fields[1]
	}
	return "default"
}

func (c *CSS) getSelectorElementType(selector string) string {
	fieldsWithoutSubclass := strings.FieldsFunc(selector, func(c rune) bool {
		return c == ':'
	})
	if len(fieldsWithoutSubclass) > 0 {

		fields := strings.FieldsFunc(fieldsWithoutSubclass[0], func(c rune) bool {
			return c == '#'
		})

		if len(fields) == 2 {
			return fields[0]
		}
		if len(fields) == 1 {
			return fields[0]
		}
	}
	return "Control"
}
