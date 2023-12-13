package models

type Damage struct {
	// Minimum extra dmg to be added or subtracted
	MinDmg float32 `bson:"mindmg" json:"mindmg"`
	// Maximum extra dmg to be added or subtracted
	MaxDmg float32 `bson:"maxdmg" json:"maxdmg"`
	// Critical dmg multiplier. this is a percentage
	CritVal float32 `bson:"cv" json:"cv"`
	// Percentage of the chance to proc the critical hit
	CritChance float32 `bson:"cc" json:"cc"`
}
