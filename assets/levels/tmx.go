package levels

import (
	"fmt"
	"strconv"

	"github.com/abelroes/gmtk2024/src/vector"
)

type Map struct {
	ObjectGroups []ObjectGroup `xml:"objectgroup"`
}

type ObjectGroup struct {
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Objects []Object `xml:"object"`
}

type Object struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`

	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Width  float64 `xml:"width,attr"`
	Height float64 `xml:"height,attr"`
	Props  *Props  `xml:"properties"`
}

type Props struct {
	Props []Property `xml:"property"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

func (group *ObjectGroup) FindObjectByName(name string) *Object {
	for _, obj := range group.Objects {
		if obj.Name == name {
			return &obj
		}
	}
	return nil
}

func (object Object) TopLeftPos() vector.Vector2 {
	return vector.New(object.X, object.Y-object.Height)
}

func (object Object) CenterPos() vector.Vector2 {
	return vector.New(object.X+object.Width/2, object.Y-object.Height/2)
}

func (props *Props) GetProp(name string) *Property {
	for _, p := range props.Props {
		if p.Name == name {
			return &p
		}
	}

	return nil
}

func (props *Props) GetPropString(name string) (string, error) {
	p := props.GetProp(name)
	if p == nil {
		return "", fmt.Errorf("prop %s not found", name)
	}
	return p.Value, nil
}

func (props *Props) GetPropInt(name string) (int64, error) {
	p := props.GetProp(name)
	if p == nil {
		return 0, fmt.Errorf("prop %s not found", name)
	}
	return strconv.ParseInt(p.Value, 10, 64)
}

func (props *Props) GetPropFloat(name string) (float64, error) {
	p := props.GetProp(name)
	if p == nil {
		return 0, fmt.Errorf("prop %s not found", name)
	}
	return strconv.ParseFloat(p.Value, 64)
}
