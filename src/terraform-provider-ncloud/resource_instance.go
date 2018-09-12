package main

import (
	"fmt"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

func disassociatePublicIPInstance(client *sdk.Conn, data *schema.ResourceData, publicIP string) error {
	_, err := client.DisassociatePublicIP(publicIP)
	if err != nil {
		return fmt.Errorf("Failed to disassociate public IP %s", err)
	}

	data.SetPartial("public_ip_instance")
	return nil
}

func associatePublicIPInstance(client *sdk.Conn, data *schema.ResourceData, publicIP string) error {
	associateReqParams := new(sdk.RequestAssociatePublicIP)
	associateReqParams.PublicIPInstanceNo = publicIP
	associateReqParams.ServerInstanceNo = data.Id()

	_, err := client.AssociatePublicIP(associateReqParams)
	if err != nil {
		return fmt.Errorf("Failed to associate public IP %s", err)
	}

	data.SetPartial("public_ip_instance")
	return nil
}

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name",
			},
			"zone_number": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone number (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#zones)",
			},
			"server_product_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Product code (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#servers-server_product_code)",
			},
			"server_image_product_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Server image code (see https://github.com/Wizcorp/terraform-provider-ncloud/blob/master/Services.md#images-server_image_product_code)",
			},
			"login_keyname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "login keyname",
			},
			"termination_protection": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "login keyname",
				Default:     false,
			},
			"user_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "script to run at first boot",
				Default:     false,
			},
			"public_ip_instance": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Public IP ID",
				Default:     "",
			},
			"acg_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ACG",
				Default:     "",
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
	createReqParams.ServerName = data.Get("name").(string)
	createReqParams.ZoneNo = data.Get("zone_number").(string)
	createReqParams.ServerImageProductCode = data.Get("server_image_product_code").(string)
	createReqParams.ServerProductCode = data.Get("server_product_code").(string)
	createReqParams.LoginKeyName = data.Get("login_keyname").(string)
	createReqParams.UserData = data.Get("user_data").(string)
	createReqParams.ServerCreateCount = 1
	createReqParams.IsProtectServerTermination = "false"

	acgID := data.Get("acg_id").(string)
	if acgID != "" {
		createReqParams.AccessControlGroupConfigurationNoList = []string{
			acgID,
		}
	}

	if data.Get("termination_protection").(bool) {
		createReqParams.IsProtectServerTermination = "true"
	}

	response, err := client.CreateServerInstances(createReqParams)
	if err != nil {
		if response != nil && response.ReturnCode == 23006 {
			// Try again in a few seconds
			time.Sleep(1 * time.Second)

			return resourceInstanceCreate(data, meta)
		}

		return fmt.Errorf("Failed to create server: %s", err)
	}

	if response.TotalRows < 1 {
		return fmt.Errorf("Received no servers in the API response")
	}

	serverInfo := response.ServerInstanceList[0]
	data.SetId(serverInfo.ServerInstanceNo)

	data.SetPartial("name")
	data.SetPartial("zone_number")
	data.SetPartial("server_image_product_code")
	data.SetPartial("server_product_code")
	data.SetPartial("login_keyname")
	data.SetPartial("termination_protection")
	data.SetPartial("user_data")

	if acgID != "" {
		data.SetPartial("acg_id")
	}

	listReqParams := new(sdk.RequestGetServerInstanceList)
	listReqParams.ServerInstanceNoList = []string{
		serverInfo.ServerInstanceNo,
	}

	waitForServerStatus(client, data.Id(), "RUN")

	publicIP := data.Get("public_ip_instance").(string)
	if publicIP != "" {
		err = associatePublicIPInstance(client, data, publicIP)
		if err != nil {
			return err
		}
	}

	return resourceInstanceRead(data, meta)
}

func resourceInstanceRead(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)
	serverInfo, err := getServerInfo(client, data.Id())

	if err != nil {
		return err
	}

	data.Set("public_ip", serverInfo.PublicIP)
	data.SetPartial("public_ip")
	data.Set("private_ip", serverInfo.PrivateIP)
	data.SetPartial("private_ip")

	data.Partial(false)
	return nil
}

func resourceInstanceUpdate(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)

	if data.HasChange("server_product_code") {
		changeSpecsReqParams := new(sdk.RequestChangeServerInstanceSpec)
		changeSpecsReqParams.ServerInstanceNo = data.Id()
		changeSpecsReqParams.ServerProductCode = data.Get("server_product_code").(string)

		_, err := client.ChangeServerInstanceSpec(changeSpecsReqParams)
		if err != nil {
			return err
		}

		data.SetPartial("server_product_code")
	}

	if data.HasChange("public_ip_instance") {
		// client.DisassociatePublicIP
		old, new := data.GetChange("public_ip_instance")

		if new.(string) != "" {
			return associatePublicIPInstance(client, data, new.(string))
		} else if old != new {
			return disassociatePublicIPInstance(client, data, old.(string))
		}

		data.SetPartial("public_ip_instance")
	}

	data.Partial(false)

	return nil
}

func resourceInstanceDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)
	publicIP := data.Get("public_ip").(string)
	zoneNo := data.Get("zone_number").(string)

	publicIPReqParams := new(sdk.RequestPublicIPInstanceList)
	publicIPReqParams.IsAssociated = "true"

	publicIPListResponse, err := client.GetPublicIPInstanceList(publicIPReqParams)
	if err != nil {
		return fmt.Errorf("Failed to verify IP association for servers %s: %s", data.Id(), err)
	}

	for _, publicIPInstance := range publicIPListResponse.PublicIPInstanceList {
		if publicIPInstance.PublicIP == publicIP {
			_, err = client.DisassociatePublicIP(publicIPInstance.PublicIPInstanceNo)
			if err != nil {
				return fmt.Errorf("Failed to disassociate IP with ID %s from server %s: %s", publicIP, data.Id(), err)
			}

			waitForPublicIPDetach(client, zoneNo, publicIPInstance.PublicIPInstanceNo)

			break
		}
	}

	stopReqParams := new(sdk.RequestStopServerInstances)
	stopReqParams.ServerInstanceNoList = []string{
		data.Id(),
	}

	stopResponse, err := client.StopServerInstances(stopReqParams)
	if err != nil {
		if stopResponse.ReturnCode != 25041 {
			return fmt.Errorf("Failed to stop servers %s: %s", data.Id(), err)
		}
	}

	waitForServerStatus(client, data.Id(), "NSTOP")

	terminateReqParams := new(sdk.RequestTerminateServerInstances)
	terminateReqParams.ServerInstanceNoList = []string{
		data.Id(),
	}

	_, err = client.TerminateServerInstances(terminateReqParams)
	if err != nil {
		return fmt.Errorf("Failed to terminate servers %s: %s", data.Id(), err)
	}

	waitForServerStatus(client, data.Id(), "TERMT")

	return nil
}
