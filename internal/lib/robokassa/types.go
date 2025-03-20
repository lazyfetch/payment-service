package robokassa

// Payload may have errors cuz i haven't access to API Robokassa
// Im not enjoy to register my personal data in their service, sry >:(
type JWT struct {
	MerchantLogin  string  `json:"MerchantLogin"`
	InvoiceType    string  `json:"InvoiceType"`
	OutSum         float64 `json:"OutSum"`
	ShpUsername    string  `json:"Shp_username"`
	ShpUserID      string  `json:"Shp_userid"`
	ShpDescription string  `json:"Shp_description"`
}

// Response need for validate request data so naming shitted shit anywayy im corect it soon
type Response struct {
}
