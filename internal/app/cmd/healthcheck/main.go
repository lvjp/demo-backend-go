package healthcheck

import (
	"fmt"
	"net/http"
	"os"

	"go.lvjp.me/demo-backend-go/internal/app/config"
	"go.lvjp.me/demo-backend-go/pkg/buildinfo"
)

func Run() {
	config, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load configuration: %v\n", err)
		os.Exit(1)
	}

	endpoint := fmt.Sprintf("http://%s/api/v0/misc/version", *config.Server.ListenAddress)

	req, err := http.NewRequest(http.MethodHead, endpoint, nil)
	if err != nil {
		fmt.Println("Cannot forge the healthcheck request:", err.Error())
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "demo-backend-go/"+buildinfo.Get().Revision)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// HEAD request does not have a body
	resp.Body.Close()

	fmt.Println("Endpoint:", endpoint)
	fmt.Println("Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected status code:", resp.StatusCode)
		os.Exit(1)
	}
}
