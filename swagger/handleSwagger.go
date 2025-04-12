package swagger

import (
	"fmt"
	"net/http"
	"os"
)

// handleSwagger serves the Swagger UI and the Swagger spec (swagger.yaml)
func HandleSwagger(w http.ResponseWriter, r *http.Request) {
	swaggerFilePath := "/Users/iyadi/playground/eth-toy-client/eth-toy-client/swagger/swagger.yaml" // Path to the swagger.yaml file

	// Check if the Swagger file exists
	if _, err := os.Stat(swaggerFilePath); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("Swagger file not found at %s", swaggerFilePath), http.StatusNotFound)
		return
	}

	// Serve the swagger.yaml file
	http.ServeFile(w, r, swaggerFilePath)
}
