package gameobjects

// Game config constants
const (
	// DB connection URI
	// TODO: make this part of env of course
	GameDbURI string = "mongodb://dungeonmaster:m123123123@localhost:27017/"
	// Max connection pool
	ConnPoolMaxSize uint64 = 100
	// database name
	GameDB string = "tanothDB"
	// Player collection
	ColPlayer string = "players"
)

// Constants defining slot values
const (
	// Item worn for head gear
	SlotHead int = 0

	// Item worn as body armour
	SlotArmour int = 1

	// Primary weapon in main hand
	SlotPrimaryW int = 2

	// Secondary weapon or shield hand
	SlotSecondaryW int = 3

	// Boots
	SlotBoots int = 4

	// Ring
	SlotRing int = 5

	// Amulet
	SlotAmulet int = 6

	// Potion buff slot
	SlotBuff1 int = 7

	// Rune buff slot
	SlotBuff2 int = 8

	// Gloves slot
	SlotGloves int = 9

	// Indicator for max
	SlotMaxVal int = 10
)

// constants for rarity
const (
	RarityCommon   int = 0 // Gray
	RarityRare     int = 1 // Green
	RarityEpic     int = 3 // Blue
	RarityImmortal int = 4 // YEllow - ONLY ONE of this
	RarityMaxVal   int = 5
)

// modifier target
const (
	TargetSelf     int = 0
	TargetOpponent int = 1
	TargetAoe      int = 2
)

// modifier proc types
const (
	Permanent int = 0
)
