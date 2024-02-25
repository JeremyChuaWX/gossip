package password

import "golang.org/x/crypto/bcrypt"

const BCRYPT_COST = 12

func Hash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		return nil, &PasswordError{
			message: "Error hashing password",
			error:   err,
		}
	}
	return hash, nil
}

func Compare(hash []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return &PasswordError{
			message: "Passwords do not match",
			error:   err,
		}
	}
	return nil
}
