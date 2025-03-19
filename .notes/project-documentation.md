# Documentación del Proyecto Crypto Dashboard

## Descripción General
Este proyecto es un dashboard para monitorear balances de criptomonedas utilizando la API de Buda.com.

## Estructura del Proyecto
```
.
├── cmd/           # Punto de entrada de la aplicación
├── internal/      # Código interno de la aplicación
│   ├── handlers/  # Manejadores HTTP
│   ├── models/    # Modelos de datos
│   ├── services/  # Servicios de negocio
│   └── views/     # Vistas y templates
├── static/        # Archivos estáticos
├── templates/     # Templates HTML
└── .notes/        # Documentación y notas
```

## Componentes Principales

### Servicio Buda
El servicio `Buda` en `internal/services/buda.go` maneja la comunicación con la API de Buda.com:
- Autenticación mediante API key y secret
- Manejo de requests HTTP
- Procesamiento de respuestas
- Gestión de balances

### Handlers
- `DashboardHandler`: Maneja la visualización del dashboard principal
- Muestra los balances de criptomonedas del usuario

### Modelos
- `Balance`: Estructura que representa un balance de criptomoneda
- `DashboardData`: Estructura para los datos del dashboard

## Configuración
El proyecto utiliza variables de entorno para la configuración:
- `BUDA_API_KEY`: API key de Buda.com
- `BUDA_API_SECRET`: API secret de Buda.com
- `BUDA_BASE_URL`: URL base de la API de Buda.com

## Desarrollo
Para ejecutar el proyecto localmente:
1. Copiar `.env.example` a `.env`
2. Configurar las variables de entorno
3. Ejecutar `go run cmd/main.go`

## Docker
El proyecto incluye soporte para Docker:
- `Dockerfile`: Configuración para construir la imagen
- `docker-compose.yml`: Configuración para ejecutar el servicio
- `.dockerignore`: Archivos a ignorar en la construcción

## Próximos Pasos
- [ ] Implementar más endpoints de la API de Buda
- [ ] Mejorar el manejo de errores
- [ ] Agregar tests unitarios y de integración
- [ ] Implementar caché para las respuestas de la API
- [ ] Mejorar la UI/UX del dashboard 