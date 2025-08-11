package model

type Release struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	ArtistCredit []Artist `json:"artist-credit"`
	Date         string   `json:"date,omitempty"`
	Country      string   `json:"country,omitempty"`
}

type Artist struct {
	Name string `json:"name"`
}

type MusicBrainzAlbumResponse struct {
	Created  string    `json:"created"`
	Count    int       `json:"count"`
	Offset   int       `json:"offset"`
	Releases []Release `json:"releases"`
}
