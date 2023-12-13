# Fixtrue Descriptions

This folder describes how the items fixtures are defined and what the json fields actually mean.

## main_weapons.json

This file contains fixutres for the main hand weapons. The constant for this is defined as `gameobjects.SlotPrimaryW`. The list includes objects with the following fields:

- **name** - the generic name of the item
- **info** - is an array which contains stats info for the item. Each Element index of the Array represents a stat as described:
  - [0] - Slot identifier
  - [1]