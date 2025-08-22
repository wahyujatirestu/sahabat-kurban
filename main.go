// @title Sahabat Kurban
// @version 1.0
// @description REST API untuk Sahabat Kurban
// @termsOfService http://swagger.io/terms/

// @contact.name Restu Adi Wahyujati
// @contact.email restu@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" + a space + your JWT token.

package main

func main() {
	server := NewServer()
	defer server.Close()

	server.Run()
}