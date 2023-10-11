package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"text/template"
)

const URL = "https://groupietrackers.herokuapp.com/api"

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for the main page
	*/
	tmpl := template.Must(template.ParseFiles("index.html"))
	fmt.Println("r.Method:", r.Method)
	var data JSON
	resp, err := http.Get(URL)
	if err != nil {
		ErrorHandler(w, r, err, "Impossible d'ouvrir l'API à l'adresse suivante:'"+URL+"'", http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&data)
	url := data.Artists
	Art, err := http.Get(url)
	if err != nil {
		ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
	}
	defer Art.Body.Close()
	var Artists []Artists
	// var Artist []Artist
	err = json.NewDecoder(Art.Body).Decode(&Artists)
	if r.Method == "POST" {
		r.ParseForm()
		mCheckBox := GetCheckBoxValue(r)
		startCreationYear := r.FormValue("startCreationYearRange")
		endCreationYear := r.FormValue("endCreationYearRange")
		startFirstAlbumYear := r.FormValue("startFirstAlbumYearRange")
		endFirstAlbumYear := r.FormValue("endFirstAlbumYearRange")
		mLocations := GetLocations(r)
		Artists = SortBand(Artists, mCheckBox, mLocations, Atoi(startFirstAlbumYear), Atoi(endFirstAlbumYear), Atoi(startCreationYear), Atoi(endCreationYear), w, r)
	}
	err = tmpl.Execute(w, Artists)
	if err != nil {
		ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
	}
}
func PageArtistHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for secondary page
	*/
	artistID := r.FormValue("ID")
	fmt.Println("artistsID:", artistID)
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("templates/band.html"))
		var data JSON
		resp, err := http.Get(URL)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&data)
		urlArtist := data.Artists + "/" + artistID
		Art, err := http.Get(urlArtist)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		}
		defer Art.Body.Close()
		var Locations Locations
		urlLocations := data.Locations + "/" + artistID
		Loc, err := http.Get(urlLocations)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		}
		err = json.NewDecoder(Loc.Body).Decode(&Locations)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		}
		defer Loc.Body.Close()
		var artist Artist
		err = json.NewDecoder(Art.Body).Decode(&artist)
		urlRelation := "https://groupietrackers.herokuapp.com/api/relation/" + artistID
		Rel, err := http.Get(urlRelation)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		} else {
			fmt.Println("Tentative d'ouverture:", urlRelation)
		}
		var relation Relation
		err = json.NewDecoder(Rel.Body).Decode(&relation)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		}
		defer Rel.Body.Close()
		var outputBand OutputBand
		outputBand.ID = artist.ID
		outputBand.Image = artist.Image
		outputBand.FirstAlbum = artist.FirstAlbum
		outputBand.Name = artist.Name
		outputBand.Members = artist.Members
		outputBand.CreationDate = artist.CreationDate
		outputBand.Locations = Format_Locations_From_Array(Locations.Lcs)
		outputBand.RelationA = Format_Date(relation.DatesLocations)
		err = tmpl.Execute(w, outputBand)
		if err != nil {
			ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
			return
		}
	case "POST":
		// do nothing
	default:
		http.Redirect(w, r, "400", 400)
	}
}
func Format_Date(m map[string][]string) []string {
	/*
		Utilisé pour formater vers une chaîne les concerts avec les dates
	*/
	a := make([]string, 0)
	var line string
	for k, v := range m {
		line = "\n" + Format_Location_From_String(k) + " : "
		for i, d := range v {
			line += d
			if i < len(v)-1 {
				line += "  " // Ajoutez un espace entre les dates
			}
		}
		a = append(a, line)
	}
	return a
}
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error, message string, errorStatus int) {
	w.WriteHeader(errorStatus)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)
}
func GetCheckBoxValue(r *http.Request) []bool {
	/*
		Called for getting checkbox value into an array,
		this array is returned.
		Theses checkbox elements are accessed by an request object
	*/
	a := make([]bool, 0)
	fmt.Println("firstValue:", r.FormValue("checkboxOneMember"))
	a = append(a, (string(r.FormValue("checkboxOneMember")) == "on"))
	a = append(a, (string(r.FormValue("checkboxTwoMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxThreeMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxFourMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxFiveMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxSixMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxSevenMembers")) == "on"))
	a = append(a, (string(r.FormValue("checkboxEightMembers")) == "on"))
	return a
}
func SortBand(aArtists []Artists, mCheckBox []bool, mLocations map[string]string, startFirstAlbumYear, endFirstAlbumYear, startCreationDateYear, endCreationDateYear int, w http.ResponseWriter, r *http.Request) []Artists {
	/*
		Called for getting all artists sorted.
		All artists that have criteria that matched with those passed as arguments are appended in the array
	*/
	var arrayArtists []Artists
	var isLocationsMatched = false
	arrayMembersValuesPossibiles := GetAllMembersValuePossibles(mCheckBox)
	arrayLocationsModel := GetLocationsAskedByUser(mLocations)
	for _, artist := range aArtists {
		doesFirstAlbumYearMatched := GetIsDateValid(artist.FirstAlbum, startFirstAlbumYear, endFirstAlbumYear)
		doesCreationYearMatched := GetIsDateValid(artist.CreationDate, startCreationDateYear, endCreationDateYear)
		isLocationsMatched = ManageLocation(artist.Locations, arrayLocationsModel, w, r)
		isEnoughMembers := GetMembersMatched(arrayMembersValuesPossibiles, len(artist.Members))
		if len(arrayLocationsModel) == 0 {
			isLocationsMatched = true
		}
		if doesFirstAlbumYearMatched && doesCreationYearMatched && isEnoughMembers && isLocationsMatched {
			arrayArtists = append(arrayArtists, artist)
		}
	}
	return arrayArtists
}
func ToLower(s string) string {
	var result string = ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(rune(char + 32))
		} else {
			result += string(char)
		}
	}
	return result
}
func GetMembersMatched(aCheckBox []int, countMembers int) bool {
	/*
		Called for getting if the number of members correspond to one of theses asked by the user
	*/
	for _, year := range aCheckBox {
		if year == countMembers {
			return true
		}
	}
	return false
}
func GetAllMembersValuePossibles(mCheckBoxIndexMembers []bool) []int {
	/*
		Called for getting all checkbox values from an array passed as argument into array of int
	*/
	a := make([]int, 0)
	var index int = 0
	for _, b := range mCheckBoxIndexMembers {
		if b {
			a = append(a, (index + 1))
		}
		index++
	}
	return a
}
func GetIsDateValid(year interface{}, startYear, endYear int) bool {
	switch year := year.(type) {
	case string:
		yearInt := GetYear(year)
		return yearInt >= startYear && yearInt <= endYear
	case int:
		return year >= startYear && year <= endYear
	default:
		return false
	}
}
func GetYear(year string) int {
	reg := regexp.MustCompile(`(\d{2})-(\d{2})-(?P<year>\d{4})`)
	a := reg.FindStringSubmatch(year)
	return Atoi(a[len(a)-1])
}
func Atoi(s string) int {
	var result int = 0
	var factor int = 1
	for i := len(s) - 1; i >= 0; i-- {
		digit := int(s[i] - '0') // Convertir le caractère en chiffre
		result += digit * factor
		factor *= 10
	}
	return result
}
