package responses

import (
	"github.com/DatGamet/go-discord-wrapper/types/common"
)

type InteractionDataMessageComponent struct {
	CustomID      string               `json:"custom_id"`
	ComponentType common.ComponentType `json:"component_type"`
	Values        *[]interface{}       `json:"values,omitempty"`
	Resolved      *common.ResolvedData `json:"resolved,omitempty"`
}

func (d *InteractionDataMessageComponent) GetType() common.InteractionDataType {
	return common.InteractionDataTypeMessageComponent
}
