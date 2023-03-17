package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ipReq struct {
	out chan string
}

type ipCommand struct {
	input    chan ipReq
	portSpec string
	cachedIP string
}

func (i *ipCommand) Serve(ctx context.Context) error {
	for {
		select {
		case m := <-i.input:
			if err := i.consumeInput(ctx, m); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (i *ipCommand) consumeInput(ctx context.Context, r ipReq) error {
	if i.cachedIP == "" {
		if err := i.requestIP(ctx); err != nil {
			select {
			case r.out <- "[ip] failed to discover because " + err.Error():
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		}
	}

	response := fmt.Sprintf("[ip] `%s%s`", i.cachedIP, i.portSpec)
	select {
	case r.out <- response:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (i *ipCommand) requestIP(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api64.ipify.org", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	i.cachedIP = string(body)
	return nil
}
