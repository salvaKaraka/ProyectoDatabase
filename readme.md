# Consultas a base de datos usando lenguaje natural

Funciona con LangChain para conectar gpt-4o con una base de datos relacional.
El agente resultante toma consultas en lenguaje natural y responde de la misma manera, utilizando la base de datos.

# Funcionamiento de la API

La api provee 4 endpoints principales, de los cuales
GET /webhook/whatsapp -> es utilizado por whatsapp para poder corroborar el token
POST /webhook/whatsapp -> es utilizado por whatsapp para mandar los eventos a los cuales se suscribió la cuenta
POST /webhook/slack -> es utilizado por slack para corroborar el webhook al cual va a mandar los eventos como para mandar los mismos
POST /webhook/bot -> es utilizado por el proceso aparte que actua como bot de consultas en NQL

# Funcionamiento del flujo

El servidor va a estar escuchando en los 4 endpoints correspondientes, para poder comunicarse con el otro proceso va a mandar un json al estilo de

type requestPayload struct {
Messanger string `json:"messanger"`
Recipient string `json:"recipient"`
Message string `json:"message"`
SessionID string `json:"session_id"`
}

El cual luego de hacer la primera iteración el bot responde con un post con

type requestPayload struct {
Messanger string `json:"messanger"`
Recipient string `json:"recipient"`
Message string `json:"message"`
}

Donde messanger es si es whatsapp o slack, recipient es si es un numero, un canal, etc. y message el texto que manda el bot
