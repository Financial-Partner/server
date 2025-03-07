package entities

type Character struct {
	ID       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	ImageURL string `bson:"image_url" json:"image_url"`
}
