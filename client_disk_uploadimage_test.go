package ovirtclient_test

import (
	"fmt"
	"testing"

	ovirtclient "github.com/ovirt/go-ovirt-client"
)

func assertCanUploadDiskImage(t *testing.T, helper ovirtclient.TestHelper, disk ovirtclient.Disk) {
	fh, stat := getTestImageFile(t)
	defer func() {
		_ = fh.Close()
	}()

	originalSize := disk.ProvisionedSize()
	if originalSize < uint64(stat.Size()) {
		if _, err := disk.Update(
			ovirtclient.UpdateDiskParams().MustWithProvisionedSize(uint64(stat.Size())),
		); err != nil {
			t.Fatalf("Failed to resize disk from %d to %d bytes. (%v)", originalSize, stat.Size(), err)
		}
	}

	client := helper.GetClient()

	if err := client.UploadToDisk(disk.ID(), uint64(stat.Size()), fh); err != nil {
		t.Fatalf("Failed to upload disk image to disk %s. (%v)", disk.ID(), err)
	}
}

func TestImageUploadDiskCreated(t *testing.T) {
	t.Parallel()
	fh, stat := getTestImageFile(t)
	defer func() {
		_ = fh.Close()
	}()

	helper := getHelper(t)
	client := helper.GetClient()

	imageName := fmt.Sprintf("client_test_%s", helper.GenerateRandomID(5))

	uploadResult, err := client.UploadToNewDisk(
		helper.GetStorageDomainID(),
		ovirtclient.ImageFormatRaw,
		uint64(stat.Size()),
		ovirtclient.CreateDiskParams().MustWithSparse(true).MustWithAlias(imageName),
		fh,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to upload image (%w)", err))
	}
	disk, err := client.GetDisk(uploadResult.Disk().ID())
	if err != nil {
		t.Fatal(fmt.Errorf("failed to fetch disk after image upload (%w)", err))
	}
	if err := client.RemoveDisk(disk.ID()); err != nil {
		t.Fatal(err)
	}
}

func TestImageUploadToExistingDisk(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)
	client := helper.GetClient()

	imageName := fmt.Sprintf("client_test_%s", helper.GenerateRandomID(5))

	disk, err := client.CreateDisk(
		helper.GetStorageDomainID(),
		ovirtclient.ImageFormatRaw,
		uint64(512),
		ovirtclient.CreateDiskParams().MustWithSparse(true).MustWithAlias(imageName),
	)
	if disk != nil {
		defer func() {
			_ = disk.Remove()
		}()
	}
	if err != nil {
		t.Fatal(err)
	}

	assertCanUploadDiskImage(t, helper, disk)
}
