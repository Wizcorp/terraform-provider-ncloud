package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourceLoginKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoginKeyCreate,
		Read:   resourceLoginKeyRead,
		Delete: resourceLoginKeyDelete,
		Schema: map[string]*schema.Schema{
			"server_product_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product code (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#servers-server_product_code)",
			},
			"server_image_product_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server image code (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#images-server_image_product_code)",
			},
		},
	}
}

func resourceLoginKeyCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestCreateServerInstance)
	reqParams.ServerImageProductCode = data.Get("server_image_product_code").(string)
	reqParams.ServerProductCode = data.Get("server_product_code").(string)
	reqParams.ServerCreateCount = 1

	_, err := client.CreateServerInstances(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to create servers %s", err)
	}

	return nil
}

func resourceLoginKeyRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestGetServerInstanceList)

	_, err := client.GetServerInstanceList(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to create servers %s", err)
	}

	return nil
}

func resourceLoginKeyDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestTerminateServerInstances)
	reqParams.ServerInstanceNoList = []string{}

	_, err := client.TerminateServerInstances(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to create servers %s", err)
	}

	return nil
}
