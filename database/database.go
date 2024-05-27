package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"
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
			Chirps:        make(map[int]Chirp),
			Users:         make(map[int]User),
			NextID:        1,
			NextUserID:    1,
			RefreshTokens: make(map[string]RefreshToken),
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

func GetChirpsByAuthorID(authorID int) ([]Chirp, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	chirps := make([]Chirp, 0)
	for _, chirp := range db.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
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

func CreateChirp(body string, userId int) (Chirp, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	cleanedBody := replaceProfaneWords(body)
	if len(cleanedBody) > 140 {
		return Chirp{}, ErrChirpTooLong
	}

	chirp := Chirp{
		ID:       db.NextID,
		Body:     cleanedBody,
		AuthorID: userId,
	}

	db.Chirps[chirp.ID] = chirp
	db.NextID++

	err := saveDatabase()
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func DeleteChirp(id int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, ok := db.Chirps[id]
	if !ok {
		return ErrChirpNotFound
	}

	delete(db.Chirps, id)

	err := saveDatabase()
	if err != nil {
		return err
	}

	return nil
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
		ID:          db.NextUserID,
		Email:       email,
		Password:    string(hashedPassword),
		IsChirpyRed: false,
	}

	db.Users[user.ID] = user
	db.NextUserID++

	err = saveDatabase()
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UpgradeUserToChirpyRed(userID int) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	user, ok := db.Users[userID]
	if !ok {
		return ErrUserNotFound
	}

	user.IsChirpyRed = true
	db.Users[userID] = user

	err := saveDatabase()
	if err != nil {
		return err
	}

	return nil
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

func UpdateUser(id int, email, password string) (User, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	user, ok := db.Users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}

	if email != "" {
		user.Email = email
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}
		user.Password = string(hashedPassword)
	}

	db.Users[id] = user

	err := saveDatabase()
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateRefreshToken(userID int, expiresIn time.Duration) (string, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	refreshToken := RefreshToken{
		Token:     hex.EncodeToString(token),
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiresIn),
	}

	db.RefreshTokens[refreshToken.Token] = refreshToken

	err = saveDatabase()
	if err != nil {
		return "", err
	}

	return refreshToken.Token, nil
}

func GetRefreshToken(token string) (RefreshToken, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	refreshToken, ok := db.RefreshTokens[token]
	if !ok {
		return RefreshToken{}, errors.New("refresh token not found")
	}

	return refreshToken, nil
}

func DeleteRefreshToken(token string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, ok := db.RefreshTokens[token]
	if !ok {
		return errors.New("refresh token not found")
	}

	delete(db.RefreshTokens, token)

	err := saveDatabase()
	if err != nil {
		return err
	}

	return nil
}
