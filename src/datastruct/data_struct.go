package datastruct

// Conf is for reading the configuration file.
type Conf struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	BackupPath string `json:"backupPath"`
	MailFrom   string `json:"mailFrom"`
	Host       string `json:"host"`
	MailTo     string `json:"mailTo"`
	Subject    string `json:"subject"`
	MailType   string `json:"mailType"`
}

// WaferMap for reading STIF map.
type WaferMap struct {
	LOT     string
	WAFER   string
	PRODUCT string
	READER  string
	// STEP    []float64
	// FLAT    int
	// REF     []float64
	// TARG    []int
	// TARGBC  int
	// FRST    []int
	// STRP    []int
	// PRQUAD  int
	// COQUAD  int
	// NULBC   int
	// GOODS   int
	DATE []string
	DIM  []int
	MAP  [][]byte
}

// ListLine is for each line of wafer in the LIS file.
type ListLine struct {
	MAPPINGFILENAME string
	WAFERID         string
	GOOD            int
	YIELD           float64
}

// ListFile is for the LIS file.
type ListFile struct {
	LOT        string
	WAFERCOUNT int
	GROSS      int
	TOTALPASS  int
	LISTLINE   []ListLine
}

// MapContent contains the name and contents of a map file.
type MapContent struct {
	NAME    string
	CONTENT []byte
}
