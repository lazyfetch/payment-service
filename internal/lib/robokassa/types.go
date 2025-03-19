package robokassa

type Payload struct {
	MerchantLogin  string  `json:"MerchantLogin"`
	InvoiceType    string  `json:"InvoiceType"`
	OutSum         float64 `json:"OutSum"`
	ShpUsername    string  `json:"Shp_username"`
	ShpUserID      string  `json:"Shp_userid"`
	ShpDescription string  `json:"Shp_description"`
}
