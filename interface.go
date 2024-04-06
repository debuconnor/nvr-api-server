package roomapi

type Parser interface {
	Parse(string)
}

type Browser interface {
	Scrape()
	Retrieve()
}

type Product interface {
	Get()
	Save()
	Delete()
	Update()
	Browser
}

type User Product

type Website interface {
	Product
	Parser
}
