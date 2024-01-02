import json
import threading

from flask import Flask, request, jsonify
import requests
from threading import Thread
from kafka.consumer_kafka import KafkaConsumer

app = Flask(__name__)
consumer = None
host_auth = "mainserver"
port_auth = "25536"


def thread_handle_login(data, shared_variable, route_list):
    print(f"Received data as dictionary: {data}")

    email = data['Email']
    payload = json.dumps(data)
    try:
        resp = requests.post("http://" + host_auth + ":" + port_auth + "/api/v1/authentication", payload,
                             headers={"Content-Type": "application/json"})
        print(resp.text)
        shared_variable.append(resp)
        if resp.status_code == 200:
            route_data = {
                "Email": email
            }
            payload = json.dumps(route_data)
            route_request = requests.post("http://" + host_auth + ":" + port_auth + "/api/v1/enableRoute", payload,
                                          headers={"Content-Type": "application/json"})
            routes = route_request.json()

            for route in routes:
                route_list.append(route['Nome'])

            shared_variable.append(route_list)

    except:
        # server di autenticazione offline
        resp = "503"
        shared_variable.append(resp)


def kafka_login(topics, chat_id):
    global consumer
    consumer = KafkaConsumer('kafka:9093', chat_id, topics)
    consumer.start_consumer()


@app.route('/api/send', methods=['POST'])
def handle_login():
    # ricevo i dati in formato json
    data = request.json
    print(data)
    shared_variable = []
    route_list = []
    # imposto l'host e la porta per accedere al server auth
    # lancio un nuovo thread per utente
    thread = Thread(target=thread_handle_login, args=(data, shared_variable, route_list))
    thread.start()
    thread.join()

    with threading.Lock():
        response = shared_variable[0] if shared_variable else None

        if response == "503":
            return jsonify(response)

        elif response.text == "Authorized":
            kafka_login(shared_variable[1], data['Id_tg'])

            return jsonify(response.text)

        else:
            return jsonify(response.text)


def kafka_logout():
    global consumer
    consumer.logout_command()


@app.route('/api/logout', methods=['POST'])
def handle_logout():
    print(request.data)
    data = request.json
    payload = json.dumps(data)

    route = requests.post("http://" + host_auth + ":" + port_auth + "/api/v1/disableRoute", payload,
                          headers={"Content-Type": "application/json"})
    if route.text == 'ok':
        kafka_logout()
        return 'logout effettuato'

    else:
        return 'logout non effettuato'


if __name__ == '__main__':
    flask_host = "server_flask"
    flask_port = 8888
    app.run(flask_host, flask_port)
