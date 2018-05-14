package main

import (
	"fmt"
	"os"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
)

func main() {
	accessKey := os.Getenv("NCLOUD_ACCESS_KEY")
	secretKey := os.Getenv("NCLOUD_SECRET_KEY")

	client := sdk.NewConnection(accessKey, secretKey)

	regionsList, err := client.GetRegionList()

	if err != nil {
		fmt.Printf("Failed to retrieve region list, %v", err)
		os.Exit(1)
	}

	serversReq := new(sdk.RequestGetServerProductList)
	serversReq.ServerImageProductCode = "SPSW0LINUX000046"
	serversList, err := client.GetServerProductList(serversReq)
	if err != nil {
		fmt.Printf("Failed to server types list, %v", err)
		os.Exit(1)
	}

	imagesReq := new(sdk.RequestGetServerImageProductList)
	imagesList, err := client.GetServerImageProductList(imagesReq)
	if err != nil {
		fmt.Printf("Failed to images types list, %v", err)
		os.Exit(1)
	}

	fmt.Print("# Regions, servers and images\n\n")

	fmt.Print("## Regions\n\n")
	fmt.Print("| # | Code       | Name   |\n")
	fmt.Print("| - | ---------- | ------ |\n")
	for _, region := range regionsList.RegionList {
		fmt.Printf("| %s | %s\t | %s\t |\n", region.RegionNo, region.RegionCode, region.RegionName)
	}
	fmt.Print("\n")

	fmt.Print("## Servers (server_product_code)\n\n")
	fmt.Print("| Code             | Description                             |\n")
	fmt.Print("| ---------------- | --------------------------------------- |\n")
	for _, product := range serversList.Product {
		fmt.Printf("| %s | %s\t |\n", product.ProductCode, product.ProductName)
	}
	fmt.Print("\n")

	fmt.Print("## Images (server_image_product_code)\n\n")
	fmt.Print("| Code             | Description                             |\n")
	fmt.Print("| ---------------- | --------------------------------------- |\n")
	for _, product := range imagesList.Product {
		fmt.Printf("| %s | %s |\n", product.ProductCode, product.ProductDescription)
	}
	fmt.Print("\n")
}
