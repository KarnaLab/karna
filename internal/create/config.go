package create

import "fmt"

func checkRequirements(APIEndpoint, verb *string) (err error) {
	if *APIEndpoint != "EDGE" || *APIEndpoint != "REGIONAL" || *APIEndpoint != "PRIVATE" {
		return fmt.Errorf("API endpoint must be 'EDGE' or 'REGIONAL' or 'PRIVATE'")
	}

	verbsAuthorized := [7]string{"ANY", "GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"}

	found := false

	for _, item := range verbsAuthorized {
		if item == *verb {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("HTTP verb must be 'ANY' or 'GET' or 'PUT' or 'PATCH' or 'POST' or 'DELETE' or 'OPTIONS'")
	}

	return
}
