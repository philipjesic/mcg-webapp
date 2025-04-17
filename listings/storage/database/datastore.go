package storage

type DataStore interface {
	GetListingByID(id string) (Listing, error)
	GetListings(ids []string) ([]Listing, error)
	InsertListing(l Listing) error
}

type Listing struct {
	ID          string `bson:"id"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
}
