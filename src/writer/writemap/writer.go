package writemap

import (
	"bytes"
	"convertor/stiftomap"
	"datastruct"
	"fmt"
)

// WriteMap is for writing the MAP.
func WriteMap(waferMap datastruct.WaferMap, m *datastruct.MapContent, listFile *datastruct.ListFile, listLine *datastruct.ListLine) {
	var buf bytes.Buffer
	template := "FAB ID                 : %s\r\nFAB PRODUCT ID         : %s\r\nFAB Lot ID             : %s\r\nCP Vendor ID           : %s\r\nCUSTOMER ID            : %s\r\nCUSTOMER PRODUCT ID    : %s\r\nCUSTOMER LOT ID        : %s\r\nWAFER ID               : %s\r\nOCR_ID                 : %s\r\nDATE                   : %s %s\r\nGROSS DIE              : %d\r\nPASS DIE               : %d\r\nINKLESS NOTCH          : DOWN\r\nXINCREAMENT            : RIGHT\r\nYINCREAMENT            : DOWN\r\n\r\n\r\nINKLESS_START\r\n"
	g, p, newMap := stiftomap.ConvertMap(waferMap.MAP)
	m.NAME = fmt.Sprintf("%s-%s.MAP", waferMap.LOT, waferMap.WAFER)
	listFile.LOT = waferMap.LOT
	listFile.WAFERCOUNT++
	listFile.GROSS = g
	listFile.TOTALPASS += p
	listLine.MAPPINGFILENAME = m.NAME
	listLine.WAFERID = waferMap.WAFER
	listLine.GOOD = p
	listLine.YIELD = float64(p) * 100.0 / float64(g)
	// listFile.LISTLINE = append(listFile.LISTLINE, datastruct.ListLine{MAPPINGFILENAME: m.NAME, WAFERID: waferMap.WAFER, GOOD: p, YIELD: float64(p) * 100.0 / float64(g)})
	buf.WriteString(fmt.Sprintf(template, "FAB", waferMap.PRODUCT, waferMap.LOT, "Vendor", "CUSTOMER", "CUSTOMER PRODUCT", "CUSTOMER LOT", waferMap.WAFER, waferMap.READER, waferMap.DATE[0], waferMap.DATE[1], g, p))
	buf.Write(newMap)
	m.CONTENT = buf.Bytes()
	// return datastruct.MapContent{NAME: mapFileName, CONTENT: buf.Bytes()}
}

// WriteList is for writing the LIS of the lot.
func WriteList(m *datastruct.MapContent, listFile datastruct.ListFile) {
	var buf bytes.Buffer
	template := "FAB LOT ID               : %s\r\nCP Vendor ID             : %s\r\nCUSTOMER ID              : Hisilicon\r\nCUSTOMER PRODUCT ID      : %s\r\nCUSTOMER LOT ID          : %s\r\nWAFER COUNT              : %d\r\nGROSS DIE                : %d\r\nTOTAL PASS DIE           : %d\r\n\r\nMAPPING_FILE_NAME    WAFER_ID    GOOD    YIELD"
	line := "\r\n%-21s%-12s%-8d%.2f%%"
	buf.WriteString(fmt.Sprintf(template, listFile.LOT, "Vender", "Customer Product", "Customer Lot", listFile.WAFERCOUNT, listFile.GROSS, listFile.TOTALPASS))
	for _, l := range listFile.LISTLINE {
		if l.MAPPINGFILENAME != "" {
			buf.WriteString(fmt.Sprintf(line, l.MAPPINGFILENAME, l.WAFERID, l.GOOD, l.YIELD))
		}
	}
	m.NAME = fmt.Sprintf("%s.LIS", listFile.LOT)
	m.CONTENT = buf.Bytes()
	// return datastruct.MapContent{NAME: fmt.Sprintf("%s.LIS", listFile.LOT), CONTENT: buf.Bytes()}
}
