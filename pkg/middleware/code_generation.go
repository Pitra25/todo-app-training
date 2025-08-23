package middleware

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/sirupsen/logrus"
)

func CodeGeneration(lenght int) (int, error) {
	codeStr := ""
	for range lenght {
		codeStr += fmt.Sprint(rand.IntN(9))
	}

	codeInt, err := strconv.Atoi(codeStr)
	if err != nil {
		logrus.Error("(pkg/middleware/code_generation) failed to generate code")
		return 0, err
	}

	return codeInt, nil
}
