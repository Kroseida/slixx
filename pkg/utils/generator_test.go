package utils_test

import (
	"kroseida.org/slixx/pkg/utils"
	"testing"
)

func Test_GenerateSecureToken(t *testing.T) {
	for i := 0; i < 500000; i++ {
		token := []string{
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(5),
			utils.GenerateSecureToken(10),
			utils.GenerateSecureToken(10),
			utils.GenerateSecureToken(10),
			utils.GenerateSecureToken(10),
			utils.GenerateSecureToken(10),
			utils.GenerateSecureToken(15),
			utils.GenerateSecureToken(15),
			utils.GenerateSecureToken(15),
			utils.GenerateSecureToken(15),
			utils.GenerateSecureToken(15),
			utils.GenerateSecureToken(20),
			utils.GenerateSecureToken(20),
			utils.GenerateSecureToken(20),
		}
		for i := 0; i < len(token); i++ {
			for j := i + 1; j < len(token); j++ {
				if token[i] == token[j] {
					t.Error("Token collision")
					return
				}
			}
		}
	}
}
