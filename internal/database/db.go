package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp    `json:"chirps"`
	Users         map[int]User     `json:"users"`
	RevokedTokens map[string]Token `json:"revoked_tokens"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return &DB{}, err
	}

	return db, nil
}

// createDB creates a new database with proper structure
func (db *DB) createDB() error {
	return db.writeDB(
		DBStructure{
			Chirps:        map[int]Chirp{},
			Users:         map[int]User{},
			RevokedTokens: map[string]Token{},
		},
	)
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
