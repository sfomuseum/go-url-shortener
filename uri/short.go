package uri

import (
	"math/rand"
)

const chars string = "qwrtypsdfghjklzxcvbnm0123456789"

func GenerateShortURI(length int) string {

	id := ""

	for len(id) < length {

		i := rand.Intn(len(chars))
		id = id + string(chars[i])
	}

	return id
}
