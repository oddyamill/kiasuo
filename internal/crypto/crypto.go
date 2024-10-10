package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/kiasuo/bot/internal/helpers"
	"io"
)

var key []byte

func init() {
	if helpers.IsTesting() {
		// fake key for tests
		key = helpers.StringToBytes("d813c247bb06c2c8ac84ef4231658b2b")
	} else {
		key = helpers.StringToBytes(helpers.GetEnv("CRYPTO"))
	}
}

type Crypt struct {
	Encrypted string
}

func Encrypt(text string) Crypt {
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], helpers.StringToBytes(text))

	return Crypt{hex.EncodeToString(ciphertext)}
}

func (c *Crypt) Decrypt() string {
	if c.Encrypted == "" {
		return ""
	}

	ciphertext, err := hex.DecodeString(c.Encrypted)

	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic(errors.New("ciphertext too short"))
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return helpers.BytesToString(ciphertext)
}

func (c *Crypt) Scan(src any) error {
	text, ok := src.(string)

	if !ok {
		return nil
	}

	c.Encrypted = text
	return nil
}
