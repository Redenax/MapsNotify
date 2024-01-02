import json
import requests

host = 'server_flask'
port = '8888'


class ConnectionServer:
    def __init__(self, email, psw, chat_id):
        self.email = email
        self.psw = psw
        self.chat_id = str(chat_id)

    def connection_to_server(self):
        # imposto url per mandare i dati inseriti dal'utente al server flask
        url = "http://" + host + ":" + port + "/api/send"
        user_data = {
            "Email": self.email,
            "Password": self.psw,
            "Id_tg": self.chat_id
        }
        # i dati vengono salvati in un dizionario il quale viene trasformato in json
        message = json.dumps(user_data)
        response = requests.post(url, message, headers={"Content-Type": "application/json"})
        print(response.text)
        if response.text != "503":
            resp = json.loads(response.text)
            return resp
        else:
            resp = json.loads(response)
            return resp

    def logout(self):
        url = "http://" + host + ":" + port + "/api/logout"
        data = {
            "email": self.email
        }
        message = json.dumps(data)
        response = requests.post(url, message, headers={"Content-Type": "application/json"})
        print(response.text)
        if response.text != "503":
            return response.text
        else:
            return response.text

