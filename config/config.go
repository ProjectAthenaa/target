package config

import (
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
)

var Module *sonic.Module

func init() {
	fieldKey := "LOOKUP_link"

	Module = &sonic.Module{
		Name:     string(product.SiteTarget),
		Accounts: true,
		Fields: []*sonic.ModuleField{
			{
				Validation: "https:\\/\\/www\\.target\\.com\\/p\\/.*\\/-\\/A-\\d+",
				Type:       sonic.FieldTypeText,
				Label:      "Product Link",
				FieldKey:   &fieldKey,
			},
		},
	}
}
