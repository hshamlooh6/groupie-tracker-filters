package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"groupie-tracker/models"
)

const (
	artistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	relationsURL = "https://groupietrackers.herokuapp.com/api/relation"
)

func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get(artistsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func GetRelations() ([]models.Relation, error) {
	resp, err := http.Get(relationsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Index []models.Relation `json:"index"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Index, nil
}

func GetArtistDetail(id int) (*models.ArtistDetail, error) {
	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}

	var found *models.Artist
	for _, a := range artists {
		if a.ID == id {
			found = &a
			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("artist with ID %d not found", id)
	}

	relations, err := GetRelations()
	if err != nil {
		return nil, err
	}

	var relation models.Relation
	for _, rel := range relations {
		if rel.ID == id {
			relation = rel
			break
		}
	}

	return &models.ArtistDetail{
		Artist:   *found,
		Relation: relation,
	}, nil
}
