package helper

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func HashPassword(password string) (string, string, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return "", "", err
	}
	passwordWithSalt := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return string(hash), salt, nil
}

func CheckPassword(storedHash, storedSalt, password string) bool {
	passwordWithSalt := password + storedSalt
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(passwordWithSalt))
	return err == nil
}


// GerarSenhaPadrao gera uma senha padrão usando o primeiro nome em minúsculas e o ano atual.
func GerarSenhaPadrao(nome string) string {
	fields := strings.Fields(nome)
	nomeNovo := RemoveAcentos(fields[0])
	currentYear := strconv.Itoa(time.Now().Year())
	return strings.ToLower(nomeNovo) + currentYear
}

// RemoveAcentos remove acentos de uma string
func RemoveAcentos(s string) string {
	t := norm.NFD.String(s)
	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // Mn: marcas não espaçadoras
			return -1
		}
		return r
	}, t)
}
