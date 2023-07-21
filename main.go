package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"bufio"
	"os"
)

const (
	MapWidth  = 7
	MapHeight = 7
)

func createMap() [][]bool {
	mappa := make([][]bool, MapWidth)
	for i := range mappa {
		mappa[i] = make([]bool, MapHeight)
	}

	// Generazione casuale degli ostacoli nella mappa
	for i := 0; i < MapWidth; i++ {
		for j := 0; j < MapHeight; j++ {
			if rand.Intn(30) == 0 {
				mappa[i][j] = true
			}
		}
	}

	return mappa
}

type Zombie struct {
	Name     string
	Health   int
	Position Position
}

type Character struct {
	Name     string
	Health   int
	Position Position
}

type Position struct {
	X int
	Y int
}

var (
	mappa      [][]bool
	characters []Character
	zombies    []Zombie
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Creazione della mappa
	mappa = createMap()

	// Creazione dei personaggi
	characters = append(characters, createCharacter("Sylas", 100))
	// characters = append(characters, createCharacter("Elena", 80))

	// Creazione degli zombie
	zombies = append(zombies, createZombie("Walker", 10))
	zombies = append(zombies, createZombie("Runner", 15))

	// Simulazione del gioco
	scanner := bufio.NewScanner(os.Stdin)
	for !isGameOver() {
		// Stampa lo stato del gioco
		printGameState()

		// Leggi l'input dell'utente per il movimento del personaggio
		fmt.Print("Inserisci la direzione del movimento (WASD): ")
		scanner.Scan()
		direction := scanner.Text()
		moveCharacter(&characters[0], direction)

		// Movimento casuale degli zombie
		moveZombies()

		// Controlla gli scontri
		checkCollisions()

		time.Sleep(time.Second) // Pausa di 1 secondo tra ogni iterazione
	}
}

func createCharacter(name string, health int) Character {
	position := getRandomPosition()
	return Character{Name: name, Health: health, Position: position}
}

func createZombie(name string, health int) Zombie {
	position := getRandomPosition()
	return Zombie{Name: name, Health: health, Position: position}
}

func getRandomPosition() Position {
	return Position{
		X: rand.Intn(MapWidth),
		Y: rand.Intn(MapHeight),
	}
}

// func moveCharacter() {
// 	for i := range characters {
// 		direction := getRandomDirection()
// 		newX := characters[i].Position.X + direction.X
// 		newY := characters[i].Position.Y + direction.Y

// 		if isValidPosition(newX, newY) {
// 			characters[i].Position.X = newX
// 			characters[i].Position.Y = newY
// 		}
// 	}
// }

func moveCharacter(character *Character, direction string) {
	switch strings.ToLower(direction) {
	case "w":
		moveCharacterTo(character, character.Position.X-1, character.Position.Y)
	case "a":
		moveCharacterTo(character, character.Position.X, character.Position.Y-1)
	case "s":
		moveCharacterTo(character, character.Position.X+1, character.Position.Y)
	case "d":
		moveCharacterTo(character, character.Position.X, character.Position.Y+1)
	default:
		fmt.Println("Movimento non valido.")
	}
}

func moveCharacterTo(character *Character, x, y int) {
	if isValidPosition(x, y) {
		character.Position.X = x
		character.Position.Y = y
	}
}

func moveZombies() {
	for i := range zombies {
		direction := getRandomDirection()
		newX := zombies[i].Position.X + direction.X
		newY := zombies[i].Position.Y + direction.Y

		if isValidPosition(newX, newY) {
			zombies[i].Position.X = newX
			zombies[i].Position.Y = newY
		} else {
			closestCharacter := getClosestCharacter(&zombies[i])
			if closestCharacter != nil {
				zombies[i].Position.X = closestCharacter.Position.X
				zombies[i].Position.Y = closestCharacter.Position.Y
			}
		}
	}
}

// func isCharacterAdjacent(zombie *Zombie, character *Character) bool {
// 	return (zombie.Position.X == character.Position.X && math.Abs(float64(zombie.Position.Y-character.Position.Y)) == 1) ||
// 		(zombie.Position.Y == character.Position.Y && math.Abs(float64(zombie.Position.X-character.Position.X)) == 1)
// }

func getClosestCharacter(zombie *Zombie) *Character {
	var closestDistance float64 = -1
	var closestChar *Character

	for i := range characters {
		distance := math.Abs(float64(zombie.Position.X-characters[i].Position.X)) +
			math.Abs(float64(zombie.Position.Y-characters[i].Position.Y))
		if closestDistance == -1 || distance < closestDistance {
			closestDistance = distance
			closestChar = &characters[i]
		}
	}
	return closestChar
}

func getRandomDirection() Position {
	directions := []Position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	return directions[rand.Intn(len(directions))]
}

func isValidPosition(x, y int) bool {
	return x >= 0 && x < MapWidth && y >= 0 && y < MapHeight && !mappa[x][y]
}

func checkCollisions() {
	for characterIndex, character := range characters {
		for zombieIndex, zombie := range zombies {
			if character.Position.X == zombie.Position.X && character.Position.Y == zombie.Position.Y {
				fmt.Printf("%s e %s si sono scontrati nella stessa posizione!\n", character.Name, zombie.Name)
				fight(characterIndex, zombieIndex)
			}
		}
	}
}

func fight(characterIndex, zombieIndex int) {
	character := &characters[characterIndex]
	zombie := &zombies[zombieIndex]

	fmt.Printf("%s attacca lo zombie %s!\n", character.Name, zombie.Name)
	damage := rand.Intn(20)
	zombie.Health -= damage
	fmt.Printf("%d danni!\n", damage)

	if zombie.Health <= 0 {
		fmt.Printf("Lo zombie %s è stato sconfitto!\n", zombie.Name)
		zombie.Health = 0
		removeZombie(zombieIndex)
	} else {
		fmt.Printf("Lo zombie %s ha ancora %d di salute e ora attacca!\n", zombie.Name, zombie.Health)
		zombieDamage := 15
		character.Health -= zombieDamage
		fmt.Printf("%d danni!\n", zombieDamage)

		fmt.Printf("%s ha ancora %d di salute.\n", character.Name, character.Health)
		if character.Health <= 0 {
			fmt.Printf("%s è stato sconfitto!\n", character.Name)
			character.Health = 0
			removeCharacter(characterIndex)
		}
	}
}

func removeZombie(index int) {
	zombies = append(zombies[:index], zombies[index+1:]...)

}

func removeCharacter(index int) {
	characters = append(characters[:index], characters[index+1:]...)

}

func isGameOver() bool {
	// Il gioco termina quando tutti gli zombie vengono sconfitti
	for _, zombie := range zombies {
		if zombie.Health > 0 {
			return false
		}
	}
	return true
}

func printGameState() {
	fmt.Println("======= GAME STATE =======")
	for i := 0; i < MapWidth; i++ {
		for j := 0; j < MapHeight; j++ {
			if mappa[i][j] {
				fmt.Print("X ")
			} else {
				hasCharacter := hasCharacterAtPosition(i, j)
				hasZombie := hasZombieAtPosition(i, j)

				if hasCharacter && hasZombie {
					fmt.Print("cz ")
				} else if hasCharacter {
					fmt.Print("c ")
				} else if hasZombie {
					fmt.Print("z ")
				} else {
					fmt.Print(". ")
				}
			}
		}
		fmt.Println()
	}

	fmt.Println()

	for i := range characters {
		character := &characters[i]
		fmt.Printf("%s - Posizione: (%d, %d) - Salute: %d\n", character.Name, character.Position.X, character.Position.Y, character.Health)
	}

	fmt.Println()

	for i := range zombies {
		zombie := &zombies[i]
		fmt.Printf("%s - Posizione: (%d, %d) - Salute: %d\n", zombie.Name, zombie.Position.X, zombie.Position.Y, zombie.Health)
	}

	fmt.Println("==========================")

	if isGameOver() {
		fmt.Println("Tutti gli zombie sono stati sconfitti!")
	}
}

func hasCharacterAtPosition(x, y int) bool {
	for _, character := range characters {
		if character.Position.X == x && character.Position.Y == y {
			return true
		}
	}
	return false
}

func hasZombieAtPosition(x, y int) bool {
	for _, zombie := range zombies {
		if zombie.Position.X == x && zombie.Position.Y == y {
			return true
		}
	}
	return false
}
