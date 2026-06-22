X Package is meant to be gateway package to access all the infrastructure components and most used components in the application.
Example:

```go
import "github.com/imohamedsheta/xapp/app/x"

func main() {
	// for app config
	cfg := x.Config()

	// for database connection
	db := x.DB()

	// for redis connection
	redis := x.Redis()

	// for queue connection
	queue := x.Queue()

	// for websocket manager
	ws := x.Websocket()

	// for notification manager
	notify := x.Notify()

	// for logger manager
	logger := x.Logger()

	// for validator wrapper over go-playground/validator
	validator := x.Valid()

	// for ioc container
	app := x.App[any]()

	// for error
	error := x.Error()

	// for social authentication
	socialite := x.Social()

	// for inertia
	inertia := x.Inertia()

}
```