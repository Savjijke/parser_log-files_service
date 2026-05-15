package csvparser

type Node struct {
	NodeDesc        string `csv:"NodeDesc"`
	NumPorts        int    `csv:"NumPorts"`
	NodeType        int    `csv:"NodeType"`
	ClassVersion    int    `csv:"ClassVersion"`
	BaseVersion     int    `csv:"BaseVersion"`
	SystemImageGUID string `csv:"SystemImageGUID"`
	NodeGUID        string `csv:"NodeGUID"`
	PortGUID        string `csv:"PortGUID"`
}

type Port struct {
	NodeGUID      string `csv:"NodeGuid"`
	PortGUID      string `csv:"PortGuid"`
	PortNum       int    `csv:"PortNum"`
	PortState     int    `csv:"PortState"`
	PortPhyState  int    `csv:"PortPhyState"`
	LinkSpeedActv int    `csv:"LinkSpeedActv"`
	LinkWidthActv int    `csv:"LinkWidthActv"`
}


type SystemGeneralInformation struct {
	NodeGuid     string `csv:"NodeGuid"`
	SerialNumber string `csv:"SerialNumber"`
	PartNumber   string `csv:"PartNumber"`
	Revision     string `csv:"Revision"`
	ProductName  string `csv:"ProductName"`
}

type SwitchSettings struct {
	NodeGUID string

	Endianness             int
	EnableEndiannessPerJob int
	ReproducibilityDisable int
}


