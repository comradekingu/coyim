package gui

import (
	"strings"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFieldTextMulti struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldTextMultiValue

	textView gtki.TextView `gtk-widget:"room-config-text-multi-field-textview"`
}

func newRoomConfigFormTextMulti(f *muc.RoomConfigFormField, value *muc.RoomConfigFieldTextMultiValue) hasRoomConfigFormField {
	field := &roomConfigFieldTextMulti{value: value}
	field.roomConfigFormField = newRoomConfigFormField(f, "MUCRoomConfigFormFieldTextMulti")

	panicOnDevError(field.builder.bindObjects(field))

	tb, _ := g.gtk.TextBufferNew(nil)
	field.textView.SetBuffer(tb)

	v := strings.Join(value.Text(), "\n")
	tb.SetText(v)

	return field
}

// collectFieldValue MUST be called from the UI thread
func (f *roomConfigFieldTextMulti) collectFieldValue() {
	sp := strings.Split(getTextViewText(f.textView), "\n")
	f.value.SetText(sp)
}
