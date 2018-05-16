package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Delete: resourceInstanceDelete,
		Schema: map[string]*schema.Schema{
			"zone_number": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone number (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#zones)",
			},
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
			"login_keyname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "login keyname",
			},
			"termination_protection": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    false,
				Description: "login keyname",
				Default:     false,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInstanceCreate(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)

	createReqParams := new(sdk.RequestCreateServerInstance)
	createReqParams.ZoneNo = data.Get("zone_number").(string)
	createReqParams.ServerImageProductCode = data.Get("server_image_product_code").(string)
	createReqParams.ServerProductCode = data.Get("server_product_code").(string)
	createReqParams.LoginKeyName = data.Get("login_keyname").(string)
	createReqParams.IsProtectServerTermination = data.Get("termination_protection").(bool)
	createReqParams.ServerCreateCount = 1

	response, err := client.CreateServerInstances(createReqParams)
	if err != nil {
		return fmt.Errorf("Failed to create server %s", err)
	}

	if response.TotalRows < 1 {
		return fmt.Errorf("Received no servers in the API response")
	}

	serverInfo := response.ServerInstanceList[0]
	data.SetId(serverInfo.ServerInstanceNo)

	startReqParams := new(sdk.RequestStartServerInstances)
	startReqParams.ServerInstanceNoList = []string{
		data.Id(),
	}

	_, err = client.StartServerInstances(startReqParams)
	if err != nil {
		return fmt.Errorf("Failed to start server %s", err)
	}

	return resourceInstanceRead(data, meta)
}

func resourceInstanceRead(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)

	readReqParams := new(sdk.RequestGetServerInstanceList)

	response, err := client.GetServerInstanceList(readReqParams)
	if err != nil {
		return fmt.Errorf("Failed to read server info %s", err)
	}

	if response.TotalRows < 1 {
		return fmt.Errorf("Received no servers in the API response")
	}

	serverInfo := response.ServerInstanceList[0]

	data.Set("public_ip", serverInfo.PublicIP)
	data.SetPartial("public_ip")
	data.Set("private_ip", serverInfo.PrivateIP)
	data.SetPartial("private_ip")

	data.Partial(false)
	return nil
}

func resourceInstanceDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	terminateReqParams := new(sdk.RequestTerminateServerInstances)
	terminateReqParams.ServerInstanceNoList = []string{
		data.Id(),
	}

	_, err := client.TerminateServerInstances(terminateReqParams)
	if err != nil {
		return fmt.Errorf("Failed to terminate servers %s", err)
	}

	return nil
}
