package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateOrderID() string {
	today := time.Now().Format("20060102")
	random := uuid.New().String()[:8]
	return fmt.Sprintf("ORDER-%s-%s", today, random)
}
