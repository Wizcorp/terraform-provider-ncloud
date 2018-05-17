package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
)

func getServerInfo(client *sdk.Conn, serverID string) (*sdk.ServerInstance, error) {
	readReqParams := new(sdk.RequestGetServerInstanceList)
	readReqParams.ServerInstanceNoList = []string{
		serverID,
	}

	response, err := client.GetServerInstanceList(readReqParams)
	if err != nil {
		return nil, fmt.Errorf("Failed to read server info for server %s: %s", serverID, err)
	}

	if response.TotalRows < 1 {
		return nil, fmt.Errorf("Received no servers in the API response")
	}

	if response.TotalRows == 0 {
		return nil, nil
	}

	return &response.ServerInstanceList[0], nil
}

func getPublicIPInfo(client *sdk.Conn, publicIPID string) (*sdk.PublicIPInstance, error) {
	reqParams := new(sdk.RequestPublicIPInstanceList)
	reqParams.PublicIPInstanceNoList = []string{
		publicIPID,
	}

	response, err := client.GetPublicIPInstanceList(reqParams)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch public IP info %s", err)
	}

	if response.TotalRows == 0 {
		return nil, nil
	}

	return &response.PublicIPInstanceList[0], nil
}

func waitForPublicIPDetach(client *sdk.Conn, publicIPID string) error {
	for {
		IPInfo, err := getPublicIPInfo(client, publicIPID)
		if err != nil {
			return fmt.Errorf("Failed to list server: %s", err)
		}

		if IPInfo.ServerInstance.ServerInstanceNo == "" {
			return nil
		}
	}
}
func waitForServerStatus(client *sdk.Conn, serverID string, status string) error {
	for {
		serverInfo, err := getServerInfo(client, serverID)
		if err != nil {
			return fmt.Errorf("Failed to list server: %s", err)
		}

		if serverInfo == nil || serverInfo.ServerInstanceStatus.Code == status {
			return nil
		}
	}
}
