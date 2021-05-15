package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Provider properties
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"zixar_user":     userServer(),
			"zixar_customer": customerServer(),
		},
	}
}
