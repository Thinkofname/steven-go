// +build mobile

package platform

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/geom"
)

func run(handler Handler) {
	app.Run(app.Callbacks{
		Start: handler.Start,
		Draw:  handler.Draw,
	})
}

func size() (int, int) {
	return int(geom.Width.Px()), int(geom.Height.Px())
}
