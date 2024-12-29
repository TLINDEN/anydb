/*
Copyright © 2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/term"
)

const (
	ArgonMem      uint32 = 64 * 1024
	ArgonIter     uint32 = 5
	ArgonParallel uint8  = 2
	ArgonSaltLen  int    = 16
	ArgonKeyLen   uint32 = 32
	B64SaltLen    int    = 22
)

type Key struct {
	Salt []byte
	Key  []byte
}

// called from interactive thread, hides  input and returns clear text
// password
func AskForPassword() ([]byte, error) {
	fmt.Fprint(os.Stderr, "Password: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Fprintln(os.Stderr)

	return pass, nil
}

// We're  using the  Argon2id  key derivation  algorithm  to derive  a
// secure  key from  the given  password. This  is important,  because
// users might  use unsecure  passwords. The resulting  encrypted data
// might of course easily be decrypted  using brute force methods if a
// weak password  was used, but  that would  cost, because of  the key
// derivation. It does several rounds  of hash calculations which take
// a considerable  amount of cpu  time. For  our legal user  that's no
// problem because it's being executed  only once, but an attacker has
// to do it in a forever loop, which will take a lot of time.
func DeriveKey(password []byte, salt []byte) (*Key, error) {
	if salt == nil {
		// none given, new password
		newsalt, err := GetRandom(ArgonSaltLen, ArgonSaltLen)
		if err != nil {
			return nil, err
		}

		salt = newsalt
	}

	hash := argon2.IDKey(
		[]byte(password), salt,
		ArgonIter,
		ArgonMem,
		ArgonParallel,
		ArgonKeyLen,
	)

	return &Key{Key: hash, Salt: salt}, nil
}

// Retrieve a random chunk of given size
func GetRandom(size int, capacity int) ([]byte, error) {
	buf := make([]byte, size, capacity)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve random bytes: %w", err)
	}

	return buf, nil
}

// Encrypt clear  text given in  attr using ChaCha20  and auhtenticate
// using the mac Poly1305. The cipher text will be put into attr, thus
// modifying it.
//
// The cipher text consists of:
// password-salt) + (12 byte nonce + ciphertext + 16 byte mac)
func Encrypt(pass []byte, attr *DbAttr) error {
	key, err := DeriveKey(pass, nil)
	if err != nil {
		return err
	}

	aead, err := chacha20poly1305.New(key.Key)
	if err != nil {
		return fmt.Errorf("failed to create AEAD cipher: %w", err)
	}

	total := aead.NonceSize() + len(attr.Val) + aead.Overhead()

	nonce, err := GetRandom(aead.NonceSize(), total)
	if err != nil {
		return err
	}

	cipher := aead.Seal(nonce, nonce, attr.Val, nil)

	attr.Val = append(attr.Val, key.Salt...)
	attr.Val = append(attr.Val, cipher...)

	attr.Encrypted = true

	return nil
}

// Do the reverse
func Decrypt(pass []byte, cipherb []byte) ([]byte, error) {
	if len(cipherb) < B64SaltLen {
		return nil, fmt.Errorf("encrypted cipher block too small")
	}

	key, err := DeriveKey(pass, cipherb[0:B64SaltLen])
	if err != nil {
		return nil, err
	}

	cipher := cipherb[B64SaltLen:]

	aead, err := chacha20poly1305.New(key.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD cipher: %w", err)
	}

	if len(cipher) < aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := cipher[:aead.NonceSize()], cipher[aead.NonceSize():]

	return aead.Open(nil, nonce, ciphertext, nil)
}
