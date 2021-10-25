package application

import "context"

type (
	IdentityLoginUseCase    func(ctx context.Context, request *LoginIdentityRequest) (*LoginIdentityResponse, error)
	IdentityLogoutUseCase   func(ctx context.Context, request *LogoutIdentityRequest) (*LogoutIdentityResponse, error)
	RegisterIdentityUseCase func(ctx context.Context, request *RegisterIdentityRequest) (*RegisterIdentityResponse, error)
)
