# Mensaje Servicio

Servicio de mensajes tipo Twitter desarrollado en Go.

## Requisitos

- Go 1.22 o superior

## Instalación
```
go mod tidy
```

## Desarrollo Local

### Iniciar el servicio

```
go build -o main .
./main
```

El servicio estará disponible en `http://localhost:80` (puerto por defecto).

### Variables de entorno

```bash
export PORT=8080
export AWS_REGION=us-east-1
export DDB_TABLE_MENSAJES=mensajes
export DDB_TABLE_SEGUIDORES=seguidores
export DDB_TABLE_TIMELINE=timeline
```

## Testing

### Ejecutar todos los tests

```bash
go test ./... -v
```

### Ejecutar tests con coverage

```bash
go test ./... -cover
```

## Endpoints

- `POST /message` - Crear mensaje
- `GET /message` - Obtener mensajes del usuario
- `POST /follow` - Seguir usuario
- `GET /timeline` - Obtener timeline del usuario 
