package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"server_image_product_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server image code, e.g",
			},
			"server_product_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product code, e.g. Server Specification under https://www.ncloud.com/charge/calc",
			},
		},
	}
}

func resourceServerCreate(data *schema.ResourceData, meta interface{}) error {
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

func resourceServerRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestGetServerInstanceList)

	_, err := client.GetServerInstanceList(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to create servers %s", err)
	}

	return nil
}

func resourceServerDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestTerminateServerInstances)
	reqParams.ServerInstanceNoList = []string{}

	_, err := client.TerminateServerInstances(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to create servers %s", err)
	}

	return nil
}
