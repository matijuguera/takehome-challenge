package houseresponse

type Houses struct {
	Id        int
	Address   string
	Homeowner string
	Price     int
	PhotoURL  string
}

type HouseResponse struct {
	Houses []Houses
	Ok     bool
}
