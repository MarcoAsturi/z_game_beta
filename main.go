package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
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
	Weapon   Weapon
}

type Position struct {
	X int
	Y int
}

type ZombiePosition struct {
	Position    // Includi i campi X e Y della struct Position
	ZombieIndex int
}

type Weapon struct {
	Name   string
	Damage int
	Range  int
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
	characters = append(characters, createCharacter("Elsa", 80))

	// Creazione degli zombie
	zombies = append(zombies, createZombie("Walker", 10))
	zombies = append(zombies, createZombie("Runner", 15))

	// WaitGroup per sincronizzare il completamento del movimento degli zombie
	var wg sync.WaitGroup

	// Simulazione del gioco
	for !isGameOver() {
		// Stampa lo stato del gioco
		printGameState()

		// Movimento dei personaggi
		for i := 0; i < len(characters); i++ {
			currentCharacter := &characters[i]

			// Stampiamo il messaggio per il personaggio corrente
			fmt.Printf("%s, fai la tua mossa\n", currentCharacter.Name)

			// Leggi l'input dell'utente per il movimento del personaggio
			fmt.Print("Inserisci la direzione del movimento (WASD): ")
			var direction string
			fmt.Scan(&direction)

			// Fai muovere il personaggio corrente
			moveCharacter(i, direction)
		}

		// Movimento degli zombie in modo concorrente
		for i := range zombies {
			wg.Add(1)
			go moveZombieConcurrently(&zombies[i], &wg)
		}

		// Aspetta che tutti gli zombie abbiano completato il movimento
		wg.Wait()

		// Controlla gli scontri
		checkCollisions()

		time.Sleep(time.Second / 2) // Pausa di 1/2 secondo tra ogni iterazione
	}
}

func moveZombieConcurrently(zombie *Zombie, wg *sync.WaitGroup) {
	defer wg.Done()

	// Trova personaggio più vicino
	closestCharacter := getClosestCharacter(zombie)

	// Caso 1: Lo zombie si trova nella stessa zona di un personaggio, resta fermo
	if closestCharacter != nil && isCharacterInSameZoneAsZombie(zombie, closestCharacter) {
		return
	} else if closestCharacter != nil && isCharacterAdjacentToZombie(zombie, closestCharacter) {
		// Caso 2: Lo zombie si trova in una posizione adiacente al personaggio, si muove nella stessa zona del personaggio
		newX := closestCharacter.Position.X
		newY := closestCharacter.Position.Y

		if isValidPosition(newX, newY) {
			zombie.Position.X = newX
			zombie.Position.Y = newY
		}
		return
	} else {
		// Caso 3: In tutti gli altri casi, lo zombie esegue un movimento randomico
		direction := getRandomDirection()
		newX := zombie.Position.X + direction.X
		newY := zombie.Position.Y + direction.Y

		if isValidPosition(newX, newY) {
			zombie.Position.X = newX
			zombie.Position.Y = newY
		}
	}
}

func isCharacterInSameZoneAsZombie(zombie *Zombie, character *Character) bool {
	return (zombie.Position.X/3 == character.Position.X/3) && (zombie.Position.Y/3 == character.Position.Y/3)
}

func isCharacterAtPosition(x, y int) bool {
	for _, character := range characters {
		if character.Position.X == x && character.Position.Y == y {
			return true
		}
	}
	return false
}

func moveZombieTowardsCharacter(zombie *Zombie, character *Character) {
	deltaX := character.Position.X - zombie.Position.X
	deltaY := character.Position.Y - zombie.Position.Y

	if deltaX != 0 {
		zombie.Position.X += deltaX / int(math.Abs(float64(deltaX)))
	}
	if deltaY != 0 {
		zombie.Position.Y += deltaY / int(math.Abs(float64(deltaY)))
	}
}

func isCharacterAdjacentToZombie(zombie *Zombie, character *Character) bool {
	return (zombie.Position.X == character.Position.X && math.Abs(float64(zombie.Position.Y-character.Position.Y)) == 1) ||
		(zombie.Position.Y == character.Position.Y && math.Abs(float64(zombie.Position.X-character.Position.X)) == 1)
}

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

func getRandomDirection() ZombiePosition {
	directions := []ZombiePosition{
		{Position: Position{X: 0, Y: 1}},
		{Position: Position{X: 0, Y: -1}},
		{Position: Position{X: 1, Y: 0}},
		{Position: Position{X: -1, Y: 0}},
	}

	randomIndex := rand.Intn(len(directions))
	return directions[randomIndex]
}

func isValidPosition(x, y int) bool {
	return x >= 0 && x < MapWidth && y >= 0 && y < MapHeight && !mappa[x][y]
}

func createCharacter(name string, health int) Character {
	position := getRandomPosition()
	weapon := Weapon{
		Name:   "Sword",
		Damage: 20,
		Range:  0,
	}
	return Character{Name: name, Health: health, Position: position, Weapon: weapon}
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

func moveCharacter(characterIndex int, direction string) {
	currentCharacter := &characters[characterIndex]

	switch strings.ToLower(direction) {
	case "w":
		moveCharacterTo(currentCharacter, currentCharacter.Position.X-1, currentCharacter.Position.Y)
	case "a":
		moveCharacterTo(currentCharacter, currentCharacter.Position.X, currentCharacter.Position.Y-1)
	case "s":
		moveCharacterTo(currentCharacter, currentCharacter.Position.X+1, currentCharacter.Position.Y)
	case "d":
		moveCharacterTo(currentCharacter, currentCharacter.Position.X, currentCharacter.Position.Y+1)
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

func checkCollisions() {
	for characterIndex, character := range characters {
		for zombieIndex, zombie := range zombies {
			if character.Position.X == zombie.Position.X && character.Position.Y == zombie.Position.Y {
				printMap()
				fmt.Printf("%s e %s si sono scontrati nella stessa posizione!\n", character.Name, zombie.Name)
				time.Sleep(time.Second / 2)
				fight(characterIndex, zombieIndex)
			}
		}
	}
}

func fight(characterIndex, zombieIndex int) {
	character := &characters[characterIndex]
	zombie := &zombies[zombieIndex]

	fmt.Printf("%s attacca lo zombie %s!\n", character.Name, zombie.Name)
	time.Sleep(time.Second / 2)
	damage := rand.Intn(20) + character.Weapon.Damage
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
	time.Sleep(time.Second / 2)
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

func printMap() {
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
