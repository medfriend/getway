package dto

type TrazaDTO struct {
	UsuarioID     string `json:"usuario_id"`
	Accion        string `json:"accion"`
	Estado        int    `json:"estado"`
	Endpoint      string `json:"endpoint"`
	Payload       string `json:"payload"`
	Ip            string `json:"ip"`
	Error         string `json:"error"`
	Duracion      string `json:"duracion"`
	Coleccion     string `json:"coleccion"`
	Microservicio string `json:"microservicio"`
	CollectionId  string `json:"collection_id"`
}
