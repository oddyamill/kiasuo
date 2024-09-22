package crypto

import "testing"

func TestCrypto(t *testing.T) {
	text := "как стать python hacker домашних условиях?"
	encrypted := Encrypt(text)

	if encrypted.Encrypted == "" {
		t.Errorf("Encrypt().Value = %s; want non-empty string", encrypted)
	}

	t.Logf("Encrypt().Encrypted = %s", encrypted.Encrypted)

	decrypted := encrypted.Decrypt()

	if decrypted == "" {
		t.Errorf("Crypt.Decrypt() = %s; want non-empty string", decrypted)
	}

	if decrypted != text {
		t.Errorf("Crypt.Decrypt() = %s; want %s", decrypted, text)
	}
}
