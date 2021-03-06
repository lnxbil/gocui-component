package component

import "github.com/jroimartin/gocui"

type Key interface{}
type Handler func(g *gocui.Gui, v *gocui.View) error
type Handlers map[Key]Handler

type Component interface {
	GetLabel() string
	GetPosition() *Position
	GetType() ComponentType
	Focus()
	UnFocus()
	Draw()
	Close()
	AddHandlerOnly(Key, Handler)
}

type Attributes struct {
	textColor      gocui.Attribute
	textBgColor    gocui.Attribute
	hilightColor   gocui.Attribute
	hilightBgColor gocui.Attribute
}

type Position struct {
	X, Y int
	W, H int
}

type ComponentType int

const (
	TypeInputField ComponentType = iota
	TypeSelect
	TypeButton
	TypeCheckBox
	TypeRadio
)
