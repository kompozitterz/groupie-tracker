package utils

func SortCreationDate(arrayAllArtists []Artist, from_date, to_date int) []Artist {
	/*
		Return all bands that had started after from_date
	*/
	arrayArtists := make([]Artist, 0)
	for _, artist := range arrayAllArtists {
		if artist.CreationDate > from_date && artist.CreationDate < to_date {
			arrayArtists = append(arrayArtists, artist)
		}
	}
	return arrayArtists
}
