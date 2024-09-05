package dto

type DeleteObjectRequest struct {
	Quiet  bool `xml:"Quiet"`
	Object struct {
		Key string `xml:"Key"`
	} `xml:"Object"`
}

type DeleteResult struct {
	DeletedResult []Deleted `xml:"Deleted"`
}

type Deleted struct {
	Key string `xml:"Key"`
}
