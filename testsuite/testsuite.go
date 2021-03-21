package store_testsuite

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	store "github.com/SWRMLabs/ss-store"
	"github.com/google/uuid"
)

type Testsuite int

const (
	Basic Testsuite = iota
	Advanced
)

type testStruct struct {
	Namespace string
	Id        string
	RandStr   string
	CreatedAt int64
	UpdatedAt int64
}

type testFactory struct{}

func (f testFactory) Factory() store.SerializedItem {
	return &testStruct{
		Namespace: "StreamSpace",
	}
}

func (t *testStruct) GetNamespace() string { return t.Namespace }

func (t *testStruct) GetId() string { return t.Id }

func (t *testStruct) Marshal() ([]byte, error) { return json.Marshal(t) }

func (t *testStruct) Unmarshal(val []byte) error { return json.Unmarshal(val, t) }

func (t *testStruct) SetCreated(unixTime int64) { t.CreatedAt = unixTime }

func (t *testStruct) SetUpdated(unixTime int64) { t.UpdatedAt = unixTime }

func (t *testStruct) GetCreated() int64 { return t.CreatedAt }

func (t *testStruct) GetUpdated() int64 { return t.UpdatedAt }

func RunTestsuite(t *testing.T, impl store.Store, suite Testsuite) {
	r := &runner{
		t: t,
		s: impl,
	}
	switch suite {
	case Basic:
		r.run(
			TestNilStore,
			TestSimpleCRUD,
			TestSortNaturalLIST,
		)
	case Advanced:
		r.run(
			TestNilStore,
			TestSimpleCRUD,
			TestSortNaturalLIST,
			TestSortCreatedAscLIST,
			TestSortCreatedDscLIST,
			TestSortUpdatedAscLIST,
			TestSortUpdatedDscLIST,
			TestFilterLIST,
		)
	}
}

type implTest func(t *testing.T, s store.Store)

type runner struct {
	t *testing.T
	s store.Store
}

func (r *runner) run(tests ...implTest) {
	for _, t := range tests {
		tName := runtime.FuncForPC(reflect.ValueOf(t).Pointer()).Name()
		tArr := strings.Split(tName, ".")
		fmt.Printf("\t== Running %s\n", tArr[len(tArr)-1])
		t(r.t, r.s)
		fmt.Printf("\t== Passed  %s\n", tArr[len(tArr)-1])
	}
}

func TestNilStore(t *testing.T, s store.Store) {
	if s == nil {
		t.Fatal("Store should not be nil")
	}
}

func TestSimpleCRUD(t *testing.T, s store.Store) {
	// Create new object
	d := &testStruct{
		Namespace: "SS",
		Id:        "04791e92-0b85-11ea-8d71-362b9e155667",
		RandStr:   "totally random",
	}
	err := s.Create(d)
	if err != nil {
		t.Fatalf(err.Error())
	}
	// Read object and verify contents
	nd := &testStruct{
		Namespace: "SS",
		Id:        "04791e92-0b85-11ea-8d71-362b9e155667",
	}
	err = s.Read(nd)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if nd.RandStr != d.RandStr {
		t.Fatalf("Incorrect contents during read")
	}
	// Update object and verify contents again on reading
	d.RandStr = "not totally random"
	err = s.Update(d)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = s.Read(nd)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if nd.RandStr != d.RandStr || nd.RandStr != "not totally random" {
		t.Fatalf("Incorrect contents during read after update")
	}
	// Delete object and make sure its not readable
	err = s.Delete(d)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = s.Read(d)
	if err == nil {
		t.Fatalf("Expected error on read after delete")
	}
}

func TestSortNaturalLIST(t *testing.T, s store.Store) {
	// Create some dummies with StreamSpace namespace
	for i := 0; i < 5; i++ {
		d := testStruct{
			Namespace: "StreamSpace",
			Id:        uuid.New().String(),
			RandStr:   fmt.Sprintf("random %d", i),
		}
		err := s.Create(&d)
		if err != nil {
			t.Fatalf(err.Error())
		}
		// Required for varying timestamps
		<-time.After(time.Second)
	}
	// Create some dummies with Other namespace
	for i := 0; i < 5; i++ {
		d := testStruct{
			Namespace: "Other",
			Id:        uuid.New().String(),
			RandStr:   fmt.Sprintf("random %d", i),
		}
		err := s.Create(&d)
		if err != nil {
			t.Fatalf(err.Error())
		}
		// Required for varying timestamps
		<-time.After(time.Second)
	}
	// Sort '0' is Natural
	opts := store.ListOpt{
		Page:  0,
		Limit: 2,
	}
	// Verify no. of entries with "StreamSpace" NS in list
	var ssEntries, otherEntries = 0, 0
	for i := 0; i < 3; i++ {
		ds, err := s.List(&testFactory{}, opts)
		if err != nil {
			t.Fatalf(err.Error())
		}
		// Although we say limit is '3', we should only get 1 item in the last
		// iteration as there are in all only 10 items
		if (i != 2 && len(ds) != 2) || (i == 2 && len(ds) != 1) {
			t.Fatalf("Unexpected entries in query i: %d Count: %d", i, len(ds))
		}
		for _, v := range ds {
			if v.GetNamespace() == "StreamSpace" {
				ssEntries++
			}
			if v.GetNamespace() == "Other" {
				otherEntries++
			}
		}
		opts.Page++
	}
	if ssEntries != 5 || otherEntries != 0 {
		t.Fatalf("Incorrect entries in List")
	}
}

// Test uses entries from the previous List test
func TestSortCreatedAscLIST(t *testing.T, s store.Store) {
	opts := store.ListOpt{
		Page:  0,
		Limit: 10,
		Sort:  store.SortCreatedAsc,
	}
	ds, err := s.List(&testFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(ds) != 5 {
		t.Fatal("Invalid no of entries", len(ds))
	}
	var created int64 = 0
	for i := 0; i < len(ds); i++ {
		if ds[i].(store.TimeTracker).GetCreated() < created {
			t.Fatalf("Found older created timestamp in ASC List")
		}
		created = ds[i].(store.TimeTracker).GetCreated()
	}
}

// Test uses entries from the previous List test
func TestSortCreatedDscLIST(t *testing.T, s store.Store) {
	opts := store.ListOpt{
		Page:  0,
		Limit: 10,
		Sort:  store.SortCreatedDesc,
	}
	ds, err := s.List(&testFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(ds) != 5 {
		t.Fatal("Invalid no of entries", len(ds))
	}
	var created int64 = 0
	for i := 0; i < len(ds); i++ {
		if ds[i].(store.TimeTracker).GetCreated() > created && i != 0 {
			t.Fatalf("Found older created timestamp in ASC List")
		}
		created = ds[i].(store.TimeTracker).GetCreated()
	}
}

// Test uses entries from the previous List test
func TestSortUpdatedAscLIST(t *testing.T, s store.Store) {
	opts := store.ListOpt{
		Page:  0,
		Limit: 10,
		Sort:  store.SortUpdatedAsc,
	}
	ds, err := s.List(&testFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(ds) != 5 {
		t.Fatal("Invalid no of entries", len(ds))
	}
	var updated int64 = 0
	for i := 0; i < len(ds); i++ {
		if ds[i].(store.TimeTracker).GetUpdated() < updated {
			t.Fatalf("Found older updated timestamp in ASC List")
		}
		updated = ds[i].(store.TimeTracker).GetUpdated()
	}
}

// Test uses entries from the previous List test
func TestSortUpdatedDscLIST(t *testing.T, s store.Store) {
	opts := store.ListOpt{
		Page:  0,
		Limit: 10,
		Sort:  store.SortUpdatedDesc,
	}
	ds, err := s.List(&testFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(ds) != 5 {
		t.Fatal("Invalid no of entries", len(ds))
	}
	var updated int64 = 0
	for i := 0; i < len(ds); i++ {
		if ds[i].(store.TimeTracker).GetUpdated() > updated && i != 0 {
			t.Fatalf("Found older updated timestamp in ASC List")
		}
		updated = ds[i].(store.TimeTracker).GetUpdated()
	}
}

type filterRandStr struct {
	str string
}

func (f filterRandStr) Compare(i store.SerializedItem) bool {
	st, ok := i.(*testStruct)
	if !ok {
		return false
	}
	return st.RandStr == f.str
}

// Test uses entries from the previous List test
func TestFilterLIST(t *testing.T, s store.Store) {
	opts := store.ListOpt{
		Page:  0,
		Limit: 3,
		Filter: filterRandStr{
			str: "random 3",
		},
	}
	ds, err := s.List(&testFactory{}, opts)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(ds) != 1 {
		t.Fatalf("Filter should find only 2 entries Found: %d", len(ds))
	}
	if ds[0].(*testStruct).RandStr != "random 3" {
		t.Fatalf("Invalid filter value in results")
	}
}
