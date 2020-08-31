package store

import (
	"io"
)

const (
	// SortNatural use natural order
	SortNatural Sort = iota
	// SortCreatedDesc created newest to oldest
	SortCreatedDesc
	// SortCreatedAsc created oldest to newset
	SortCreatedAsc
	// SortUpdatedDesc updated newest to oldest
	SortUpdatedDesc
	// SortUpdatedAsc updated oldest to newset
	SortUpdatedAsc
)

type (
	// Store is a generic KV Store interface which provides an easier
	// interface to access underlying database. It is mainly used to abstract
	// the database used underneath so we can have a uniform API to use for clients
	Store interface {
		Create(Item) error
		Read(Item) error
		Update(Item) error
		Delete(Item) error
		List(Factory, ListOpt) (int, Items, error)

		io.Closer
	}

	// Item is a generic object which can be used to interact with the store.
	// Users can create their own 'Item' for using the store
	Item interface {
		GetNamespace() string
		GetId() string
	}

	// Serializable interface for the Items to store/retrieve data to/from DB as bytes
	Serializable interface {
		Marshal() ([]byte, error)
		Unmarshal([]byte) error
	}

	// Items stands for list of items. Used by the list operation
	Items []Item

	// Sort is an enum for using different sorting methods on the query
	Sort int

	// ListOpt provides different options for querying the DB
	// Pagination can be used if supported by underlying DB
	ListOpt struct {
		Page    int64
		Limit   int64
		Sort    Sort
		Version int64
		Filters Filter
	}

	SerializedItem interface {
		Item
		Serializable
	}

	Factory interface {
		Factory() SerializedItem
	}

	Filter interface {
		Compare(SerializedItem) bool
	}

	// TimeTracker interface implements basic time tracking functionality
	// for the objects. If Item supports this interface, additional indexes
	// can be maintained to support queries based on this
	TimeTracker interface {
		SetCreated(t int64)
		GetCreated() int64
		SetUpdated(t int64)
		GetUpdated() int64
	}

	// IDSetter interface can be used by the DB to provide new IDs for objects.
	// If Item supports this, when we Create the new item we can set a unique ID
	// based on different DB implementations
	IDSetter interface {
		SetID(string)
	}
)
