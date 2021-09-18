package application

type IdentityLoginUseCase func(request *LoginIdentityRequest) (*LoginIdentityResponse, error)
type IdentityLogoutUseCase func(request *LogoutIdentityRequest) (*LogoutIdentityResponse, error)
type RegisterIdentityUseCase func(request *RegisterIdentityRequest) (*RegisterIdentityResponse, error)
