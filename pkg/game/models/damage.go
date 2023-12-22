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

// Create new stats. Params are:
// min = Minimum Damage
// max = Maximum Damage
// crit = Critical damage multiplier
// chance = Crit proc rate chance
//
// Return *models.Damage
func NewDamage(min, max, crit, chance float32) *Damage {
	return &Damage{
		MinDmg:     min,
		MaxDmg:     max,
		CritVal:    crit,
		CritChance: chance,
	}
}
