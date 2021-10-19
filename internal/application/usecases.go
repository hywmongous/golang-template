package application

type (
	IdentityLoginUseCase    func(request *LoginIdentityRequest) (*LoginIdentityResponse, error)
	IdentityLogoutUseCase   func(request *LogoutIdentityRequest) (*LogoutIdentityResponse, error)
	RegisterIdentityUseCase func(request *RegisterIdentityRequest) (*RegisterIdentityResponse, error)
)
