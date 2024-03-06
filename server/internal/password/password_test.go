package password

import "testing"

func TestPassword(t *testing.T) {
	password := "password123"
	encodedHash, err := Hash(password)
	if err != nil {
		t.Fatal("failed to hash password", err)
	}
	err = Verify(password, encodedHash)
	if err != nil {
		t.Fatal("failed to verify correct password", err)
	}
	err = Verify("wrongpassword", encodedHash)
	if err.Error() != invalidPasswordError.Error() {
		t.Fatal("failed to verify wrong password", err)
	}
}
