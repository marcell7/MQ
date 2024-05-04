package broker

import (
	"fmt"
	"time"
)

// Basic id generator
func generateId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
