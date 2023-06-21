package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	log "github.com/sirupsen/logrus"
)

func FetchSubnets() ([]string, []string, error) {
	var publicsubnets []string
	var privatesubnets []string

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Errorf("Failed to load AWS config: %v", err)
		return publicsubnets, privatesubnets, err
	}

	client := ec2.NewFromConfig(cfg)

	params := &ec2.DescribeSubnetsInput{}

	response, err := client.DescribeSubnets(context.TODO(), params)
	if err != nil {
		log.Errorf("Got an error retrieving information about your Amazon EC2 subnets: %v", err)
		return publicsubnets, privatesubnets, err
	}

	for _, subnet := range response.Subnets {
		log.Infof("Found subnet with id %s, CIDR block %s, and default for AZ %t\n",
			*subnet.SubnetId, *subnet.CidrBlock, *subnet.DefaultForAz)
	}

	routeTablesResult, err := client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{})
	if err != nil {
		log.Errorf("Error describing route tables: %v", err)
		return publicsubnets, privatesubnets, err
	}

	routeTableMap := make(map[string]*types.RouteTable)
	for _, routeTable := range routeTablesResult.RouteTables {
		routeTableCopy := routeTable
		for _, assoc := range routeTable.Associations {
			if assoc.SubnetId != nil {
				routeTableMap[*assoc.SubnetId] = &routeTableCopy
			}
		}
	}

	for _, subnet := range response.Subnets {
		routeTable := routeTableMap[*subnet.SubnetId]
		if routeTable == nil {
			log.Infof("Subnet ID: %s is private (no route table found)\n", *subnet.SubnetId)
		} else {
			isPublic := false
			for _, route := range routeTable.Routes {
				if route.GatewayId != nil && *route.GatewayId != "local" {
					isPublic = true
					break
				}
			}
			if isPublic {
				log.Infof("Subnet ID: %s is public\n", *subnet.SubnetId)
				publicsubnets = append(publicsubnets, *subnet.SubnetId)
			} else {
				log.Infof("Subnet ID: %s is private\n", *subnet.SubnetId)
				privatesubnets = append(privatesubnets, *subnet.SubnetId)
			}
		}
	}

	return publicsubnets, privatesubnets, nil
}
