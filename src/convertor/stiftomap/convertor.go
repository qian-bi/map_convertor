package stiftomap

// ConvertMap is to convert STIF map to MAP
func ConvertMap(m [][]byte) (int, int, []byte) {
	var g, p int
	var newMap []byte
	var r, c = len(m), len(m[0])
	newMap = make([]byte, r*(c+2))
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if m[i][j] <= 121 {
				newMap[i*(c+2)+j] = 'X'
				g++
			} else if m[i][j] <= 126 {
				newMap[i*(c+2)+j] = '.'
			} else {
				newMap[i*(c+2)+j] = '1'
				g++
				p++
			}
		}
		newMap[i*(c+2)+c] = '\r'
		newMap[i*(c+2)+c+1] = '\n'
	}
	return g, p, newMap
}
