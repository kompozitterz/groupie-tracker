package utils

import (
	"encoding/json"
	"net/http"
)

/*
This file is dedicated for methods that will be used for Locations from API see 'struct.go'
*/
func Format_LocationsModel_To_Sort(locationscs []string) []string {
	/*
		Called for formatting element from locationslcs that contains all locations from the front (values from select elements)
		Iterate all elements and format them:
			for example format the following:
				"Mexico, Playa Del Carmen"
			By:
				"Mexico, Playa_Del_Carmen"
	*/
	var isFirstWordCaught = false
	arrayResult := make([]string, 0)
	var state string
	var country string
	for _, loca := range locationscs {
		for _, char := range loca {
			if isFirstWordCaught {
				if char == ' ' && len(state) > 0 {
					state += "_"
				} else if char != ',' {
					state += string(char)
				}
			} else {
				country += string(char)
				if char == ',' {
					isFirstWordCaught = true
				}
			}
		}
		arrayResult = append(arrayResult, country+state)
	}
	return arrayResult
}
func Format_Locations_To_Sort(locationscs []string) []string {
	/*
		Format string for sorting them.
		Format element from locationscs for example,
		in the argument the array could contains an element like that:
			"georgia-usa"
		So this method will format this as:
			"USA, Georgia"
		Theses models are from the API
	*/
	var firstWord string
	var secondWord string
	var isFirstWordCaptured bool = false
	arrayResult := make([]string, 0)
	for _, loca := range locationscs {
		for _, char := range loca {
			if char == '-' {
				isFirstWordCaptured = true
			} else if !isFirstWordCaptured {
				firstWord += string(char)
			} else {
				secondWord += string(char)
			}
		}
		arrayResult = append(arrayResult, secondWord+", "+firstWord)
		isFirstWordCaptured = false
		firstWord = ""
		secondWord = ""
	}
	return arrayResult
}
func Format_Location_From_String(s string) string {
	/*
		Called for formatting locations string output
	*/
	var result string = "   "
	isSpaceBefore := false
	result = "  "
	for index, char := range s {
		if index == 0 {
			result += string(rune(char - 32))
		} else if isSpaceBefore {
			result += string(rune(char - 32))
			isSpaceBefore = false
		} else if char != '_' && char != '-' {
			result += string(char)
		} else if char == '_' || char == '-' {
			result += "   "
			isSpaceBefore = true
		}
	}
	return result
}
func Format_Locations_From_Array(l []string) []string {
	/*
		Called for formatting locations array output
	*/
	var result string = "    "
	isSpaceBefore := false
	for indexE, element := range l {
		result = "   "
		for index, char := range element {
			if index == 0 {
				result += string(rune(char - 32))
			} else if isSpaceBefore {
				result += string(rune(char - 32))
				isSpaceBefore = false
			} else if char != '_' && char != '-' {
				result += string(char)
			} else if char == '_' || char == '-' {
				result += "   "
				isSpaceBefore = true
			}
		}
		l[indexE] = result
	}
	return l
}
func GetLocations(r *http.Request) map[string]string {
	/*
		Called for getting locations field values from front
		Take Args As:
			request witch can help to get value
	*/
	mLocation := make(map[string]string, 0)
	mLocation["europe"] = r.FormValue("locationSelectEurope")
	mLocation["usa"] = r.FormValue("locationSelectAmerica")
	mLocation["asia"] = r.FormValue("locationSelectAsia")
	mLocation["oceania"] = r.FormValue("locationSelectOceania")
	return mLocation
}
func GetLocationsAskedByUser(mLocations map[string]string) []string {
	/*
		Called for getting array of all locations asked by the user,
		Reformat them and return the map as an array
	*/
	result := make([]string, 0)
	var element string = ""
	for _, loca := range mLocations {
		for _, char := range loca {
			if char == ' ' {
				element += ", "
			} else {
				element += string(char)
			}
		}
		if element != "Choose" {
			result = append(result, element)
		}
		element = ""
	}
	return result
}
func ManageLocation(urlLocations string, locationsModels []string, w http.ResponseWriter, r *http.Request) bool {
	/*
		Called for managing matchs between locations of an artist (api) and locations asked by the user (front),
		Return an boolean if it matched or not
	*/
	var Locations Locations
	Loc, err := http.Get(urlLocations)
	if err != nil {
		ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
		//return
	}
	err = json.NewDecoder(Loc.Body).Decode(&Locations)
	if err != nil {
		ErrorHandler(w, r, err, "Impossible d'ouvrir l'API", http.StatusInternalServerError)
		//return
	}
	defer Loc.Body.Close()
	locationsArtists := Format_Locations_To_Sort(Locations.Lcs)                //Format each elements from the locations API
	locationsModelsFormatted := Format_LocationsModel_To_Sort(locationsModels) //Format each elements from the front select element
	for _, locationModel := range locationsModelsFormatted {
		for _, locationArtist := range locationsArtists {
			if ToLower(locationModel) == ToLower(locationArtist) {
				return true
			}
		}
	}
	return false
}
