package gui

import (
	"errors"

	"github.com/coyim/gotk3adapter/gtk_mock"
	. "gopkg.in/check.v1"
)

type MUCRoomConfigListAddComponentSuite struct{}

var _ = Suite(&MUCRoomConfigListAddComponentSuite{})

type jidEntryMock struct {
	defaultText string
	*gtk_mock.MockEntry
}

func (e *jidEntryMock) GetText() (string, error) {
	return e.defaultText, nil
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_jidCanBeAdded(c *C) {
	la := &mucRoomConfigListAddComponent{
		formItems: []*mucRoomConfigListFormItem{
			{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
			{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}},
			{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
		},
	}

	c.Assert(la.jidCanBeAdded("bla"), Equals, false)
	c.Assert(la.jidCanBeAdded("foo"), Equals, false)
	c.Assert(la.jidCanBeAdded("baz"), Equals, false)
	c.Assert(la.jidCanBeAdded("hi"), Equals, true)
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_removeItemByIndex(c *C) {
	la := &mucRoomConfigListAddComponent{
		contentBox: &gtk_mock.MockBox{},
		formItems: []*mucRoomConfigListFormItem{
			{
				form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}},
				box:  &gtk_mock.MockBox{},
			},
			{
				form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}},
				box:  &gtk_mock.MockBox{},
			},
			{
				form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}},
				box:  &gtk_mock.MockBox{},
			},
		},
	}

	la.removeItemByIndex(0)
	c.Assert(la.formItems, HasLen, 2)
	c.Assert(la.formItems, DeepEquals, []*mucRoomConfigListFormItem{
		{
			form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}},
			box:  &gtk_mock.MockBox{},
		},
		{
			form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}},
			box:  &gtk_mock.MockBox{},
		},
	})

	la.removeItemByIndex(1)
	c.Assert(la.formItems, HasLen, 1)
	c.Assert(la.formItems, DeepEquals, []*mucRoomConfigListFormItem{
		{
			form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}},
			box:  &gtk_mock.MockBox{},
		},
	})
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_forEachForm(c *C) {
	la := &mucRoomConfigListAddComponent{
		formItems: []*mucRoomConfigListFormItem{
			{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
			{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}},
			{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
		},
	}

	result := []string{}
	la.forEachForm(func(form *roomConfigListForm) bool {
		jid, _ := form.jidEntry.GetText()
		result = append(result, jid)
		return true
	})

	c.Assert(result, HasLen, 3)
	c.Assert(result, DeepEquals, []string{"bla", "foo", "baz"})

	result2 := []string{}
	la.forEachForm(func(form *roomConfigListForm) bool {
		jid, _ := form.jidEntry.GetText()
		if jid != "baz" {
			result2 = append(result2, jid)
			return true
		}
		return false
	})

	c.Assert(result2, HasLen, 2)
	c.Assert(result2, DeepEquals, []string{"bla", "foo"})
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_areAllFormsFilled(c *C) {
	la := &mucRoomConfigListAddComponent{
		form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}},
	}

	c.Assert(la.areAllFormsFilled(), Equals, false)

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo@domain.org"}}
	c.Assert(la.areAllFormsFilled(), Equals, true)

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}}
	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}
	c.Assert(la.areAllFormsFilled(), Equals, true)

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}}
	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}
	c.Assert(la.areAllFormsFilled(), Equals, false)

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla@domain.org"}}
	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}
	c.Assert(la.areAllFormsFilled(), Equals, false)
}

type configListAddComponentButtonMock struct {
	*gtk_mock.MockButton
	isEnabled bool
	label     string
}

func (b *configListAddComponentButtonMock) SetSensitive(v bool) {
	b.isEnabled = v
}

func (b *configListAddComponentButtonMock) IsSensitive() bool {
	return b.isEnabled
}

func (b *configListAddComponentButtonMock) SetLabel(l string) {
	b.label = l
}

func (b *configListAddComponentButtonMock) GetLabel() (string, error) {
	return b.label, nil
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_refresh(c *C) {
	la := &mucRoomConfigListAddComponent{
		form:            &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}},
		removeAllButton: &configListAddComponentButtonMock{},
		applyButton:     &configListAddComponentButtonMock{},
	}

	var applyLabel string

	la.refresh()
	c.Assert(la.removeAllButton.IsSensitive(), Equals, false)
	c.Assert(la.applyButton.IsSensitive(), Equals, false)
	applyLabel, _ = la.applyButton.GetLabel()
	c.Assert(applyLabel, Equals, "[localized] Add")

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}
	la.refresh()
	c.Assert(la.removeAllButton.IsSensitive(), Equals, false)
	c.Assert(la.applyButton.IsSensitive(), Equals, true)
	c.Assert(applyLabel, Equals, "[localized] Add")

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}}
	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}
	la.refresh()
	c.Assert(la.removeAllButton.IsSensitive(), Equals, true)
	c.Assert(la.applyButton.IsSensitive(), Equals, true)
	applyLabel, _ = la.applyButton.GetLabel()
	c.Assert(applyLabel, Equals, "[localized] Add all")
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_enableApplyIfConditionsAreMet(c *C) {
	la := &mucRoomConfigListAddComponent{
		form:        &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}},
		applyButton: &configListAddComponentButtonMock{},
	}

	la.enableApplyIfConditionsAreMet()
	c.Assert(la.applyButton.IsSensitive(), Equals, false)

	la.form = &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}
	la.enableApplyIfConditionsAreMet()
	c.Assert(la.applyButton.IsSensitive(), Equals, true)

	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}
	la.enableApplyIfConditionsAreMet()
	c.Assert(la.applyButton.IsSensitive(), Equals, true)

	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}
	la.enableApplyIfConditionsAreMet()
	c.Assert(la.applyButton.IsSensitive(), Equals, false)
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_isFormValid(c *C) {
	la := &mucRoomConfigListAddComponent{}

	var isValid bool
	var err error

	isValid, err = la.isFormValid(&roomConfigListForm{jidEntry: &jidEntryMock{defaultText: ""}})
	c.Assert(isValid, Equals, false)
	c.Assert(err, Not(IsNil))
	c.Assert(err, Equals, errEmptyMemberIdentifier)

	isValid, err = la.isFormValid(&roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "@bla@foo"}})
	c.Assert(isValid, Equals, false)
	c.Assert(err, Not(IsNil))
	c.Assert(err, Equals, errInvalidMemberIdentifier)

	isValid, err = la.isFormValid(&roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla@foo"}})
	c.Assert(isValid, Equals, true)
	c.Assert(err, IsNil)
}

type configListAddComponentDialogMock struct {
	*gtk_mock.MockDialog
	isVisible bool
}

func (d *configListAddComponentDialogMock) IsVisible() bool {
	return d.isVisible
}

func (d *configListAddComponentDialogMock) Show() {
	d.isVisible = true
}

func (d *configListAddComponentDialogMock) Destroy() {
	d.isVisible = false
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_close(c *C) {
	la := &mucRoomConfigListAddComponent{
		dialog: &configListAddComponentDialogMock{isVisible: true},
	}

	c.Assert(la.dialog.IsVisible(), Equals, true)

	la.close()
	c.Assert(la.dialog.IsVisible(), Equals, false)
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_show(c *C) {
	la := &mucRoomConfigListAddComponent{
		dialog: &configListAddComponentDialogMock{},
	}

	c.Assert(la.dialog.IsVisible(), Equals, false)

	la.show()
	c.Assert(la.dialog.IsVisible(), Equals, true)
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_hasItems(c *C) {
	la := &mucRoomConfigListAddComponent{}

	c.Assert(la.hasItems(), Equals, false)

	la.formItems = []*mucRoomConfigListFormItem{
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "foo"}}},
		{form: &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "baz"}}},
	}

	c.Assert(la.hasItems(), Equals, true)
}

func (s *MUCRoomConfigListAddComponentSuite) Test_mucRoomConfigListAddComponent_friendlyErrorMessage(c *C) {
	la := &mucRoomConfigListAddComponent{}
	form := &roomConfigListForm{jidEntry: &jidEntryMock{defaultText: "bla"}}

	c.Assert(la.friendlyErrorMessage(form, nil), Equals, "[localized] Invalid form values.")
	c.Assert(la.friendlyErrorMessage(form, errors.New("bla")), Equals, "[localized] Invalid form values.")
	c.Assert(la.friendlyErrorMessage(form, errEmptyMemberIdentifier), Equals, "[localized] You must enter the account address.")
	c.Assert(la.friendlyErrorMessage(form, errInvalidMemberIdentifier), Equals, "[localized] You must provide a valid account address.")
	c.Assert(la.friendlyErrorMessage(form, errRoomConfigListFormNotFilled), Equals, "[localized] Please, fill in the form fields.")
}
