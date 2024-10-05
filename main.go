package main

import (
	"fmt"
	"strings"
)

const (
	kitchen  = "кухня"
	corridor = "коридор"
	room     = "комната"
	street   = "улица"
)

type Location struct {
	name  string
	desc  string
	items map[string]string
	ways  map[string]string
}

type Character struct {
	curRoom  *Location
	inv      map[string]string
	backpack bool
	door     bool
}

var (
	character *Character
	locations map[string]*Location
)

func initGame() {
	locations = make(map[string]*Location)
	locations[kitchen] = &Location{
		name: kitchen,
		desc: "кухня, ничего интересного. можно пройти - коридор",
		items: map[string]string{
			"чай": "На столе стоит чашка чая.",
		},
		ways: map[string]string{
			corridor: corridor,
		},
	}
	locations[corridor] = &Location{
		name:  corridor,
		desc:  "ничего интересного. можно пройти - кухня, комната, улица",
		items: make(map[string]string),
		ways: map[string]string{
			kitchen: kitchen,
			room:    room,
			street:  street,
		},
	}
	locations[room] = &Location{
		name: room,
		desc: "ты в своей комнате. можно пройти - коридор",
		items: map[string]string{
			"ключи":     "ключи лежат на столе.",
			"конспекты": "конспекты на столе.",
			"рюкзак":    "рюкзак висит на стуле.",
		},
		ways: map[string]string{
			corridor: corridor,
		},
	}
	locations[street] = &Location{
		name:  street,
		desc:  "на улице весна. можно пройти - домой",
		items: make(map[string]string),
		ways: map[string]string{
			"домой": corridor,
		},
	}
	character = &Character{
		curRoom:  locations[kitchen],
		inv:      make(map[string]string),
		backpack: false,
		door:     false,
	}
}

func handleLookAround() string {
	if character.curRoom.name == room {
		if len(character.curRoom.items) == 0 {
			return "пустая комната. можно пройти - коридор"
		}
		var itemList []string
		if _, exists := character.curRoom.items["ключи"]; exists {
			itemList = append(itemList, "ключи")
		}
		if _, exists := character.curRoom.items["конспекты"]; exists {
			itemList = append(itemList, "конспекты")
		}
		if _, exists := character.curRoom.items["рюкзак"]; exists {
			return fmt.Sprintf("на столе: %s, на стуле: рюкзак. можно пройти - коридор", strings.Join(itemList, ", "))
		}
		return fmt.Sprintf("на столе: %s. можно пройти - коридор", strings.Join(itemList, ", "))
	}
	if character.curRoom.name == kitchen && !character.backpack {
		return "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
	}
	if character.curRoom.name == kitchen && character.backpack {
		return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
	}
	return character.curRoom.desc
}

func handleMove(parts []string) string {
	if len(parts) < 2 {
		return "куда идти?"
	}
	if parts[1] == street && !character.door {
		return "дверь закрыта"
	}
	nextRoomName, exists := character.curRoom.ways[parts[1]]
	if !exists {
		return fmt.Sprintf("нет пути в %s", parts[1])
	}
	character.curRoom = locations[nextRoomName]
	return character.curRoom.desc
}

func handleTakeItem(parts []string) string {
	if len(parts) < 2 {
		return "что взять?"
	}
	if !character.backpack {
		return "некуда класть"
	}
	item, exists := character.curRoom.items[parts[1]]
	if !exists {
		return "нет такого"
	}
	character.inv[parts[1]] = item
	delete(character.curRoom.items, parts[1])
	if character.curRoom.name == kitchen && len(character.curRoom.items) == 0 {
		character.curRoom.desc = "кухня, ничего интересного. можно пройти - коридор"
	}
	return fmt.Sprintf("предмет добавлен в инвентарь: %s", parts[1])
}

func handleWearItem(parts []string) string {
	if len(parts) < 2 {
		return "что надеть?"
	}
	if parts[1] == "рюкзак" {
		if _, exists := character.curRoom.items["рюкзак"]; exists {
			character.backpack = true
			delete(character.curRoom.items, "рюкзак")
			return "вы надели: рюкзак"
		}
		return "нет такого"
	}
	return "неизвестная команда"
}

func handleApplyItem(parts []string) string {
	if len(parts) < 3 {
		return "не к чему применить"
	}
	item, inInv := character.inv[parts[1]]
	if !inInv {
		return fmt.Sprintf("нет предмета в инвентаре - %s", parts[1])
	}
	if parts[2] == "дверь" && item == "ключи лежат на столе." && character.curRoom.name == corridor {
		character.door = true
		return "дверь открыта"
	}
	return "не к чему применить"
}

func handleCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "неизвестная команда"
	}
	switch parts[0] {
	case "осмотреться":
		return handleLookAround()
	case "идти":
		return handleMove(parts)
	case "взять":
		return handleTakeItem(parts)
	case "надеть":
		return handleWearItem(parts)
	case "применить":
		return handleApplyItem(parts)
	default:
		return "неизвестная команда"
	}
}

func main() {
	initGame()
	fmt.Println("Начало игры!")
	fmt.Println(character.curRoom.desc)
}
