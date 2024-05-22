package database

import (
	"encoding/json"
	"os"
	"sync"
)

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

func loadDatabase() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := os.Stat(databaseFile)
	if os.IsNotExist(err) {
		db = &Database{
			Chirps: make(map[int]Chirp),
			NextID: 1,
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
