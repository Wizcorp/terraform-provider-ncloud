package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// NCloudProvider provides the integration between the
// NCloud Go SDK and Terraform
func NCloudProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("NCLOUD_ACCESS_KEY_ID", nil),
				InputDefault: "",
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("NCLOUD_SECRET_ACCESS_KEY", nil),
				InputDefault: "",
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"NCLOUD_REGION",
					"NCLOUD_DEFAULT_REGION",
				}, nil),
				InputDefault: "KO-1",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ncloud_instance":  resourceInstance(),
			"ncloud_login_key": resourceLoginKey(),
			"ncloud_public_ip": resourcePublicIP(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	client := sdk.NewConnection(accessKey, secretKey)

	_, err := client.GetRegionList()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to NCloud APIs, %s", err)
	}

	return client, nil
}
