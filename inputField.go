package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Margin struct {
	Top  int
	Left int
}

type InputField struct {
	*gocui.Gui
	Label *Label
	Field *Field
}

type Label struct {
	Text      string
	Width     int
	DrawFrame bool
	*Position
	*Attributes
	Margin *Margin
}

type Field struct {
	Text      string
	Width     int
	DrawFrame bool
	*Position
	*Attributes
	Handlers Handlers
	Margin   *Margin
	*Validate
}

var labelPrefix = "label"

// NewInputField new input label and field
func NewInputField(gui *gocui.Gui, labelText string, x, y, labelWidth, fieldWidth int) *InputField {
	// label psition
	lp := &Position{
		x,
		y,
		x + labelWidth + 1,
		y + 2,
	}

	// field position
	fp := &Position{
		lp.W,
		lp.Y,
		lp.W + fieldWidth,
		lp.H,
	}

	// new label
	label := &Label{
		Text:     labelText,
		Width:    labelWidth,
		Position: lp,
		Attributes: &Attributes{
			TextColor:   gocui.ColorYellow,
			TextBgColor: gocui.ColorBlack,
		},
		DrawFrame: false,
		Margin: &Margin{
			Top:  0,
			Left: 0,
		},
	}

	// new field
	field := &Field{
		Width:    fieldWidth,
		Position: fp,
		Attributes: &Attributes{
			TextColor:   gocui.ColorBlack,
			TextBgColor: gocui.ColorCyan,
			FgColor:     gocui.ColorBlack,
			BgColor:     gocui.ColorCyan,
		},
		Handlers:  make(Handlers),
		DrawFrame: false,
		Margin: &Margin{
			Top:  0,
			Left: 0,
		},
		Validate: &Validate{
			Gui:       gui,
			Name:      label.Text + "errMsg",
			Validator: func(text string) bool { return true },
			Position: &Position{
				X: fp.W,
				Y: fp.Y,
				W: fp.X + fp.W,
				H: fp.H,
			},
		},
	}

	// new input field
	i := &InputField{
		Gui:   gui,
		Label: label,
		Field: field,
	}

	return i
}

// AddFieldTextAttribute add field colors
func (i *InputField) AddFieldAttribute(textColor, textBgColor, fgColor, bgColor gocui.Attribute) *InputField {
	i.Field.Attributes = &Attributes{
		TextColor:   textColor,
		TextBgColor: textBgColor,
		FgColor:     fgColor,
		BgColor:     bgColor,
	}
	return i
}

// AddLabelAttribute add label colors
func (i *InputField) AddLabelAttribute(textColor, textBgColor gocui.Attribute) *InputField {
	i.Label.Attributes = &Attributes{
		TextColor:   textColor,
		TextBgColor: textBgColor,
	}

	return i
}

// AddHandler add keybinding
func (i *InputField) AddHandler(key Key, handler Handler) *InputField {
	i.Field.Handlers[key] = handler
	return i
}

// AddMarginTop add margin top
func (i *InputField) AddMarginTop(top int) *InputField {
	i.Label.Margin.Top += top
	i.Field.Margin.Top += top
	return i
}

// AddMarginLeft add margin left
func (i *InputField) AddMarginLeft(left int) *InputField {
	i.Label.Margin.Left += left
	i.Field.Margin.Left += left
	return i
}

// AddValidator add input validator
func (i *InputField) AddValidator(errMsg string, validator Validator) *InputField {
	v := i.Field.Validate
	v.ErrMsg = errMsg
	v.Validator = validator
	v.W += len(errMsg)
	return i
}

// SetLabelBorder draw label border
func (i *InputField) SetLabelBorder() *InputField {
	i.Label.DrawFrame = true
	return i
}

// SetFieldBorder draw field border
func (i *InputField) SetFieldBorder() *InputField {
	i.Field.DrawFrame = true
	return i
}

// Edit input field editor
func (i *InputField) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(+1, 0, false)
	}

	// get field text
	text, _ := v.Line(0)
	i.Field.Text = text

	// validate input
	i.Field.IsValid = i.Field.Validator(text)

	if !i.Field.IsValid {
		i.Field.DispValidateMsg()
	} else {
		i.Field.CloseValidateMsg()
	}
}

// GetFieldText get input field text
func (i *InputField) GetFieldText() string {
	return i.Field.Text
}

// IsValid valid field data will be return true
func (i *InputField) IsValid() bool {
	return i.Field.Validate.IsValid
}

// Draw draw label and field
func (i *InputField) Draw() *InputField {
	// draw label
	x, y, w, h := i.addMargin(i.Label)
	if v, err := i.Gui.SetView(labelPrefix+i.Label.Text, x, y, w, h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = i.Label.DrawFrame

		v.FgColor = i.Label.TextColor | gocui.AttrBold
		v.BgColor = i.Label.TextBgColor

		fmt.Fprint(v, i.Label.Text+": ")
	}

	// draw input
	x, y, w, h = i.addMargin(i.Field)
	if v, err := i.Gui.SetView(i.Label.Text, x, y, w, h); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = i.Field.DrawFrame
		v.Highlight = true

		v.SelBgColor = i.Field.BgColor
		v.SelFgColor = i.Field.FgColor

		v.FgColor = i.Field.TextColor
		v.BgColor = i.Field.TextBgColor

		v.Editable = true
		v.Editor = i
	}

	// set keybindings
	if i.Field.Handlers != nil {
		for key, handler := range i.Field.Handlers {
			if err := i.Gui.SetKeybinding(i.Label.Text, key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}

	// focus input field
	i.Gui.SetCurrentView(i.Label.Text)
	i.Gui.SetViewOnTop(i.Label.Text)

	return i
}

func (i *InputField) addMargin(view interface{}) (int, int, int, int) {
	switch v := view.(type) {
	case *Field:
		p := v.Position
		m := v.Margin
		return p.X + m.Left, p.Y + m.Top, p.W + m.Left, p.H + m.Top
	case *Label:
		p := v.Position
		m := v.Margin
		return p.X + m.Left, p.Y + m.Top, p.W + m.Left, p.H + m.Top
	default:
		panic("Unkown type")
	}
}

// Close close input field
func (i *InputField) Close() {
	views := []string{
		i.Label.Text,
		labelPrefix + i.Label.Text,
	}

	for _, v := range views {
		if err := i.DeleteView(v); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
	}

	if i.Field.Handlers != nil {
		i.DeleteKeybindings(i.Label.Text)
	}
}