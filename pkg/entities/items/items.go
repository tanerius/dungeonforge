package items

// ItemType type
type ItemType string

// ItemTypes is a slice of ItemType
type ItemTypes []ItemType

const (
	ItemTypeCurrency    ItemType = "currency"
	ItemTypeConsumable           = "consumable"
	ItemTypeArmor                = "armor"
	ItemTypeWeapon               = "weapon"
	ItemTypeCollectible          = "collectible"
	ItemTypeQuest                = "quest"

	ItemTypeCraftingMaterial = "crafting_material"
)

//ItemSubType type
type ItemSubType string

//ItemSubTypes ...
type ItemSubTypes []ItemSubType

const (
	// weapons
	ItemSubTypeSword        ItemSubType = "sword"
	ItemSubTypeTwoHandSword ItemSubType = "twohandsword"
	ItemSubTypeAxe                      = "axe"
	ItemSubTypeSpear                    = "spear"

	// shields
	ItemSubTypeShield = "shield"
)

//ItemSlot type
type ItemSlot string

//ItemSlots type
type ItemSlots []ItemSlot

const (
	ItemSlotInventory ItemSlot = "inventory"
	ItemSlotContainer          = "container"
	ItemSlotPurse              = "purse"
	ItemSlotHead               = "head"
	ItemSlotChest              = "chest"
	ItemSlotLegs               = "legs"
	ItemSlotBoots              = "boots"
	ItemSlotNeck               = "neck"
	ItemSlotRing1              = "ring1"
	ItemSlotRing2              = "ring2"
	ItemSlotHands              = "hands"
	ItemSlotMainHand           = "main_hand"
	ItemSlotOffHand            = "off_hand"
	ItemSlotBuff               = "buff"
)

//ItemQuality ...
type ItemQuality string

//ItemQualities type
type ItemQualities []ItemQuality

const (
	ItemQualityNormal    ItemQuality = "normal"
	ItemQualityMagic                 = "magic"
	ItemQualityRare                  = "rare"
	ItemQualityLegendary             = "legendary"
	ItemQualityMythic                = "mythic"
)
