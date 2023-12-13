package models

type Item struct {
	Entity     `bson:",inline" json:",inline"`
	MarketInfo `bson:",inline" json:",inline"`

	// Indicates if an item can be dropped in game
	Droppable bool `bson:"droppable" json:"droppable"`

	// can an item be buffed by a rune
	Buffable bool `bson:"buffable" json:"buffable"`

	// Which slot the item occupies when equipped
	Slot int `bson:"slot" json:"slot"`

	// Rarity of an item
	Rarity int `bson:"rarity" json:"rarity"`

	// Minimum level required to equip item
	MinLevel int `bson:"minlevel" json:"minlevel"`

	// On which level should we stop dropping this item.
	StopDropLevel int `bson:"stopdrop" json:"-"`

	// How the item changes the wielder's stats
	StatsWielder *Stats `bson:"statswielder" json:"statswielder"`

	// How the item changes the enemy's stats when equipped by the wielder.
	// These are for items that not only buff the wielder but also debuff the enemy.
	//
	// Example: To reduce 10% of enemy strength, set StatsEnemy.Str = -0.1
	// Example: To increase enemy's dexterity by 60, set StatsEnemy.Str = 60.0
	StatsEnemy *Stats `bson:"statsenemy" json:"statsenemy"`

	// Damage modifier of the item
	*Damage `bson:"damage" json:"damage"`
}

func (i *Item) Clone() *Item {
	if i.StatsWielder == nil {
		i.StatsWielder = NewStats(0, 0, 0, 0, 0)
	}

	if i.StatsEnemy == nil {
		i.StatsEnemy = NewStats(0, 0, 0, 0, 0)
	}

	if i.Damage == nil {
		i.Damage = &Damage{}
	}

	return &Item{
		Entity:        i.Entity,
		MarketInfo:    i.MarketInfo,
		Droppable:     i.Droppable,
		Buffable:      i.Buffable,
		Slot:          i.Slot,
		Rarity:        i.Rarity,
		MinLevel:      i.MinLevel,
		StopDropLevel: i.StopDropLevel,
		StatsWielder:  NewStats(i.StatsWielder.Str, i.StatsWielder.Dex, i.StatsWielder.Con, i.StatsWielder.Int, i.StatsWielder.Arm),
		StatsEnemy:    NewStats(i.StatsEnemy.Str, i.StatsEnemy.Dex, i.StatsEnemy.Con, i.StatsEnemy.Int, i.StatsEnemy.Arm),
		Damage: &Damage{
			MinDmg:     i.Damage.MinDmg,
			MaxDmg:     i.Damage.MaxDmg,
			CritVal:    i.Damage.CritVal,
			CritChance: i.Damage.CritChance,
		},
	}
}

func (i *Item) GetId() string {
	return i.Id.Hex()
}

func (i *Item) GetHumanId() string {
	return i.HrId
}
