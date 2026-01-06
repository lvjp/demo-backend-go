package healthcheck

import (
	"fmt"
	"net/http"

	"go.lvjp.me/demo-backend-go/internal/app/appcontext"
	"go.lvjp.me/demo-backend-go/pkg/buildinfo"
)

func Run(ctx *appcontext.AppContext) error {
	endpoint := fmt.Sprintf("http://%s/api/v0/misc/version", *ctx.Config.Server.ListenAddress)

	req, err := http.NewRequest(http.MethodHead, endpoint, nil)
	if err != nil {
		return fmt.Errorf("healthcheck request forging: %v", err)
	}

	req.Header.Set("User-Agent", "demo-backend-go/"+buildinfo.Get().Revision)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("healthcheck request execution: %v", err)
	}

	// HEAD request does not have a body
	resp.Body.Close()

	fmt.Fprintln(ctx.Output, "Endpoint:", endpoint)
	fmt.Fprintln(ctx.Output, "Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthcheck unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
