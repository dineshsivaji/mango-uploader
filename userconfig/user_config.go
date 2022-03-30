package userconfig

// import (
// 	"crypto/aes"
// 	"encoding/hex"
// 	"fmt"
// )

// func CheckError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func Manage() {

// 	// cipher key
// 	key := "~!@#$%^&*()_+M@nG0uP70@d3R-=[]{}"

// 	// plaintext
// 	pt := "dinesh"

// 	c := EncryptAES([]byte(key), pt)

// 	// plaintext
// 	fmt.Println(pt)

// 	// ciphertext
// 	fmt.Println(c)

// 	// decrypt
// 	DecryptAES([]byte(key), c)
// }

// func EncryptAES(key []byte, plaintext string) string {

// 	c, err := aes.NewCipher(key)
// 	CheckError(err)

// 	out := make([]byte, len(plaintext))

// 	c.Encrypt(out, []byte(plaintext))

// 	return hex.EncodeToString(out)
// }

// func DecryptAES(key []byte, ct string) {
// 	ciphertext, _ := hex.DecodeString(ct)

// 	c, err := aes.NewCipher(key)
// 	CheckError(err)

// 	pt := make([]byte, len(ciphertext))
// 	c.Decrypt(pt, ciphertext)

// 	s := string(pt[:])
// 	fmt.Println("DECRYPTED:", s)
// }

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

var masterKey string = "~!@#$%^&*()_+M@nG0uP70@d3R-=[]{}"

func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func Manage() {
	// var first string
	// fmt.Println("Enter the password : ")
	// // Taking input from user
	// fmt.Scanln(&first)
	// data := []byte(first)

	// key := []byte("~!@#$%^&*()_+M@nG0uP70@d3R-=[]{}")
	// fmt.Println(hex.EncodeToString(key))
	data := []byte("DINESH")
	ciphertext, err := Encrypt([]byte(masterKey), data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ciphertext: %s\n", hex.EncodeToString(ciphertext))

	plaintext, err := Decrypt([]byte(masterKey), ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("plaintext: %s\n", plaintext)
}
