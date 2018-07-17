package main

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func retryResourcePublicIPCreate(client *sdk.Conn, params *sdk.RequestCreatePublicIPInstance, count int) (*sdk.PublicIPInstance, error) {

	response, err := client.CreatePublicIPInstance(params)
	if err != nil {
		if response.ReturnCode == 10101 && count != 0 {
			// Hack - retry later, we might have more servers
			// and therefore it might succeed
			time.Sleep(5 * time.Second)

			return retryResourcePublicIPCreate(client, params, count-1)
		}

		return nil, fmt.Errorf("Failed to create public IP %s", err)
	}

	if response.TotalRows < 1 {
		return nil, fmt.Errorf("Received no IPs in the API response")
	}

	return &response.PublicIPInstanceList[0], nil
}

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourcePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Delete: resourcePublicIPDelete,
		Schema: map[string]*schema.Schema{
			"region_number": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region number (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#regions)",
			},
			"zone_number": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone number (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#regions)",
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

	reqParams := new(sdk.RequestCreatePublicIPInstance)
	reqParams.ZoneNo = data.Get("zone_number").(string)
	reqParams.RegionNo = data.Get("region_number").(string)

	ipInfo, err := retryResourcePublicIPCreate(client, reqParams, 5)
	if err != nil {
		return err
	}

	data.SetId(ipInfo.PublicIPInstanceNo)
	data.SetPartial("region_number")
	data.SetPartial("zone_number")

	return resourcePublicIPRead(data, meta)
}

func resourcePublicIPRead(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)
	zoneNo := data.Get("zone_number").(string)

	ipInfo, err := getPublicIPInfo(client, zoneNo, data.Id())
	if err != nil {
		return fmt.Errorf("Failed to fetch public IP info %s", err)
	}

	data.Set("public_ip", ipInfo.PublicIP)
	data.SetPartial("public_ip")

	return nil
}

func resourcePublicIPDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)
	zoneNo := data.Get("zone_number").(string)

	disassociateResponse, err := client.DisassociatePublicIP(data.Id())
	if err != nil {
		if disassociateResponse.ReturnCode != 28102 {
			return fmt.Errorf("Failed to disassociate IP with ID %s: %s", data.Id(), err)
		}
	}

	waitForPublicIPDetach(client, zoneNo, data.Id())

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
