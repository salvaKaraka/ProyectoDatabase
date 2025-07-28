from flask import Flask, request, jsonify
import threading
import requests

app = Flask(__name__)

# Mapa de inputs pendientes
pending_inputs = {}


@app.route("/input", methods=["POST"])
def receive_input():
    data = request.get_json()
    session_id = data.get("session_id")
    message = data.get("message")

    if session_id in pending_inputs:
        pending_inputs[session_id]["message"] = message
        pending_inputs[session_id]["event"].set()  # desbloquea
        return jsonify({"status": "ok"})
    else:
        return jsonify({"error": "no pending input"}), 404


def run_flask():
    app.run(debug=False, use_reloader=False, port=3000)


def input_http(session_id):
    event = threading.Event()
    pending_inputs[session_id] = {"event": event, "message": None}
    print(f"[Esperando input para session: {session_id}]", flush=True)
    event.wait()  # espera sincrónicamente
    msg = pending_inputs[session_id]["message"]
    del pending_inputs[session_id]
    return msg


def send_post(url, data):
    try:
        response = requests.post(url, json=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Error en POST: {e}")
        return None


def main_conversacion():
    print("Esperando nombre...", flush=True)
    nombre = input_http("user123")  # despues tiene que cambiar el mapping del a sesíon
    print(f"[Nombre recibido]: {nombre}", flush=True)

    data = {"message": nombre}
    url = "http://localhost:5000/webhook/bot"
    send_post(url, data)

    print("Esperando edad...", flush=True)
    edad = input_http("user123")
    print(f"[Edad recibida]: {edad}", flush=True)


if __name__ == "__main__":
    # Flask en hilo separado
    flask_thread = threading.Thread(target=run_flask)
    flask_thread.start()

    # Lógica principal
    main_conversacion()
