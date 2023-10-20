package database

import "time"

func (db *DB) RevokeToken(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.RevokedTokens[tokenString]
	if !ok {
		dbStructure.RevokedTokens[tokenString] = Token{
			ID:        tokenString,
			RevokedAt: time.Now().UTC(),
		}
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CheckTokenRevoked(tokenString string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := dbStructure.RevokedTokens[tokenString]

	return ok, nil
}
