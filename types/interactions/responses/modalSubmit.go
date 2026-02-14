package responses

import (
	"go-discord-wrapper/types/common"
	"go-discord-wrapper/types/components"
)

type InteractionDataModalSubmit struct {
	CustomID   string                                `json:"custom_id"`
	Resolved   *common.ResolvedData                  `json:"resolved,omitempty"`
	Components *[]components.ComponentLabelComponent `json:"components,omitempty"`
}

func (d *InteractionDataModalSubmit) GetType() common.InteractionDataType {
	return common.InteractionDataTypeModalSubmit
}
