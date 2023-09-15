package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	log "github.com/sirupsen/logrus"
)

func FetchEFSMounts() map[string][]string {
	mounts := make(map[string][]string)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		log.Errorf("Failed to load AWS config: %v", err)
		return nil
	}

	efsClient := efs.NewFromConfig(cfg)

	// List file systems
	fsPaginator := efs.NewDescribeFileSystemsPaginator(efsClient, &efs.DescribeFileSystemsInput{})
	for fsPaginator.HasMorePages() {
		page, err := fsPaginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to get EFS page, %v", err)
		}

		for _, fs := range page.FileSystems {
			// For each file system, get its Name tag
			fsName := "N/A"
			tagsInput := &efs.DescribeTagsInput{
				FileSystemId: fs.FileSystemId,
			}
			tagsResponse, err := efsClient.DescribeTags(ctx, tagsInput)
			if err == nil {
				for _, tag := range tagsResponse.Tags {
					if *tag.Key == "Name" {
						fsName = *tag.Value
					}
				}
			}

			// List mount targets
			mtInput := &efs.DescribeMountTargetsInput{
				FileSystemId: fs.FileSystemId,
			}
			mtResponse, err := efsClient.DescribeMountTargets(ctx, mtInput)
			if err != nil {
				log.WithError(err).Errorf("Failed to describe mount targets for file system %s", *fs.FileSystemId)
				continue
			}

			for _, mt := range mtResponse.MountTargets {
				fmt.Printf("FileSystem Name: %s, FileSystemId: %s, MountTargetId: %s\n", fsName, *fs.FileSystemId, *mt.MountTargetId)
				mounts[fsName] = append(mounts[fsName], *mt.MountTargetId)
			}
		}
	}
	return mounts

}
