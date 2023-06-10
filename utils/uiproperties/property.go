package uiproperties

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"golang.org/x/image/colornames"
)

type IPropertyOwner interface {
	Subclass() string
}

type IPropertiesContainer interface {
	AddProperty(name string, prop *Property)
	GetProperties() []*Property
	SetPropertyValue(name string, value interface{})
	PropertyValue(name string) interface{}
	Property(name string) *Property

	SetPropertyChangeNotifier(OnPropertyChangedForEditor func(prop *Property))
	NotifyChangedToContainer(prop *Property)
}

type PropertiesContainer struct {
	properties        []*Property
	propertiesMap     map[string]*Property
	OnPropertyChanged func(prop *Property)

	OnPropertyChangedForEditor func(prop *Property)
}

type PropertyType string

const (
	PropertyTypeBool      PropertyType = "bool"
	PropertyTypeInt       PropertyType = "int"
	PropertyTypeInt32     PropertyType = "int32"
	PropertyTypeColor     PropertyType = "color"
	PropertyTypeDouble    PropertyType = "double"
	PropertyTypeString    PropertyType = "string"
	PropertyTypeMultiline PropertyType = "multiline"
)

type Property struct {
	Name          string
	DisplayName   string
	Type          PropertyType
	SubType       string
	GroupName     string
	isUnstyled    bool
	unstyledValue interface{}
	DefaultValue  interface{}
	subclasses    map[string]*PropertyValue
	visible       bool
	propertyOwner IPropertyOwner
	OnChanged     func(property *Property, oldValue interface{}, newValue interface{})
	//PropOwner     IPropertiesContainer
}

type PropertyValue struct {
	Source string
	Value  interface{}
	Score  int
}

func NewProperty(name string, propertyType PropertyType) *Property {
	var p Property
	p.Init(name, nil)
	p.Type = propertyType
	p.visible = true
	return &p
}

func (c *Property) Init(name string, w IPropertyOwner) {
	c.subclasses = make(map[string]*PropertyValue)
	c.Name = name
	c.propertyOwner = w
}

func (c *Property) SetUnstyled(isUnstyled bool) {
	c.isUnstyled = isUnstyled
}

func (c *Property) Dispose() {
	c.propertyOwner = nil
}

func (c *Property) SetVisible(visible bool) {
	c.visible = visible
}

func (c *Property) Visible() bool {
	return c.visible
}

func (c *Property) SetOwnValue(value interface{}) {
	var realValue interface{}

	if c.Type == PropertyTypeBool {
		if v, ok := value.(bool); ok {
			realValue = v
		} else {
			iValue, err := strconv.Atoi(fmt.Sprint(value))
			if err == nil {
				realValue = iValue != 0
			} else {
				if fmt.Sprint(value) == "1" || fmt.Sprint(value) == "true" {
					realValue = true
				} else {
					realValue = false
				}
			}
		}
	}

	if c.Type == PropertyTypeInt {
		if v, ok := value.(int); ok {
			realValue = v
		} else {
			iValue, err := strconv.Atoi(fmt.Sprint(value))
			if err == nil {
				realValue = iValue
			}
		}
	}

	if c.Type == PropertyTypeInt32 {
		if v, ok := value.(int32); ok {
			realValue = v
		} else {
			fValue, err := strconv.ParseFloat(fmt.Sprint(value), 64)
			if err != nil {
				//panic("prop int32 - f64")
			}
			iValue := int32(fValue)
			realValue = int32(iValue)
		}
	}

	if c.Type == PropertyTypeColor {
		if v, ok := value.(color.Color); ok {
			realValue = v
		} else {
			realValue = ParseHexColor(fmt.Sprint(value))
		}
	}

	if c.Type == PropertyTypeDouble {
		if v, ok := value.(float64); ok {
			realValue = v
		} else {
			iValue, err := strconv.ParseFloat(fmt.Sprint(value), 64)
			if err == nil {
				realValue = iValue
			}
		}
	}

	if c.Type == PropertyTypeString {
		if v, ok := value.(string); ok {
			realValue = v
		} else {
			realValue = fmt.Sprint(value)
		}
	}

	if c.Type == PropertyTypeMultiline {
		if v, ok := value.(string); ok {
			realValue = v
		} else {
			realValue = fmt.Sprint(value)
		}
	}

	if realValue == nil && value != nil {
		panic("Can not parse value for set property")
	}

	oldValue := c.unstyledValue

	if value == nil {
		c.unstyledValue = nil
		c.isUnstyled = false
	} else {
		c.unstyledValue = realValue
		c.isUnstyled = true
	}

	if c.OnChanged != nil {
		c.OnChanged(c, oldValue, c.unstyledValue)
	}
}

func (c *Property) SetStyledValue(subclass string, value string, score int) {
	var propValue *PropertyValue
	var ok bool

	valueObject, err := ParseCSSProperty(c.Type, value)
	if err != nil {
		return
	}

	propValue, ok = c.subclasses[subclass]
	if !ok {
		c.subclasses[subclass] = new(PropertyValue)
		propValue = c.subclasses[subclass]
	}

	//if c.unstyledValue == nil {
	propValue.Value = valueObject
	propValue.Score = score
	//}
}

func (c *Property) Int() int {
	v := c.Value()
	if v != nil {
		if val, ok := v.(int); ok {
			return val
		} else {
			panic(fmt.Errorf("Property is not int"))
		}
	}
	return 0
}

func (c *Property) Int32() int32 {
	v := c.Value()
	if v != nil {
		if val, ok := v.(int32); ok {
			return val
		} else {
			parsedFloat64, err := strconv.ParseFloat(fmt.Sprint(v), 64)
			if err != nil {
				return int32(parsedFloat64)
			}
			return 0
		}
	}
	return 0
}

func (c *Property) Bool() bool {
	v := c.Value()
	if v != nil {
		if val, ok := v.(bool); ok {
			return val
		} else {
			panic(fmt.Errorf("Property is not bool"))
		}
	}
	return false
}

func (c *Property) String() string {
	v := c.Value()
	if v != nil {
		if val, ok := v.(string); ok {
			return val
		} else {
			panic(fmt.Errorf("Property is not string"))
		}
	}
	return ""
}

func (c *Property) Code() string {
	v := c.Value()
	if v != nil {
		if val, ok := v.(string); ok {
			return val
		} else {
			panic(fmt.Errorf("Property is not code"))
		}
	}
	return ""
}

func (c *Property) Float64() float64 {
	v := c.Value()
	if v != nil {
		if val, ok := v.(float64); ok {
			return val
		} else {
			panic(fmt.Errorf("Property is not float64"))
		}
	}
	return 0
}

func (c *Property) Color() color.Color {

	v := c.Value()
	if v != nil {
		if val, ok := v.(color.Color); ok {
			return val
		} else {
			panic(fmt.Errorf("Property is not color"))
		}
	}
	return color.RGBA{}
}

func (c *Property) Value() interface{} {
	subclass := "default"

	var val interface{}

	if c.propertyOwner != nil {
		subclass = c.propertyOwner.Subclass()
	}

	if c.isUnstyled {
		return c.unstyledValue
	}

	val = c.ValueOwn() // Unsubclasses value
	if val == nil {
		val = c.ValueStyle(subclass)
	}
	if val == nil {
		val = c.ValueStyle("default")
	}

	return val
}

func (c *Property) DefaultValueForType(propertyType PropertyType) interface{} {
	switch propertyType {
	case PropertyTypeBool:
		return false
	case PropertyTypeInt:
		return int(0)
	case PropertyTypeInt32:
		return int32(0)
	case PropertyTypeColor:
		return colornames.Black
	case PropertyTypeDouble:
		return float64(0)
	case PropertyTypeString:
		return ""
	case PropertyTypeMultiline:
		return ""
	}
	return nil
}

func (c *Property) ValueOwn() interface{} {
	if c.isUnstyled {
		return c.unstyledValue
	}
	return nil
}

type PropertyStruct struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (c *Property) SaveToStruct() *PropertyStruct {
	var result PropertyStruct
	value := ""
	if c.ValueOwn() != nil {
		switch c.Type {
		case PropertyTypeBool:
			value = fmt.Sprint(c.ValueOwn().(bool))
		case PropertyTypeInt:
			value = fmt.Sprint(int64(c.ValueOwn().(int)))
		case PropertyTypeInt32:
			value = fmt.Sprint(c.ValueOwn().(int32))
		case PropertyTypeDouble:
			value = fmt.Sprint(c.ValueOwn().(float64))
		case PropertyTypeString:
			value = fmt.Sprint(c.ValueOwn().(string))
		case PropertyTypeMultiline:
			value = fmt.Sprint(c.ValueOwn().(string))
		case PropertyTypeColor:
			r, g, b, a := c.ValueOwn().(color.Color).RGBA()
			value = fmt.Sprintf("#%02X%02X%02X%02X", r/256, g/256, b/256, a/256)
		}
	}

	result.Name = c.Name
	result.Type = string(c.Type)
	result.Value = value

	return &result
}

func (c *Property) ValueStyle(subclass string) interface{} {
	if s, ok := c.subclasses[subclass]; ok {
		if s.Value != nil {
			return s.Value
		}
	}
	return nil
}

func (c *Property) ValueStyleScore(subclass string) int {
	if s, ok := c.subclasses[subclass]; ok {
		if s.Value != nil {
			return s.Score
		}
	}
	return 0
}

func ParseCSSProperty(propType PropertyType, value string) (interface{}, error) {
	var result interface{}
	var err error
	var resInt64 int64

	switch propType {
	case PropertyTypeBool:
		if value == "true" {
			result = true
		} else {
			if value == "false" {
				result = false
			} else {
				result = false
				err = errors.New("css type bool - mismatch types")
			}
		}
	case PropertyTypeInt:
		resInt64, err = strconv.ParseInt(value, 10, 64)
		result = int(resInt64)
	case PropertyTypeInt32:
		resInt64, err = strconv.ParseInt(value, 10, 32)
		result = int32(resInt64)
	case PropertyTypeColor:
		result = ParseHexColor(value)
	case PropertyTypeDouble:
		result, err = strconv.ParseFloat(value, 64)
	case PropertyTypeString:
		result = value
	}

	return result, err
}

// ToRGB converts the HEXColor to and RGBColor
func ParseHexColor(hexColor string) color.Color {

	var r, g, b, a uint8

	a = 255

	hexFormatRGB := "#%02x%02x%02x"
	hexFormatRGBA := "#%02x%02x%02x%02x"
	hexShortFormatRGB1 := "#%1x%1x%1x"
	hexShortFormatRGBA1 := "#%1x%1x%1x%1x"

	if len(hexColor) == 4 {
		fmt.Sscanf(hexColor, hexShortFormatRGB1, &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
	}
	if len(hexColor) == 5 {
		fmt.Sscanf(hexColor, hexShortFormatRGBA1, &a, &r, &g, &b)
		r *= 17
		g *= 17
		b *= 17
	}
	if len(hexColor) == 7 {
		fmt.Sscanf(hexColor, hexFormatRGB, &r, &g, &b)
	}
	if len(hexColor) == 9 {
		fmt.Sscanf(hexColor, hexFormatRGBA, &a, &r, &g, &b)
	}

	return color.RGBA{r, g, b, a}
}

func (c *PropertiesContainer) InitPropertiesContainer() {
	c.propertiesMap = make(map[string]*Property)
	c.properties = make([]*Property, 0)
}

func (c *PropertiesContainer) SetPropertyValue(name string, value interface{}) {
	if prop, ok := c.propertiesMap[name]; ok {
		prop.SetOwnValue(value)
		if c.OnPropertyChanged != nil {
			c.OnPropertyChanged(prop)
		}
		if c.OnPropertyChangedForEditor != nil {
			c.OnPropertyChangedForEditor(prop)
		}
	}
}

func (c *PropertiesContainer) PropertyValue(name string) interface{} {
	if prop, ok := c.propertiesMap[name]; ok {
		return prop.ValueOwn()
	}
	return nil
}

func (c *PropertiesContainer) GetProperties() []*Property {
	return c.properties
}

func (c *PropertiesContainer) Property(name string) *Property {
	for _, prop := range c.properties {
		if prop.Name == name {
			return prop
		}
	}

	return nil
}

func (c *PropertiesContainer) AddProperty(name string, prop *Property) {
	c.propertiesMap[name] = prop
	c.properties = append(c.properties, prop)
	//prop.PropOwner = c
}

func (c *PropertiesContainer) SetPropertyChangeNotifier(OnPropertyChangedForEditor func(prop *Property)) {
	c.OnPropertyChangedForEditor = OnPropertyChangedForEditor
}

func (c *PropertiesContainer) NotifyChangedToContainer(prop *Property) {
	if c.OnPropertyChangedForEditor != nil {
		c.OnPropertyChangedForEditor(prop)
	}
}
