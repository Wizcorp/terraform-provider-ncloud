package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourcePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Delete: resourcePublicIPDelete,
		Schema: map[string]*schema.Schema{
			"server_id": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Product code (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#servers-server_product_code)",
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePublicIPCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)
	data.Partial(true)

	serverID := data.Get("server_id").(string)

	readReqParams := new(sdk.RequestGetServerInstanceList)

	readResponse, err := client.GetServerInstanceList(readReqParams)
	if err != nil {
		return fmt.Errorf("Failed to read server info %s", err)
	}

	if readResponse.TotalRows < 1 {
		return fmt.Errorf("Received no servers in the API response")
	}

	serverInfo := readResponse.ServerInstanceList[0]

	reqParams := new(sdk.RequestCreatePublicIPInstance)
	reqParams.ServerInstanceNo = serverID
	reqParams.RegionNo = serverInfo.Region.RegionNo
	// API doc says we should be allowed to specify th zone
	// reqParams.ZoneNo = serverInfo.Zone.ZoneNo

	response, err := client.CreatePublicIPInstance(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to create public IP %s", err)
	}

	if response.TotalRows < 1 {
		return fmt.Errorf("Received no IPs in the API response")
	}

	ipInfo := response.PublicIPInstanceList[0]
	data.SetId(ipInfo.PublicIPInstanceNo)
	data.SetPartial("server_id")

	return resourcePublicIPRead(data, meta)
}

func resourcePublicIPRead(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)

	ipInfo, err := getPublicIPInfo(client, data.Id())
	if err != nil {
		return fmt.Errorf("Failed to fetch public IP info %s", err)
	}

	data.Set("public_ip", ipInfo.PublicIP)
	data.SetPartial("public_ip")

	return nil
}

func resourcePublicIPDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	disassociateResponse, err := client.DisassociatePublicIP(data.Id())
	if err != nil {
		if disassociateResponse.ReturnCode != 28102 {
			return fmt.Errorf("Failed to disassociate IP with ID %s: %s", data.Id(), err)
		}
	}

	waitForPublicIPDetach(client, data.Id())

	reqParams := new(sdk.RequestDeletePublicIPInstances)
	reqParams.PublicIPInstanceNoList = []string{
		data.Id(),
	}

	_, err = client.DeletePublicIPInstances(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to delete public IP %s", err)
	}

	return nil
}
