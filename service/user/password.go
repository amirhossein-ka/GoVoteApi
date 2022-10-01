package user

import "golang.org/x/crypto/bcrypt"

func hashPass(pass string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(h), err
}

// comparePass compare given password (gpass) with hashed password from db (dpass)
// return true if they are same
func comparePass(gpass, dpass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dpass), []byte(gpass))
	return err == nil
}
