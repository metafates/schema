package bench

import "github.com/metafates/schema/required"

//go:generate schemagen -type DataWithGen

type DataWithGen []struct {
	ID         required.NonZero[string]    `json:"_id"`
	Index      int                         `json:"index"`
	GUID       string                      `json:"guid"`
	IsActive   bool                        `json:"isActive"`
	Balance    string                      `json:"balance"`
	Picture    string                      `json:"picture"`
	Age        required.Positive[int]      `json:"age"`
	EyeColor   string                      `json:"eyeColor"`
	Name       string                      `json:"name"`
	Gender     string                      `json:"gender"`
	Company    string                      `json:"company"`
	Email      string                      `json:"email"`
	Phone      string                      `json:"phone"`
	Address    string                      `json:"address"`
	About      string                      `json:"about"`
	Registered string                      `json:"registered"`
	Latitude   required.Latitude[float64]  `json:"latitude"`
	Longitude  required.Longitude[float64] `json:"longitude"`
	Tags       []string                    `json:"tags"`
	Friends    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"friends"`
	Greeting      string `json:"greeting"`
	FavoriteFruit string `json:"favoriteFruit"`
}

type Data []struct {
	ID         required.NonZero[string]    `json:"_id"`
	Index      int                         `json:"index"`
	GUID       string                      `json:"guid"`
	IsActive   bool                        `json:"isActive"`
	Balance    string                      `json:"balance"`
	Picture    string                      `json:"picture"`
	Age        required.Positive[int]      `json:"age"`
	EyeColor   string                      `json:"eyeColor"`
	Name       string                      `json:"name"`
	Gender     string                      `json:"gender"`
	Company    string                      `json:"company"`
	Email      string                      `json:"email"`
	Phone      string                      `json:"phone"`
	Address    string                      `json:"address"`
	About      string                      `json:"about"`
	Registered string                      `json:"registered"`
	Latitude   required.Latitude[float64]  `json:"latitude"`
	Longitude  required.Longitude[float64] `json:"longitude"`
	Tags       []string                    `json:"tags"`
	Friends    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"friends"`
	Greeting      string `json:"greeting"`
	FavoriteFruit string `json:"favoriteFruit"`
}
