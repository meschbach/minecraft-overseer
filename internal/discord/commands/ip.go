package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
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
			fmt.Printf("[ip] failed to discover because %s\n", err.Error())
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
	expiringContext, done := context.WithTimeout(ctx, 1*time.Second)
	defer done()
	req, err := http.NewRequestWithContext(expiringContext, http.MethodGet, "https://api.ipify.org", nil)
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
