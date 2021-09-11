package endpoint

import (
	"com/justin/micro/kit/library-user-service/dto"
	"com/justin/micro/kit/library-user-service/service"
	"context"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type UserEndpoints struct {
	RegisterEndpoint    endpoint.Endpoint
	FindByIDEndpoint    endpoint.Endpoint
	FindByEmailEndpoint endpoint.Endpoint
}

func MakeRegisterEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*dto.RegisterUser)
		user, err := svc.Register(ctx, req)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func MakeFindByIDEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id, _ := strconv.ParseUint(request.(string), 10, 64)
		user, err := svc.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func MakeFindByEmailEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		email := request.(string)
		user, err := svc.FindByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}
