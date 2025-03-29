package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gacha struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ImgSrc string             `bson:"img_src" json:"img_src"`
}
