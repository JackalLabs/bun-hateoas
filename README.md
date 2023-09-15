# Bun HATEOAS

This is a HATEOAS package for [bunrouter](https://github.com/uptrace/bunrouter).

## Usage
 Very similar to how you group routes normally using bun, you can replace them with the `WithGroup` function from `bun-hateoas` and it will add the group to the map and running methods of the group handler will work like normal but it will handle adding the routes to the group in the map.
```go
import bunhateoas "github.com/JackalLabs/bun-hateoas"

hateoas := bunhateoas.NewHATEOAS()

hateoas.WithGroup(router.Use(CustomCors(nimbusUrl, true)), "", func(group *bunhateoas.Group) {
	group.GET("/version", routeHandlers.VersionHandler(Version))
	group.GET("/hook", routeHandlers.WebhookHandler())
	group.POST("/register", routeHandlers.UserRegistrationHandler(dbSession, jwtKey))
	group.POST("/login", routeHandlers.LoginHandler(dbSession, jwtKey))
	group.GET("/", hateoas.Handler())
})

hateoas.WithGroup(router.Use(CustomCors("*", false)), "/pubapi", func(group *bunhateoas.Group) {
    group.GET("/download/:id", japicore.DownloadHandler(fileIo))
    group.GET("/d/:id", japicore.DownloadHandler(fileIo))
})
```

