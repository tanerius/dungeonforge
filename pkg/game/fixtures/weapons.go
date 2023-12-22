package fixtures

import (
	"github.com/tanerius/dungeonforge/pkg/game/gameobjects"
	"github.com/tanerius/dungeonforge/pkg/game/models"
)

/*
	Some basic weapons here we can extend and store these in a DB at some point

	{
  "weapons": [
    {
      "name": "Whisperwind Blade",
      "description": "A sleek, silver sword that hums softly. Grants exceptional speed and silence, perfect for stealth attacks."
    },
    {
      "name": "Sunfire Hammer",
      "description": "A mighty hammer imbued with the essence of a sun. Emits radiant light and can unleash fiery explosions."
    },
    {
      "name": "Frostbite Dagger",
      "description": "A dagger encased in ice. Freezes surfaces and wounds, immobilizing targets."
    },
    {
      "name": "Starshard Bow",
      "description": "A bow carved from a meteorite. Arrows trail stardust, temporarily blinding enemies."
    },
    {
      "name": "Thunderclap Mace",
      "description": "A mace crackling with electricity. Strikes produce thunderclaps, potentially stunning foes."
    },
    {
      "name": "Eclipse Scythe",
      "description": "A dark-bladed scythe. Casts shadows over enemies, disorienting them in dim light."
    },
    {
      "name": "Serpent's Fang Whip",
      "description": "A serpent-like whip. Constricts and poisons enemies, causing hallucinations and weakness."
    },
    {
      "name": "Voidwalker Staff",
      "description": "A mystical staff for stepping into an ethereal plane, allowing temporary intangibility."
    },
    {
      "name": "Crimson Tide Glaive",
      "description": "A glaive dripping with blood-like substance. Drains life from targets to heal the wielder."
    },
    {
      "name": "Tempest Chakram",
      "description": "A circular, razor-sharp weapon controlling winds. Knocks enemies off balance and redirects projectiles."
    }
  ]
}
*/

const (
	WpnWoodenClub    string = "wooden_club"
	WpnRustySword    string = "rusty_sword"
	WpnIronLongsword string = "iron_longsword"

	WpnReinforcedSling string = "reinforced_sling"
	WpnOakShortbow     string = "oak_shortbow"

	WpnSunfireHammer string = "sunfire_hammer"
)

// Spawn a weapon for a user
//
// weapon = wheapon HrId
// additional string param to specify the userID
func SpawnWeapon(weapon string, a ...string) *models.EquippableItem {

	var userid string = ""

	if len(a) > 0 {
		userid = a[0]
	}

	switch weapon {
	case WpnWoodenClub:
		return &models.EquippableItem{
			Uid: userid,
			Entity: models.Entity{
				HrId: WpnWoodenClub,
				Name: "Wooden Club",
				Desc: "A simple club made from solid wood.",
			},
			MarketInfo: models.MarketInfo{
				IsMarketable: false,
				IsTradable:   false,
				BuyCostGold:  0,
				BuyCostGems:  0,
				SellCostGold: 0,
				SellCostGems: 0,
			},
			Droppable:     false,
			Buffable:      true,
			Slot:          gameobjects.SlotPrimaryW,
			Rarity:        gameobjects.RarityCommon,
			MinLevel:      0,
			StopDropLevel: 1,
			ItemStats: &models.Modifier{
				MyStatMod:      models.NewStats(0, 0, 0, 0, 0),
				EnemyStatMod:   models.NewStats(0, 0, 0, 0, 0),
				MyDamageMod:    models.NewDamage(1, 2, 0, 0),
				EnemyDamageMod: models.NewDamage(0, 0, 0, 0),
			},
			Runes: nil,
		}

		// Rusty sword
	case WpnRustySword:
		return &models.EquippableItem{
			Uid: userid,
			Entity: models.Entity{
				HrId: WpnRustySword,
				Name: "Rusty Sword",
				Desc: "An unimpressive corroded sword, suitable for those relying on brute strength.",
			},
			MarketInfo: models.MarketInfo{
				IsMarketable: false,
				IsTradable:   false,
				BuyCostGold:  30,
				BuyCostGems:  0,
				SellCostGold: 10,
				SellCostGems: 0,
			},
			Droppable:     false,
			Buffable:      true,
			Slot:          gameobjects.SlotPrimaryW,
			Rarity:        gameobjects.RarityCommon,
			MinLevel:      0,
			StopDropLevel: 3,
			ItemStats: &models.Modifier{
				MyStatMod:      models.NewStats(1, 0, 0, 0, 0),
				EnemyStatMod:   models.NewStats(0, 0, 0, 0, 0),
				MyDamageMod:    models.NewDamage(1, 4, 0, 0),
				EnemyDamageMod: models.NewDamage(0, 0, 0, 0),
			},
			Runes: nil,
		}
		// Iron longsword
	case WpnIronLongsword:
		return &models.EquippableItem{
			Uid: userid,
			Entity: models.Entity{
				HrId: WpnIronLongsword,
				Name: "Iron Longsword",
				Desc: "A reliable and sturdy longsword made of iron, perfect for new warriors.",
			},
			MarketInfo: models.MarketInfo{
				IsMarketable: false,
				IsTradable:   false,
				BuyCostGold:  40,
				BuyCostGems:  0,
				SellCostGold: 12,
				SellCostGems: 0,
			},
			Droppable:     false,
			Buffable:      true,
			Slot:          gameobjects.SlotPrimaryW,
			Rarity:        gameobjects.RarityCommon,
			MinLevel:      1,
			StopDropLevel: 3,
			ItemStats: &models.Modifier{
				MyStatMod:      models.NewStats(1, 0, 0, 0, 0),
				EnemyStatMod:   models.NewStats(0, 0, 0, 0, 0),
				MyDamageMod:    models.NewDamage(2, 4, 0, 0),
				EnemyDamageMod: models.NewDamage(0, 0, 0, 0),
			},
			Runes: nil,
		}

		// Reinforced Sling
	case WpnReinforcedSling:
		return &models.EquippableItem{
			Uid: userid,
			Entity: models.Entity{
				HrId: WpnReinforcedSling,
				Name: "Reinforced Sling",
				Desc: "A basic sling with reinforced stitching, good for hurling stones at a distance.",
			},
			MarketInfo: models.MarketInfo{
				IsMarketable: false,
				IsTradable:   false,
				BuyCostGold:  30,
				BuyCostGems:  0,
				SellCostGold: 10,
				SellCostGems: 0,
			},
			Droppable:     false,
			Buffable:      true,
			Slot:          gameobjects.SlotPrimaryW,
			Rarity:        gameobjects.RarityCommon,
			MinLevel:      1,
			StopDropLevel: 3,
			ItemStats: &models.Modifier{
				MyStatMod:      models.NewStats(0, 1, 0, 0, 0),
				EnemyStatMod:   models.NewStats(0, 0, 0, 0, 0),
				MyDamageMod:    models.NewDamage(1, 4, 0, 0),
				EnemyDamageMod: models.NewDamage(0, 0, 0, 0),
			},
			Runes: nil,
		}
		// Oak Shortbow
	case WpnOakShortbow:
		return &models.EquippableItem{
			Uid: userid,
			Entity: models.Entity{
				HrId: WpnOakShortbow,
				Name: "Oak Shortbow",
				Desc: "A basic shortbow crafted from oak wood, ideal for novice archers.",
			},
			MarketInfo: models.MarketInfo{
				IsMarketable: false,
				IsTradable:   false,
				BuyCostGold:  40,
				BuyCostGems:  0,
				SellCostGold: 12,
				SellCostGems: 0,
			},
			Droppable:     false,
			Buffable:      true,
			Slot:          gameobjects.SlotPrimaryW,
			Rarity:        gameobjects.RarityCommon,
			MinLevel:      1,
			StopDropLevel: 3,
			ItemStats: &models.Modifier{
				MyStatMod:      models.NewStats(0, 1, 0, 0, 0),
				EnemyStatMod:   models.NewStats(0, 0, 0, 0, 0),
				MyDamageMod:    models.NewDamage(2, 4, 0, 0),
				EnemyDamageMod: models.NewDamage(0, 0, 0, 0),
			},
			Runes: nil,
		}
	default:
		return nil
	}
}
