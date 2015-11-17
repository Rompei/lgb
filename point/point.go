package point

// Point : Pixel of field
type Point struct {
	X       int
	Y       int
	Str     string
	IsAlive bool
}

// NewPoint : Constructor of Point
// @Param X X座標
// @param Y Y座標
// @Param str 文字
// @Param isAlive 生きているかどうか
// return Point
func NewPoint(x, y int, str string, isAlive bool) Point {
	return Point{
		X:       x,
		Y:       y,
		Str:     str,
		IsAlive: isAlive,
	}
}
