package siesa

type PopappClient struct {
	Email   string      `json:"email"`
	ID      interface{} `json:"id"`
	Nombre  string      `json:"nombre"`
	Telefono string      `json:"telefono"`
}

type PopappProduct struct {
	ID          interface{} `json:"idProducto"`
	Nombre      string `json:"nombre"`
	PrecioUnitario int    `json:"precioUnitario"`
}

type PopappModifier struct {
	Cantidad int           `json:"cantidad"`
	Producto PopappProduct `json:"producto"`
}

type PopappItemGroup struct {
	ID        interface{}          `json:"id"`
	Modifiers []PopappModifier `json:"modifiers"`
	Nombre    string          `json:"nombre"`
	TodoElegido bool            `json:"todoElegido"`
}

type PopappItem struct {
	Producto        PopappProduct    `json:"producto"`
	TotalModifiers int               `json:"totalModifiers"`
	Observacion     string            `json:"observacion"`
	Categoria       struct {
		NombreCategoria string `json:"nombreCategoria"`
		IDCategoria     string `json:"idCategoria"`
	} `json:"categoria"`
	ItemGroups []PopappItemGroup `json:"itemGroups"`
	Cantidad   int               `json:"cantidad"`
}

type PopappTotal struct {
	DescuentosRecargos []interface{} `json:"descuentosRecargos"`
	Moneda            string        `json:"moneda"`
	FormatedTotal     string        `json:"formatedTotal"`
	TotalItems        int           `json:"totalItems"`
	Total             int           `json:"total"`
	CostoEnvio        *int          `json:"costoEnvio,omitempty"`
}

type PopappAddress struct {
	Label string `json:"label"`
}

type PopappOrder struct {
	EventType     string         `json:"eventType"`
	ID            string   		 `json:"id"`
	DisplayID     string         `json:"displayId"`
	Cliente       PopappClient   `json:"cliente"`
	PreOrder      bool           `json:"preOrder"`
	Direccion     PopappAddress  `json:"direccion"`
	Plataforma    string         `json:"plataforma"`
	Tipo          string         `json:"tipo"`
	IDStore       string         `json:"idStore"`
	Total         PopappTotal    `json:"total"`
	Items         []PopappItem   `json:"items"`
	TipoPago      string         `json:"tipoPago"`
	EstadoInterno string         `json:"estadoInterno"`
	NombreStore   string         `json:"nombreStore"`
	FechaCreacion string   		 `json:"fechaCreacion"`
	FechaPickUp   string         `json:"fechaPickUp"`
	Estado        string         `json:"estado"`
	KeyLocal      string         `json:"keyLocal"`
	Pagado        bool           `json:"pagado"`
}
