package Strutture

type Utenti struct {
	Nome     string `json:"Nome"`
	Cognome  string `json:"Cognome"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Id_tg    string `json:"Id_tg"`
	Active   bool   `json:"Active"`
}

type Routes struct {
	Nome         string `json:"Nome"`
	Partenza     string `json:"Partenza"`
	Destinazione string `json:"Destinazione"`
	Email        string `json:"Email"`
}

type Authentication struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Id_tg    string `json:"Id_tg"`
}
