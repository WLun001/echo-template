package env

import (
	"os"
)

var (
	//isProduction = strings.Contains(strings.ToLower(os.Getenv("TAIZI_APP_ENV")), "prod")

	// DefaultDB :
	DefaultDB = os.Getenv("DATABASE_NAME")
)

//func IsProduction() bool {
//	return isProduction
//}
