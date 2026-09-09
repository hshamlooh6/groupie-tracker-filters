package filters

import (
	"net/url"
	"strconv"
	"strings"

	"groupie-tracker/models"
)

// Selected holds whatever filter values were present on the request.
type Selected struct {
	CreationFrom int
	CreationTo   int
	AlbumFrom    int
	AlbumTo      int
	Members      map[int]bool
	Locations    map[string]bool
}

func (s Selected) HasMember(n int) bool      { return s.Members[n] }
func (s Selected) HasLocation(l string) bool { return s.Locations[l] }

// Parse reads filter values out of the query string. Missing values default
// to "no restriction" (min/max left at 0 means unset for range filters).
func Parse(q url.Values) Selected {
	sel := Selected{
		Members:   map[int]bool{},
		Locations: map[string]bool{},
	}
	sel.CreationFrom, _ = strconv.Atoi(q.Get("creation_from"))
	sel.CreationTo, _ = strconv.Atoi(q.Get("creation_to"))
	sel.AlbumFrom, _ = strconv.Atoi(q.Get("album_from"))
	sel.AlbumTo, _ = strconv.Atoi(q.Get("album_to"))

	for _, v := range q["members"] {
		if n, err := strconv.Atoi(v); err == nil {
			sel.Members[n] = true
		}
	}
	for _, v := range q["locations"] {
		if v != "" {
			sel.Locations[v] = true
		}
	}
	return sel
}

func firstAlbumYear(date string) int {
	parts := strings.Split(date, "-")
	year, _ := strconv.Atoi(parts[len(parts)-1])
	return year
}

// normalizeLocation turns "seattle-washington-usa" into
// "Seattle, Washington, USA".
func NormalizeLocation(raw string) string {
	raw = strings.ReplaceAll(raw, "_", " ")
	parts := strings.Split(raw, "-")
	for i, p := range parts {
		parts[i] = strings.Title(strings.ToLower(strings.TrimSpace(p)))
	}
	return strings.Join(parts, ", ")
}

// locationMatches lets a broader filter value (e.g. "Washington, USA")
// match a more specific concert location (e.g. "Seattle, Washington, USA"),
// per the hint that a city is "part of" its state/country.
func locationMatches(fullLocation, filterValue string) bool {
	full := strings.ToLower(fullLocation)
	val := strings.ToLower(filterValue)
	return full == val || strings.HasSuffix(full, ", "+val) || strings.Contains(full, val)
}

// Matches reports whether an artist satisfies every filter currently set.
func Matches(a models.Artist, rel models.Relation, sel Selected) bool {
	if sel.CreationFrom > 0 && a.CreationDate < sel.CreationFrom {
		return false
	}
	if sel.CreationTo > 0 && a.CreationDate > sel.CreationTo {
		return false
	}

	year := firstAlbumYear(a.FirstAlbum)
	if sel.AlbumFrom > 0 && year < sel.AlbumFrom {
		return false
	}
	if sel.AlbumTo > 0 && year > sel.AlbumTo {
		return false
	}

	if len(sel.Members) > 0 && !sel.Members[len(a.Members)] {
		return false
	}

	if len(sel.Locations) > 0 {
		matched := false
		for rawLocation := range rel.DatesLocations {
			normalized := NormalizeLocation(rawLocation)
			for loc := range sel.Locations {
				if locationMatches(normalized, loc) {
					matched = true
				}
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
