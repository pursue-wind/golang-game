package main

// CheckCollision 碰撞检查
func CheckCollision(obj1, obj2 GameAttr) bool {
	x1, y1, w1, h1 := obj1.X(), obj1.Y(), obj1.Width(), obj1.Height()
	x2, y2, w2, h2 := obj2.X(), obj2.Y(), obj2.Width(), obj2.Height()

	if x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2 {
		return true // 发生了碰撞
	}

	return false // 没有碰撞
}
