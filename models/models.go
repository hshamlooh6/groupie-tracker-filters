package models

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistDetail struct {
	Artist   Artist
	Relation Relation
}

type HomePage struct {
	Artists []Artist
	Query   string

	// Selected filter values, echoed back so the form stays filled in.
	CreationFrom int
	CreationTo   int
	AlbumFrom    int
	AlbumTo      int
	Members      map[int]bool
	Locations    map[string]bool

	// Choices available for the checkbox filters.
	MemberOptions   []int
	LocationOptions []string
}

func (h HomePage) HasMember(n int) bool      { return h.Members[n] }
func (h HomePage) HasLocation(l string) bool { return h.Locations[l] }

type ErrorPage struct {
	StatusCode int
	Message    string
}
