package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)
import "golang.org/x/crypto/bcrypt"

const databaseFile = "database.json"

var (
	db      *Database
	once    sync.Once
	dbMutex sync.RWMutex
)

func Init() error {
	var err error
	once.Do(func() {
		err = loadDatabase()
	})
	return err
}

func RemoveDatabase() error {
	err := os.Remove(databaseFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to remove database file: %v", err)
		return err

	}
	return nil
}

func loadDatabase() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := os.Stat(databaseFile)
	if os.IsNotExist(err) {
		db = &Database{
			Chirps:     make(map[int]Chirp),
			Users:      make(map[int]User),
			NextID:     1,
			NextUserID: 1,
		}
		return nil
	}

	data, err := os.ReadFile(databaseFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return err
	}

	return nil
}

func GetChirps() ([]Chirp, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	chirps := make([]Chirp, 0, len(db.Chirps))
	for _, chirp := range db.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func GetChirpByID(id int) (Chirp, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	chirp, ok := db.Chirps[id]
	if !ok {
		return Chirp{}, ErrChirpNotFound
	}

	return chirp, nil
}

func CreateChirp(body string) (Chirp, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	cleanedBody := replaceProfaneWords(body)
	if len(cleanedBody) > 140 {
		return Chirp{}, ErrChirpTooLong
	}

	chirp := Chirp{
		ID:   db.NextID,
		Body: cleanedBody,
	}

	db.Chirps[chirp.ID] = chirp
	db.NextID++

	err := saveDatabase()
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func saveDatabase() error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(databaseFile, data, 0644)
}

func CreateUser(email, password string) (User, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// moved here as got a write lock error using the function - need to fix
	for _, existingUser := range db.Users {
		if existingUser.Email == email {
			return User{}, ErrUserExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       db.NextUserID,
		Email:    email,
		Password: string(hashedPassword),
	}

	db.Users[user.ID] = user
	db.NextUserID++

	err = saveDatabase()
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	for _, user := range db.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrUserNotFound
}
