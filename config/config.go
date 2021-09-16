package config

import (
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
)

var Module *sonic.Module

func init() {
	fieldKey := "LOOKUP_pid"

	Module = &sonic.Module{
		Name:     string(product.SiteTarget),
		Accounts: true,
		Fields: []*sonic.ModuleField{
			{
				Validation: `\d+`,
				Type:       sonic.FieldTypeText,
				Label:      "Product ID",
				FieldKey:   &fieldKey,
			},
		},
	}
}
