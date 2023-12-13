package models

type Stats struct {
	// Strength
	// If value range is [-1, 1] then its a percentage value
	// If value range is between (,-1) to (1,) then it's an apsolute value
	// A value of 0 should indicate no change
	Str float32 `bson:"str" json:"str"`

	// Dexterity
	// If value range is [-1, 1] then its a percentage value
	// If value range is between (,-1) to (1,) then it's an apsolute value
	// A value of 0 should indicate no change
	Dex float32 `bson:"dex" json:"dex"`

	// Constitution
	// If value range is [-1, 1] then its a percentage value
	// If value range is between (,-1) to (1,) then it's an apsolute value
	// A value of 0 should indicate no change
	Con float32 `bson:"con" json:"con"`

	// Intelligence
	// If value range is [-1, 1] then its a percentage value
	// If value range is between (,-1) to (1,) then it's an apsolute value
	// A value of 0 should indicate no change
	Int float32 `bson:"int" json:"int"`

	// Armour
	// If value range is [-1, 1] then its a percentage value
	// If value range is between (,-1) to (1,) then it's an apsolute value
	// A value of 0 should indicate no change
	Arm float32 `bson:"arm" json:"arm"`
}

// Create new stats. Params are:
// s = Strength
// d = Dexterity
// c = Constitution
// i = Intelligence
// a = Armour
//
// Return *models.Stats
func NewStats(s, d, c, i, a float32) *Stats {
	return &Stats{
		Str: s,
		Dex: d,
		Con: c,
		Int: i,
		Arm: a,
	}
}
