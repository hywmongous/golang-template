package routes

type Route interface {
	Setup()
}

type Routes []Route

func Factory(
	accountRoutes AccountRoutes,
	authenticationRoutes AuthenticationRoutes,
	sessionRoutes SessionRoutes,
	ticketRoutes TicketRoutes,
) Routes {
	return Routes{
		accountRoutes,
		authenticationRoutes,
		sessionRoutes,
		ticketRoutes,
	}
}

func (
	routes Routes,
) Setup() {
	for _, route := range routes {
		route.Setup()
	}
}
