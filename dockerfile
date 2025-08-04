# Etapa 1: construimos el binario
FROM golang:1.24.5 AS builder

#Asignamos el workdir, en este caso le ponemos api
WORKDIR /api

#Copiamos el sum y el mod y le hacemos el tidy para que instale todo
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Aca importante 
# go build -> construye el binario
# -o indica el nombre del binario (y donde va a estar que por defecto de donde se llama)
# cmd/main.go es que quiero convertir en binario
# compilación sin dependencias a C
RUN CGO_ENABLED=0 GOOS=linux go build -o api_app cmd/main.go 
# Etapa 2: imagen liviana para albergar el binario
FROM gcr.io/distroless/static-debian11

WORKDIR /api

COPY --from=builder /api/api_app .
COPY --from=builder /api/.env .

# Expone el puerto que usa tu app
EXPOSE 5000

CMD ["./api_app"]
