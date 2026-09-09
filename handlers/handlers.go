package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"groupie-tracker/api"
	"groupie-tracker/filters"
	"groupie-tracker/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "Page not found")
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		fmt.Println("API ERROR:", err)
		renderError(w, http.StatusInternalServerError, "Could not load artists")
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		fmt.Println("API ERROR:", err)
		renderError(w, http.StatusInternalServerError, "Could not load relations")
		return
	}

	relByID := map[int]models.Relation{}
	for _, rel := range relations {
		relByID[rel.ID] = rel
	}

	sel := filters.Parse(r.URL.Query())

	var filtered []models.Artist
	for _, a := range artists {
		if filters.Matches(a, relByID[a.ID], sel) {
			filtered = append(filtered, a)
		}
	}

	// Build the checkbox options from the full (unfiltered) dataset.
	memberSet := map[int]bool{}
	for _, a := range artists {
		memberSet[len(a.Members)] = true
	}
	var memberOptions []int
	for n := range memberSet {
		memberOptions = append(memberOptions, n)
	}
	sort.Ints(memberOptions)

	locationSet := map[string]bool{}
	for _, rel := range relations {
		for rawLocation := range rel.DatesLocations {
			locationSet[filters.NormalizeLocation(rawLocation)] = true
		}
	}
	var locationOptions []string
	for loc := range locationSet {
		locationOptions = append(locationOptions, loc)
	}
	sort.Strings(locationOptions)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Template error")
		return
	}

	tmpl.Execute(w, models.HomePage{
		Artists:         filtered,
		Query:           "",
		CreationFrom:    sel.CreationFrom,
		CreationTo:      sel.CreationTo,
		AlbumFrom:       sel.AlbumFrom,
		AlbumTo:         sel.AlbumTo,
		Members:         sel.Members,
		Locations:       sel.Locations,
		MemberOptions:   memberOptions,
		LocationOptions: locationOptions,
	})

}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		renderError(w, http.StatusBadRequest, "Missing artist ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	artist, err := api.GetArtistDetail(id)
	if err != nil {
		renderError(w, http.StatusNotFound, "Artist not found")
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Template error")
		return
	}
	tmpl.Execute(w, artist)

}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	results, err := api.GetArtists()
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Could not search")
		return
	}

	var filtered []models.Artist
	for _, artist := range results {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			filtered = append(filtered, artist)
		}
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Template error")
		return
	}

	tmpl.Execute(w, models.HomePage{Artists: filtered, Query: query})

}

func renderError(w http.ResponseWriter, code int, msg string) {

	w.WriteHeader(code)

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, msg, code)
		return
	}

	data := models.ErrorPage{
		StatusCode: code,
		Message:    msg,
	}

	tmpl.Execute(w, data)
}
