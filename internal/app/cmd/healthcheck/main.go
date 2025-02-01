package healthcheck

import (
	"fmt"
	"net/http"
	"os"
)

func Run() error {
	endpoint := "http://:8080"
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Println("Endpoint:", endpoint)
	fmt.Println("Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
