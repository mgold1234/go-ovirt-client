package ovirtclient_test

import (
	"fmt"
	"testing"

	ovirtclient "github.com/ovirt/go-ovirt-client"
)

func TestTagCreation(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)

	tag := assertCanCreateTag(t, helper, fmt.Sprintf("test-%s", helper.GenerateRandomID(5)), "")
	tag2 := assertCanGetTag(t, helper, tag.ID())

	if tag.ID() != tag2.ID() {
		t.Fatalf("IDs of the returned tag don't match.")
	}
}

func TestAddTagToVM(t *testing.T) {
	t.Parallel()
	helper := getHelper(t)
	client := helper.GetClient()
	tagName := fmt.Sprintf("test-%s", helper.GenerateRandomID(5))

	tag := assertCanCreateTag(t, helper, tagName, "")
	vm := assertCanCreateVM(
		t,
		helper,
		tagName,
		nil,
	)
	fetchedVM, err := client.GetVM(vm.ID())
	if err != nil {
		t.Fatal(err)
	}
	if fetchedVM == nil {
		t.Fatal("returned VM is nil")
	}

	err = client.AddTagToVM(vm.ID(), tag.ID())

	if err != nil {
		t.Fatal(err)
	}

	vms, err := client.SearchVMs(ovirtclient.VMSearchParams().WithTag(tagName))
	if err != nil {
		t.Fatalf("Failed to search for VM by Tag (%v)", err)
	}
	if len(vms) != 1 {
		t.Fatalf("Incorrect number of VMs returned (%d)", len(vms))
	}
}

func assertCanGetTag(
	t *testing.T,
	helper ovirtclient.TestHelper,
	tagID string,
) ovirtclient.Tag {
	client := helper.GetClient()
	tag, err := client.GetTag(tagID)

	if err != nil {
		t.Fatalf("Failed to Get Tag (%v)", err)
	}

	return tag
}

func assertCanCreateTag(
	t *testing.T,
	helper ovirtclient.TestHelper,
	name string,
	description string,
) ovirtclient.Tag {
	client := helper.GetClient()
	tag, err := client.CreateTag(
		name,
		description,
	)
	if err != nil {
		t.Fatalf("Failed to create Tag (%v)", err)
	}

	t.Cleanup(
		func() {
			t.Logf("Cleaning up test tag %s...", tag.ID())
			if err := tag.Remove(); err != nil && !ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
				t.Fatalf("Failed to remove test VM %s (%v)", tag.ID(), err)
			}
		},
	)

	return tag
}
