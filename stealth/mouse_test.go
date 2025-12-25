package stealth

import (
	"math"
	"testing"

	"github.com/automation-poc/browser-automation/config"
)

func TestNewMouseController(t *testing.T) {
	cfg := &config.MouseMovementConfig{
		Enabled:          true,
		BezierCurves:     true,
		Overshoot:        true,
		MicroCorrections: true,
		VelocityVariance: 0.3,
	}

	mc := NewMouseController(cfg)

	if mc == nil {
		t.Fatal("Expected non-nil MouseController")
	}

	if mc.cfg != cfg {
		t.Error("Configuration not set correctly")
	}

	if mc.random == nil {
		t.Error("Random generator not initialized")
	}
}

func TestGenerateBezierPath(t *testing.T) {
	cfg := &config.MouseMovementConfig{
		Enabled:      true,
		BezierCurves: true,
	}
	mc := NewMouseController(cfg)

	startX, startY := 0.0, 0.0
	endX, endY := 100.0, 100.0

	points := mc.generateBezierPath(startX, startY, endX, endY)

	if len(points) < 20 {
		t.Errorf("Expected at least 20 points, got %d", len(points))
	}

	if points[0].X != startX || points[0].Y != startY {
		t.Errorf("First point should be start position (%f, %f), got (%f, %f)",
			startX, startY, points[0].X, points[0].Y)
	}

	lastIdx := len(points) - 1
	if points[lastIdx].X != endX || points[lastIdx].Y != endY {
		t.Errorf("Last point should be end position (%f, %f), got (%f, %f)",
			endX, endY, points[lastIdx].X, points[lastIdx].Y)
	}

	for i := 0; i < len(points)-1; i++ {
		distance := mc.distance(points[i], points[i+1])
		if distance < 0 {
			t.Error("Distance between points should be non-negative")
		}
		if distance > 100 {
			t.Errorf("Distance between consecutive points too large: %f", distance)
		}
	}
}

func TestCubicBezier(t *testing.T) {
	cfg := &config.MouseMovementConfig{}
	mc := NewMouseController(cfg)

	x0, y0 := 0.0, 0.0
	x1, y1 := 33.0, 33.0
	x2, y2 := 66.0, 66.0
	x3, y3 := 100.0, 100.0

	tests := []struct {
		t        float64
		expectedX float64
		expectedY float64
	}{
		{t: 0.0, expectedX: 0.0, expectedY: 0.0},
		{t: 1.0, expectedX: 100.0, expectedY: 100.0},
		{t: 0.5, expectedX: 50.0, expectedY: 50.0},
	}

	for _, tt := range tests {
		point := mc.cubicBezier(tt.t, x0, y0, x1, y1, x2, y2, x3, y3)

		if math.Abs(point.X-tt.expectedX) > 0.1 {
			t.Errorf("At t=%f, expected X=%f, got %f", tt.t, tt.expectedX, point.X)
		}
		if math.Abs(point.Y-tt.expectedY) > 0.1 {
			t.Errorf("At t=%f, expected Y=%f, got %f", tt.t, tt.expectedY, point.Y)
		}
	}
}

func TestCalculateVelocity(t *testing.T) {
	cfg := &config.MouseMovementConfig{
		VelocityVariance: 0.3,
	}
	mc := NewMouseController(cfg)

	total := 100

	startVelocity := mc.calculateVelocity(5, total)
	midVelocity := mc.calculateVelocity(50, total)
	endVelocity := mc.calculateVelocity(95, total)

	if startVelocity <= 0 {
		t.Error("Velocity should be positive")
	}

	if midVelocity < startVelocity*0.8 {
		t.Log("Expected mid-section velocity to be higher (but random variance may affect this)")
	}

	if endVelocity <= 0 {
		t.Error("End velocity should be positive")
	}
}

func TestDistance(t *testing.T) {
	cfg := &config.MouseMovementConfig{}
	mc := NewMouseController(cfg)

	tests := []struct {
		p1       Point
		p2       Point
		expected float64
	}{
		{Point{0, 0}, Point{3, 4}, 5.0},
		{Point{0, 0}, Point{0, 0}, 0.0},
		{Point{1, 1}, Point{4, 5}, 5.0},
	}

	for _, tt := range tests {
		result := mc.distance(tt.p1, tt.p2)
		if math.Abs(result-tt.expected) > 0.1 {
			t.Errorf("Distance from %v to %v: expected %f, got %f",
				tt.p1, tt.p2, tt.expected, result)
		}
	}
}

func TestRandomFloat(t *testing.T) {
	cfg := &config.MouseMovementConfig{}
	mc := NewMouseController(cfg)

	min := -10.0
	max := 10.0

	for i := 0; i < 100; i++ {
		result := mc.randomFloat(min, max)

		if result < min || result > max {
			t.Errorf("Random float %f is outside range [%f, %f]", result, min, max)
		}
	}
}

func BenchmarkGenerateBezierPath(b *testing.B) {
	cfg := &config.MouseMovementConfig{
		BezierCurves: true,
	}
	mc := NewMouseController(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.generateBezierPath(0, 0, 1000, 1000)
	}
}

func BenchmarkCalculateVelocity(b *testing.B) {
	cfg := &config.MouseMovementConfig{
		VelocityVariance: 0.3,
	}
	mc := NewMouseController(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.calculateVelocity(50, 100)
	}
}
