package circuitbreakers

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
)

func Hystrix(commandName, fallbackMsg string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = hystrix.Do(
				commandName,
				func() error {
					response, err = next(ctx, request)
					return err
				},
				func(err error) error {
					fmt.Println("fallbackErrorDesc", err.Error())
					response = struct {
						Fallback string `json:"fallback"`
					}{fallbackMsg}
					return nil
				})
			if err != nil {
				return nil, err
			}
			return
		}
	}
}
