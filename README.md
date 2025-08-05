# Mensaje Servicio

Servicio de mensajes tipo Twitter desarrollado en Go.

## Requisitos

- Go 1.22 o superior

## Instalación

```bash
# Clonar el repositorio
git clone <repository-url>
cd mensajeServicio

# Instalar dependencias
go mod tidy
```

## Desarrollo Local

### Iniciar el servicio

```bash
# Compilar y ejecutar
go run main.go

# O compilar y ejecutar el binario
go build -o main .
./main
```

El servicio estará disponible en `http://localhost:80` (puerto por defecto).

### Variables de entorno

```bash
# Configuración opcional
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

### Ejecutar tests de un paquete específico

```bash
go test ./message-api/controller -v
go test ./message-api/service -v
go test ./components/config -v
```

## Endpoints

- `POST /messages` - Crear mensaje
- `GET /messages` - Obtener mensajes del usuario
- `POST /follows` - Seguir usuario
- `GET /timeline` - Obtener timeline del usuario 
