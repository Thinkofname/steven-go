// +build mobile

package platform

import "golang.org/x/mobile/app"

func run(handler Handler) {
	app.Run(app.Callbacks{
		Start: handler.Start,
		Draw:  handler.Draw,
	})
}
