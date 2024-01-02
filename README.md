# MapsNotify

Ancora il client desktop non è stato sviluppato quindi per registrare gli utenti abbiamo utilizzato postman con la seguente richiesta API: 
localhost:25536/api/v1/register con body 
{"Nome":"A","Cognome":"A","Email":"A","Password":"A"} 

e un inserimento di una route:
127.0.0.1:25536/api/v1/registerRoute  con body 
{
        "Partenza":"Catania",
        "Destinazione":"Palermo",
        "Email":"A"
}
 
successivamente è possibile mediante il bot telegram effettuare l'accesso al servizio inserendo l'email e la password al seguente bot @Traffic_detection_bot.
Inoltre quando si avviano i server mediante utilizzo del docker compose sono state impostate delle wait di 2 minuti per permettere l'avvio dei server mysql.
