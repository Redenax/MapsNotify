package Strutture

type Utenti struct {
	Nome     string `json:"Nome"`
	Cognome  string `json:"Cognome"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Id_tg    string `json:"Id_tg"`
	Active   bool   `json:"Active"`
}
