package parsestif

import (
	"mapconvertor/datastruct"
	"strconv"
	"strings"
)

// Parser is for parsing STIF map.
func Parser(mapContent string, w *datastruct.WaferMap) error {
	var n int = 0
	var m []string
	var err error
	// var f float64
	var v int
	w.DATE = []string{}
	w.DIM = []int{}
	s := strings.Split(mapContent, "\r\n")
	for i := 0; i < len(s); i++ {
		m = strings.Split(s[i], "\t")
		if len(m) == 1 {
			m = strings.Split(s[i], "=")
		}
		switch m[0] {
		case "LOT":
			w.LOT = m[1]
		case "WAFER":
			w.WAFER = m[1]
		case "PRODUCT":
			w.PRODUCT = m[1]
		case "READER":
			w.READER = m[1]
		// case "XSTEP", "YSTEP":
		// 	f, err = strconv.ParseFloat(m[1], 64)
		// 	w.STEP = append(w.STEP, f)
		// case "FLAT":
		// 	w.FLAT, err = strconv.Atoi(m[1])
		// case "XREF", "YREF":
		// 	f, err = strconv.ParseFloat(m[1], 64)
		// 	w.REF = append(w.REF, f)
		// case "XBE TARG1", "YBE TARG1", "XBE TARG2", "YBE TARG2", "XBE TARG3", "YBE TARG3":
		// 	v, err = strconv.Atoi(m[1])
		// 	w.TARG = append(w.TARG, v)
		// case "TARGBC":
		// 	w.TARGBC, err = strconv.Atoi(m[1])
		// case "XFRST", "YFRST":
		// 	v, err = strconv.Atoi(m[1])
		// 	w.FRST = append(w.FRST, v)
		// case "XSTRP", "YSTRP":
		// 	v, err = strconv.Atoi(m[1])
		// 	w.STRP = append(w.STRP, v)
		// case "PRQUAD":
		// 	w.PRQUAD, err = strconv.Atoi(m[1])
		// case "COQUAD":
		// 	w.COQUAD, err = strconv.Atoi(m[1])
		// case "NULBC":
		// 	w.NULBC, err = strconv.Atoi(m[1])
		// case "GOODS":
		// 	w.GOODS, err = strconv.Atoi(m[1])
		case "DATE", "TIME", "EDATE", "ETIME":
			w.DATE = append(w.DATE, m[1])
		case "WMXDIM", "WMYDIM":
			n = i
			v, err = strconv.Atoi(m[1])
			w.DIM = append(w.DIM, v)
		}
	}
	w.MAP = make([][]byte, w.DIM[1])
	for i := 0; i < w.DIM[1]; i++ {
		w.MAP[i] = []byte(s[n+i+2])
	}
	return err
}
