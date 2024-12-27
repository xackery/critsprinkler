package monitor

/*
import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

// intersectwalk.Rectangle calculates the intersection of two walk.Rectangles.
func intersectRectangle(r1, r2 walk.Rectangle) (walk.Rectangle, bool) {
	intersect := walk.Rectangle{
		Left:   max(r1.Left, r2.Left),
		Top:    max(r1.Top, r2.Top),
		Right:  min(r1.Right, r2.Right),
		Bottom: min(r1.Bottom, r2.Bottom),
	}
	if intersect.Right > intersect.Left && intersect.Bottom > intersect.Top {
		return intersect, true
	}
	return walk.Rectangle{}, false
}

// RectangleArea calculates the area of a walk.Rectangle.
func RectangleArea(r walk.Rectangle) int32 {
	return (r.Right - r.Left) * (r.Bottom - r.Top)
}

// max returns the maximum of two int32 values.
func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two int32 values.
func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// MonitorFromWindow returns the monitor index that the window is primarily on.
func MonitorFromWindow(hwnd windows.HWND) (int, error) {
	// Get the window walk.Rectangleangle
	var windowwalk.Rectangle windows.walk.Rectangle
	if !win.GetWindowwalk.Rectangle(hwnd, &windowwalk.Rectangle) {
		return -1, fmt.Errorf("failed to get window walk.Rectangle")
	}

	// Prepare data for enumeration
	type monitorInfo struct {
		index int
		walk.Rectangle  walk.Rectangle
	}
	monitors := []monitorInfo{}
	callback := syscall.NewCallback(func(hMonitor win.HMONITOR, hdcMonitor win.HDC, lprcMonitor *win.walk.Rectangle, dwData win.LPARAM) win.BOOL {
		monitors = append(monitors, monitorInfo{
			index: len(monitors),
			walk.Rectangle: walk.Rectangle{
				Left:   lprcMonitor.Left,
				Top:    lprcMonitor.Top,
				Right:  lprcMonitor.Right,
				Bottom: lprcMonitor.Bottom,
			},
		})
		return 1 // Continue enumeration
	})

	// Enumerate monitors
	if !win.EnumDisplayMonitors(0, nil, callback, 0) {
		return -1, fmt.Errorf("failed to enumerate monitors")
	}

	// Find the monitor with the largest intersection
	var bestMonitor int = -1
	var maxArea int32 = 0
	windowwalk.Rectangle := walk.Rectangle{
		Left:   windowwalk.Rectangle.Left,
		Top:    windowwalk.Rectangle.Top,
		Right:  windowwalk.Rectangle.Right,
		Bottom: windowwalk.Rectangle.Bottom,
	}
	for _, monitor := range monitors {
		if intersect, ok := intersectwalk.Rectangle(windowwalk.Rectangle, monitor.walk.Rectangle); ok {
			area := walk.RectangleArea(intersect)
			if area > maxArea {
				maxArea = area
				bestMonitor = monitor.index
			}
		}
	}

	if bestMonitor == -1 {
		return -1, fmt.Errorf("window is not primarily on any monitor")
	}
	return bestMonitor, nil
}
*/
