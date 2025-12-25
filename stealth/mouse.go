package stealth

import (
	"math"
	"math/rand"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type MouseController struct {
	cfg    *config.MouseMovementConfig
	random *rand.Rand
}

func NewMouseController(cfg *config.MouseMovementConfig) *MouseController {
	return &MouseController{
		cfg:    cfg,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type Point struct {
	X, Y float64
}

func (mc *MouseController) MoveToElement(page *rod.Page, element *rod.Element) error {
	if !mc.cfg.Enabled {
		return element.MoveMouseToCenter()
	}

	box, err := element.Box()
	if err != nil {
		return err
	}

	targetX := box.X + box.Width/2
	targetY := box.Y + box.Height/2

	currentX, currentY := mc.getCurrentPosition(page)

	return mc.moveWithBezier(page, currentX, currentY, targetX, targetY)
}

func (mc *MouseController) getCurrentPosition(page *rod.Page) (float64, float64) {
	return 100.0, 100.0
}

func (mc *MouseController) moveWithBezier(page *rod.Page, startX, startY, endX, endY float64) error {
	if !mc.cfg.BezierCurves {
		return page.Mouse.Move(endX, endY, 1)
	}

	points := mc.generateBezierPath(startX, startY, endX, endY)

	for i, point := range points {
		steps := 1
		if i > 0 {
			distance := mc.distance(points[i-1], point)
			steps = int(distance / 10)
			if steps < 1 {
				steps = 1
			}
		}

		velocity := mc.calculateVelocity(i, len(points))
		delay := time.Duration(float64(steps) / velocity * float64(time.Millisecond) * 10)

		if err := page.Mouse.Move(point.X, point.Y, steps); err != nil {
			return err
		}

		time.Sleep(delay)
	}

	if mc.cfg.Overshoot {
		overshootX := endX + mc.randomFloat(-5, 5)
		overshootY := endY + mc.randomFloat(-5, 5)
		page.Mouse.Move(overshootX, overshootY, 1)
		time.Sleep(20 * time.Millisecond)
	}

	if mc.cfg.MicroCorrections {
		for i := 0; i < mc.random.Intn(2)+1; i++ {
			correctionX := endX + mc.randomFloat(-2, 2)
			correctionY := endY + mc.randomFloat(-2, 2)
			page.Mouse.Move(correctionX, correctionY, 1)
			time.Sleep(10 * time.Millisecond)
		}
	}

	return page.Mouse.Move(endX, endY, 1)
}

func (mc *MouseController) generateBezierPath(x0, y0, x3, y3 float64) []Point {
	numPoints := 20 + mc.random.Intn(10)

	x1 := x0 + (x3-x0)*0.3 + mc.randomFloat(-50, 50)
	y1 := y0 + (y3-y0)*0.3 + mc.randomFloat(-50, 50)
	x2 := x0 + (x3-x0)*0.7 + mc.randomFloat(-50, 50)
	y2 := y0 + (y3-y0)*0.7 + mc.randomFloat(-50, 50)

	var points []Point
	for i := 0; i <= numPoints; i++ {
		t := float64(i) / float64(numPoints)
		point := mc.cubicBezier(t, x0, y0, x1, y1, x2, y2, x3, y3)
		points = append(points, point)
	}

	return points
}

func (mc *MouseController) cubicBezier(t, x0, y0, x1, y1, x2, y2, x3, y3 float64) Point {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t

	x := uuu*x0 + 3*uu*t*x1 + 3*u*tt*x2 + ttt*x3
	y := uuu*y0 + 3*uu*t*y1 + 3*u*tt*y2 + ttt*y3

	return Point{X: x, Y: y}
}

func (mc *MouseController) calculateVelocity(current, total int) float64 {
	baseVelocity := 2.0
	progress := float64(current) / float64(total)

	if progress < 0.2 {
		baseVelocity *= (1.0 + progress*2)
	} else if progress > 0.8 {
		baseVelocity *= (1.0 + (1.0-progress)*2)
	} else {
		baseVelocity *= 2.5
	}

	variance := 1.0 + mc.randomFloat(-mc.cfg.VelocityVariance, mc.cfg.VelocityVariance)
	return baseVelocity * variance
}

func (mc *MouseController) distance(p1, p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (mc *MouseController) randomFloat(min, max float64) float64 {
	return min + mc.random.Float64()*(max-min)
}

func (mc *MouseController) HoverElement(page *rod.Page, element *rod.Element) error {
	if err := mc.MoveToElement(page, element); err != nil {
		return err
	}

	hoverDuration := time.Duration(100+mc.random.Intn(300)) * time.Millisecond
	time.Sleep(hoverDuration)

	return nil
}

func (mc *MouseController) IdleWander(page *rod.Page) error {
	if mc.random.Float64() > 0.3 {
		return nil
	}

	viewport, err := page.Eval(`() => ({width: window.innerWidth, height: window.innerHeight})`)
	if err != nil {
		return nil
	}

	width := viewport.Value.Get("width").Num()
	height := viewport.Value.Get("height").Num()

	targetX := mc.randomFloat(width*0.2, width*0.8)
	targetY := mc.randomFloat(height*0.2, height*0.8)

	currentX, currentY := mc.getCurrentPosition(page)

	return mc.moveWithBezier(page, currentX, currentY, targetX, targetY)
}

func (mc *MouseController) ClickElement(page *rod.Page, element *rod.Element) error {
	if err := mc.MoveToElement(page, element); err != nil {
		return err
	}

	thinkDelay := time.Duration(300+mc.random.Intn(1000)) * time.Millisecond
	time.Sleep(thinkDelay)

	if err := page.Mouse.Down(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	clickDuration := time.Duration(50+mc.random.Intn(100)) * time.Millisecond
	time.Sleep(clickDuration)

	return page.Mouse.Up(proto.InputMouseButtonLeft, 1)
}
