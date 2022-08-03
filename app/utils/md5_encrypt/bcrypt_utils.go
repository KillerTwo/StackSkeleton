package md5_encrypt

import "golang.org/x/crypto/bcrypt"

// BcryptEncode MD5
func BcryptEncode(source string) (encode string, err error) {
	pwd := []byte(source)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	s := string(hash)
	return s, nil
}

// BcryptCompare compare
func BcryptCompare(encode, source string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encode), []byte(source))
	if err != nil {
		return false
	}
	return true
}
