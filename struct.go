package roomapi

type Admin struct {
	Id       int
	UserId   string
	Password string
}

type Platform struct {
	Code    string
	Admin   User
	Name    string
	Session map[string]string
	Places  []Place
	Url     string
}

type Place struct {
	Id          int
	Admin       User
	Platform    Website
	Rooms       []Room
	Name        string
	Address     string
	Description string
	Status      Status
	Url         string
}

type Room struct {
	Id          int
	Place       Website
	Name        string
	Price       float64
	Description string
	Status      Status
	Url         string
}

type Reservation struct {
	Id          int
	Admin       User
	Platform    Website
	Customer    User
	Room        Website
	Payment     Payment
	Date        string
	SpendTime   int
	PersonCount int
	Status      Status
	Memo        string
	Url         string
	CreatedAt   string
	UpdatedAt   string
}

type Customer struct {
	Id        int
	Name      string
	Phone     string
	Email     string
	BlackList Status
}

type Payment struct {
	Id         int
	Admin      User
	Platform   Website
	Customer   User
	Amount     float64
	PaidAmount float64
	PaidPoint  float64
	Status     int
	CreatedAt  string
	UpdatedAt  string
}

type Status struct {
	Id   int
	Name string
}
