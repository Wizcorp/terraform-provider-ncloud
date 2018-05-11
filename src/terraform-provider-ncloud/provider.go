package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// NCloud Provider
func NCloudProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"server": resourceServer(),
		},
	}
}
