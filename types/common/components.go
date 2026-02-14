package common

type ComponentType int

const (
	ComponentTypeActionRow        ComponentType = 1
	ComponentTypeButton           ComponentType = 2
	ComponentTypeStringSelectMenu ComponentType = 3
	ComponentTypeTextInput        ComponentType = 4
	ComponentTypeUserSelectMenu   ComponentType = 5
	ComponentTypeRoleSelectMenu   ComponentType = 6
	ComponentTypeMentionableMenu  ComponentType = 7
	ComponentTypeChannelSelect    ComponentType = 8
	ComponentTypeSection          ComponentType = 9
	ComponentTypeTextDisplay      ComponentType = 10
	ComponentTypeThumbnail        ComponentType = 11
	ComponentTypeMediaGallery     ComponentType = 12
	ComponentTypeFileDisplay      ComponentType = 13
	ComponentTypeSeparator        ComponentType = 14
	ComponentTypeContainer        ComponentType = 17
	ComponentTypeLabel            ComponentType = 18
	ComponentTypeFileUpload       ComponentType = 19
	ComponentTypeRadioGroup       ComponentType = 21
	ComponentTypeCheckboxGroup    ComponentType = 22
	ComponentTypeCheckbox         ComponentType = 23
)

func (c ComponentType) IsAnySelectMenu() bool {
	return c == ComponentTypeStringSelectMenu ||
		c == ComponentTypeUserSelectMenu ||
		c == ComponentTypeRoleSelectMenu ||
		c == ComponentTypeMentionableMenu ||
		c == ComponentTypeChannelSelect
}

type AnyComponent interface {
	GetType() ComponentType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
