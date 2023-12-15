package models

//Id primitive.ObjectID `bson:"_id,omitempty" json:"-"`

type Entity struct {
	// A human readable unique ID which will be used for localization later. Index also
	HrId string `bson:"hrid" json:"hrid"`

	// Name of an item in english
	Name string `bson:"name" json:"name"`

	// Item description in english
	Desc string `bson:"desc" json:"desc"`
}
