package sillyname

import "math/rand"

var adjectives []string = []string{"Big", "Large", "Huge", "Small", "Swole", "Gay", "Lesbian", "Tiny", "Black", "White", "Pink", "Crimson", "Red", "Maroon", "Brown", "Misty", "Rose", "Salmon", "Coral", "Chocolate", "Orange", "Gold", "Ivory", "Yellow", "Olive", "Chartreuse", "Lime", "Green", "Aquamarine", "Turquoise", "Azure", "Cyan", "Teal", "Lavender", "Blue", "Navy", "Violet", "Indigo", "Dark", "Plum", "Magenta", "Purple"}
var nouns []string = []string{"Jerome", "Philipp", "Tobi", "Penis", "Dick", "Boobs", "Ass", "Monkey", "Ape", "Retard", "German", "Gopher", "Marburger", "Dortmunder", "Schwabe"}

func New() string {
	randAdjective := adjectives[rand.Intn(len(adjectives))]
	randNoun := nouns[rand.Intn(len(nouns))]
	return randAdjective + randNoun
}
