package initilizer

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	if godotenv.Load() != nil {
		fmt.Println("Faild to load env")
		return
	}
	fmt.Println("Successfully loaded env")
}
