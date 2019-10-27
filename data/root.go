package data

// PlexMessage defines the message type passed from the Plex webhook
type PlexMessage struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata struct {
		LibrarySectionType    string `json:"librarySectionType"`
		RatingKey             string `json:"ratingKey"`
		Key                   string `json:"key"`
		ParentRatingKey       string `json:"parentRatingKey"`
		GrandparentRatingKey  string `json:"grandparentRatingKey"`
		GUID                  string `json:"guid"`
		ParentGUID            string `json:"parentGuid"`
		GrandparentGUID       string `json:"grandparentGuid"`
		LibrarySectionTitle   string `json:"librarySectionTitle"`
		LibrarySectionID      int    `json:"librarySectionID"`
		LibrarySectionKey     string `json:"librarySectionKey"`
		Type                  string `json:"type"`
		Title                 string `json:"title"`
		GrandparentKey        string `json:"grandparentKey"`
		ParentKey             string `json:"parentKey"`
		GrandparentTitle      string `json:"grandparentTitle"`
		ParentTitle           string `json:"parentTitle"`
		Summary               string `json:"summary"`
		Index                 int    `json:"index"`
		ParentIndex           int    `json:"parentIndex"`
		ViewOffset            int    `json:"viewOffset"`
		LastViewedAt          int    `json:"lastViewedAt"`
		Year                  int    `json:"year"`
		Thumb                 string `json:"thumb"`
		Art                   string `json:"art"`
		ParentThumb           string `json:"parentThumb"`
		GrandparentThumb      string `json:"grandparentThumb"`
		GrandparentArt        string `json:"grandparentArt"`
		OriginallyAvailableAt string `json:"originallyAvailableAt"`
		AddedAt               int    `json:"addedAt"`
		UpdatedAt             int    `json:"updatedAt"`
	} `json:"Metadata"`
}
