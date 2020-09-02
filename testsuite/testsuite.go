package store_test_suite

import (
	"encoding/json"
	"testing"
	"time"

	store "github.com/StreamSpace/ss-store"
	"github.com/google/uuid"
)

type SuccessStruct struct {
	Namespace string
	Id        string
	CreatedAt int64
	UpdatedAt int64
}

type StreamspaceFactory struct {}

func (f StreamspaceFactory) Factory() store.SerializedItem {
	return &SuccessStruct{
		Namespace: "StreamSpace",
	}
}

func (t *SuccessStruct) GetNamespace() string { return t.Namespace }

func (t *SuccessStruct) GetId() string { return t.Id }

func (t *SuccessStruct) Marshal() ([]byte, error) { return json.Marshal(t) }

func (t *SuccessStruct) Unmarshal(val []byte) error { return json.Unmarshal(val, t) }

func (t *SuccessStruct) SetCreated(unixTime int64) { t.CreatedAt = unixTime }

func (t *SuccessStruct) SetUpdated(unixTime int64) { t.UpdatedAt = unixTime }

func (t *SuccessStruct) GetCreated() int64 { return t.CreatedAt }

func (t *SuccessStruct) GetUpdated() int64 { return t.UpdatedAt }

type Tester struct {
	Store store.Store
}

func TestNilStorage(t *testing.T, tester Tester) {
	if tester.Store == nil {
		t.Fatal("Store should not be nil")
	}
}

func TestCreation(t *testing.T, tester Tester) {
	d := SuccessStruct{
		Namespace: "StreamSpace",
		Id:        "04791e92-0b85-11ea-8d71-362b9e155667",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err := tester.Store.Create(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRead(t *testing.T, tester Tester) {
	d := SuccessStruct{
		Namespace: "StreamSpace",
		Id:        "04791e92-0b85-11ea-8d71-362b9e155667",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err := tester.Store.Read(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDelete(t *testing.T, tester Tester) {
	d := SuccessStruct{
		Namespace: "StreamSpace",
		Id:        "04791e92-0b85-11ea-8d71-362b9e155667",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err := tester.Store.Create(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = tester.Store.Update(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = tester.Store.Delete(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestUpdate(t *testing.T, tester Tester) {
	d := SuccessStruct{
		Namespace: "StreamSpace",
		Id:        "04791e92-0b85-11ea-8d71-362b9e155667",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	err := tester.Store.Create(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = tester.Store.Update(&d)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSortNaturalLIST(t *testing.T, tester Tester) {
	// Create some dummies with StreamSpace namespace
	for i := 0; i < 5; i++ {
		d := SuccessStruct{
			Namespace: "StreamSpace",
			Id:        uuid.New().String(),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		err := tester.Store.Create(&d)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	//Create some dummies with Other namespace
	for i := 0; i < 5; i++ {
		d := SuccessStruct{
			Namespace: "Other",
			Id:        uuid.New().String(),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err := tester.Store.Create(&d)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	var sort store.Sort = 0
	opts := store.ListOpt{
		Page:  0,
		Limit: 3,
		Sort:  sort,
	}

	ds, err := tester.Store.List(&StreamspaceFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for i := 0; i < len(ds); i++ {
		if ds[i].GetNamespace() != "StreamSpace" {
			t.Fatalf("Namespace of the %vth element in list dosn't match", i)
		}
	}
}


func TestSortCreatedAscLIST(t *testing.T, tester Tester) {
	var sort store.Sort = 1
	opts := store.ListOpt{
		Page:  0,
		Limit: 3,
		Sort:  sort,
	}

	ds, err := tester.Store.List(&StreamspaceFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for i := 0; i < len(ds); i++ {
		if ds[i].GetNamespace() != "StreamSpace" {
			t.Fatalf("Namespace of the %vth element in list dosn't match", i)
		}
	}
}
func TestSortCreatedDscLIST(t *testing.T, tester Tester) {
	var sort store.Sort = 2
	opts := store.ListOpt{
		Page:  0,
		Limit: 3,
		Sort:  sort,
	}

	ds, err := tester.Store.List(&StreamspaceFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for i := 0; i < len(ds); i++ {
		if ds[i].GetNamespace() != "StreamSpace" {
			t.Fatalf("Namespace of the %vth element in list dosn't match", i)
		}
	}
}
func TestSortUpdatedAscLIST(t *testing.T, tester Tester) {
	var sort store.Sort = 3
	opts := store.ListOpt{
		Page:  0,
		Limit: 3,
		Sort:  sort,
	}

	ds, err := tester.Store.List(&StreamspaceFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for i := 0; i < len(ds); i++ {
		if ds[i].GetNamespace() != "StreamSpace" {
			t.Fatalf("Namespace of the %vth element in list dosn't match", i)
		}
	}
}

func TestSortUpdatedDscLIST(t *testing.T, tester Tester) {
	var sort store.Sort = 4
	opts := store.ListOpt{
		Page:  0,
		Limit: 3,
		Sort:  sort,
	}

	ds, err := tester.Store.List(&StreamspaceFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for i := 0; i < len(ds); i++ {
		if ds[i].GetNamespace() != "StreamSpace" {
			t.Fatalf("Namespace of the %vth element in list dosn't match", i)
		}
	}
}
