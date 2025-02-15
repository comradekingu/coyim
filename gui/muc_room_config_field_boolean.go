package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldBoolean struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldBooleanValue

	grid   gtki.Grid   `gtk-widget:"room-config-field-boolean-grid"`
	toggle gtki.Switch `gtk-widget:"room-config-field-boolean"`
}

func newRoomConfigFormFieldBoolean(f *muc.RoomConfigFormField, value *muc.RoomConfigFieldBooleanValue) hasRoomConfigFormField {
	field := &roomConfigFormFieldBoolean{value: value}
	field.roomConfigFormField = newRoomConfigFormField(f, "MUCRoomConfigFormFieldBoolean")

	panicOnDevError(field.builder.bindObjects(field))

	field.toggle.SetActive(value.Boolean())

	return field
}

// refreshContent MUST NOT be called from the UI thread
func (f *roomConfigFormFieldBoolean) refreshContent() {
	doInUIThread(func() {
		f.grid.SetVAlign(gtki.ALIGN_CENTER)
	})
}

// collectFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldBoolean) collectFieldValue() {
	f.value.SetBoolean(f.toggle.GetActive())
}
