package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	svg "github.com/ajstarks/svgo"

	"github.com/jpoz/sktch"
)

const style = `style="stroke:black;fill-opacity:0.0;"`
const greenStyle = `style="stroke:green;fill-opacity:0.0;"`
const redStyle = `style="stroke:red;fill-opacity:0.0;"`

const debug = false

func main() {
	s, err := sktch.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	s.AddSketch("/", func(i sktch.Inputs, c *svg.SVG) error {
		centerX := i.Width / 2
		centerY := i.Height / 2
		maxR := min(centerX, centerY) - 10

		if debug {
			c.Circle(
				centerX,
				centerY,
				1,
				style,
			)
		}

		for j := 0; j < 39; j++ {
			startAngle := rand.Float64() * 360.0
			// endAngle := startAngle + (rand.Float64() * (360.0 - startAngle))
			endAngle := startAngle + 15
			r := rand.Float64()*float64(maxR) + 5

			sx := r * math.Sin(startAngle)
			sy := r * math.Cos(startAngle)

			ex := r * math.Sin(endAngle)
			ey := r * math.Cos(endAngle)

			c.Arc(
				centerX+int(sx),
				centerY+int(sy),
				int(r),
				int(r),
				0,
				false,
				false,
				centerX+int(ex),
				centerY+int(ey),
				style,
			)

			if debug {
				c.Circle(
					centerX+int(sx),
					centerY+int(sy),
					1,
					greenStyle,
				)
				c.Circle(
					centerX+int(ex),
					centerY+int(ey),
					1,
					redStyle,
				)
			}
		}
		// c.Arc(110, 60, 50, 50, 0, false, true, 60, 110, style)
		// c.Arc(60, 110, 50, 50, 0, false, true, 110, 60, style)
		// c.Arc(60, 110, 50, 50, 0, false, true, 110, 60, style)

		return nil
	})

	err = s.ListenAndServe()
	fmt.Println(err)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
