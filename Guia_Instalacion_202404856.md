# Guía de Instalación y Configuración
## Proyecto 1 - Sistemas Operativos 1

**Proyecto:** Desarrollo, Conexión y Gestión de Contenedores en Entornos Virtualizados
**Carnet:** Angel Emanuel Rodriguez Corado
**Carnet:** 202404856  
**Curso:** Sistemas Operativos 1  

---

## 1. Descripción del proyecto

El proyecto implementa un entorno virtualizado compuesto por tres máquinas virtuales. Cada máquina utiliza un runtime de contenedores diferente y cumple una función específica.

La distribución implementada es:

| Máquina | IP | Runtime | Servicios |
|---|---|---|---|
| VM1 | `192.168.122.210` | Containerd / nerdctl | API1, API2 |
| VM2 | `192.168.122.56` | Podman | API3 |
| VM3 | `192.168.122.126` | Docker | Zot Registry |

Las APIs fueron desarrolladas en Go y utilizan REST/HTTP y JSON para comunicarse entre ellas.

---

## 2. Arquitectura

```text
                         RED VIRTUAL
                              |
          +-------------------+-------------------+
          |                   |                   |
          v                   v                   v

   +-------------+     +-------------+     +-------------+
   |    VM1      |     |    VM2      |     |    VM3      |
   | .122.210    |     | .122.56     |     | .122.126    |
   |             |     |             |     |             |
   | Containerd  |     |   Podman    |     |   Docker    |
   |   /nerdctl  |     |             |     |             |
   |             |     |   API3      |     | Zot :5000   |
   | API1 :8081  |     |   :8083     |     |             |
   | API2 :8082  |     |             |     |             |
   +-------------+     +-------------+     +-------------+
          |                   |                   |
          +-------------------+-------------------+
                              |
                     REST / HTTP / JSON
```

---

# 3. Requisitos previos

Antes de instalar el proyecto se deben tener disponibles:

- Tres máquinas virtuales Linux.
- Conectividad de red entre las tres VMs.
- KVM/QEMU para la virtualización.
- Go.
- Containerd.
- nerdctl en VM1.
- Podman en VM2.
- Docker en VM3.
- Zot Registry en VM3.
- Acceso a GitHub para obtener el código fuente.

La distribución requerida por el proyecto es:

- VM1: Containerd con API1 y API2.
- VM2: Podman con API3.
- VM3: Docker con Zot.

---

# 4. Configuración de red

Las direcciones utilizadas en el proyecto son:

```text
VM1 = 192.168.122.210
VM2 = 192.168.122.56
VM3 = 192.168.122.126
```

Verificar la dirección IP de cada máquina con:

```bash
ip addr
```

También se puede comprobar la conectividad entre las VMs:

```bash
ping -c 4 192.168.122.210
ping -c 4 192.168.122.56
ping -c 4 192.168.122.126
```

Cada VM debe poder comunicarse con las demás.

---

# 5. VM1 - Containerd

## 5.1. Verificar Containerd

En VM1:

```bash
containerd --version
```

Verificar nerdctl:

```bash
nerdctl --version
```

Verificar que Containerd esté ejecutándose:

```bash
sudo systemctl status containerd
```

Si está detenido:

```bash
sudo systemctl start containerd
```

Para habilitarlo al iniciar:

```bash
sudo systemctl enable containerd
```

---

## 5.2. Estructura de las APIs

En VM1 se utiliza la siguiente estructura:

```text
~/VM1/
├── api1/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
│
└── api2/
    ├── Dockerfile
    ├── go.mod
    └── main.go
```

API1 utiliza el puerto `8081`.

API2 utiliza el puerto `8082`.

---

# 6. VM2 - Podman

## 6.1. Verificar Podman

En VM2:

```bash
podman --version
```

Verificar que se pueda ejecutar un contenedor:

```bash
podman ps
```

La estructura de API3 es:

```text
~/VM2/
└── api3/
    ├── Dockerfile
    ├── go.mod
    └── main.go
```

API3 utiliza el puerto `8083`.

---

# 7. VM3 - Docker y Zot Registry

## 7.1. Verificar Docker

En VM3:

```bash
docker --version
```

Verificar el servicio:

```bash
sudo systemctl status docker
```

Si está detenido:

```bash
sudo systemctl start docker
```

---

## 7.2. Directorio de datos de Zot

Se utiliza el directorio:

```text
~/zot-data
```

Este directorio se utiliza como volumen para almacenar los datos de Zot.

---

## 7.3. Ejecutar Zot

Crear el directorio si todavía no existe:

```bash
mkdir -p ~/zot-data
```

Ejecutar Zot:

```bash
docker run -d \
  --name zot \
  -p 5000:5000 \
  --restart=always \
  -v ~/zot-data:/var/lib/zot \
  ghcr.io/project-zot/zot-linux-amd64:latest
```

Verificar que Zot esté ejecutándose:

```bash
docker ps
```

También se puede verificar el catálogo desde cualquier VM:

```bash
curl http://192.168.122.126:5000/v2/_catalog
```

El registro utiliza:

```text
192.168.122.126:5000
```

---

# 8. Configuración de API1

API1 se encuentra en VM1 y utiliza:

```text
IP: 192.168.122.210
Puerto: 8081
VM: VM1
Carnet: 202404856
```

Endpoint principal:

```text
GET /health
```

Endpoints de comunicación:

```text
GET /api1/202404856/call-api2
GET /api1/202404856/call-api3
```

La URL completa es:

```text
http://192.168.122.210:8081/health
http://192.168.122.210:8081/api1/202404856/call-api2
http://192.168.122.210:8081/api1/202404856/call-api3
```

---

# 9. Configuración de API2

API2 se encuentra en VM1 y utiliza:

```text
IP: 192.168.122.210
Puerto: 8082
VM: VM1
Carnet: 202404856
```

Endpoints:

```text
GET /health
GET /api2/202404856/call-api1
GET /api2/202404856/call-api3
```

URLs completas:

```text
http://192.168.122.210:8082/health
http://192.168.122.210:8082/api2/202404856/call-api1
http://192.168.122.210:8082/api2/202404856/call-api3
```

---

# 10. Configuración de API3

API3 se encuentra en VM2:

```text
IP: 192.168.122.56
Puerto: 8083
VM: VM2
Carnet: 202404856
```

Endpoints:

```text
GET /health
GET /api3/202404856/call-api1
GET /api3/202404856/call-api2
```

URLs completas:

```text
http://192.168.122.56:8083/health
http://192.168.122.56:8083/api3/202404856/call-api1
http://192.168.122.56:8083/api3/202404856/call-api2
```

---

# 11. Construcción de imágenes

Cada API tiene su propio Dockerfile.

Los nombres utilizados para las imágenes son:

```text
api1-202404856
api2-202404856
api3-202404856
```

El tag utilizado es:

```text
v1
```

Por lo tanto, las imágenes completas son:

```text
api1-202404856:v1
api2-202404856:v1
api3-202404856:v1
```

El proyecto establece que las imágenes deben utilizar el formato:

```text
API#-#CARNET
```

---

# 12. Registro privado Zot

El registro privado utilizado es:

```text
192.168.122.126:5000
```

Las imágenes almacenadas en Zot son:

```text
192.168.122.126:5000/api1-202404856:v1
192.168.122.126:5000/api2-202404856:v1
192.168.122.126:5000/api3-202404856:v1
```

Verificar el catálogo:

```bash
curl http://192.168.122.126:5000/v2/_catalog
```

La respuesta esperada debe incluir:

```json
{
  "repositories": [
    "api1-202404856",
    "api2-202404856",
    "api3-202404856"
  ]
}
```

Para consultar los tags de API3:

```bash
curl http://192.168.122.126:5000/v2/api3-202404856/tags/list
```

Respuesta esperada:

```json
{
  "name": "api3-202404856",
  "tags": [
    "v1"
  ]
}
```

El mismo procedimiento puede utilizarse para API1 y API2.

---

# 13. Ejecución de las APIs

## VM1 - API1

API1 debe ejecutarse en el puerto:

```text
8081
```

Verificar:

```bash
curl http://192.168.122.210:8081/health
```

## VM1 - API2

API2 debe ejecutarse en el puerto:

```text
8082
```

Verificar:

```bash
curl http://192.168.122.210:8082/health
```

## VM2 - API3

API3 debe ejecutarse en el puerto:

```text
8083
```

Verificar:

```bash
curl http://192.168.122.56:8083/health
```

---

# 14. Pruebas de comunicación entre APIs

## API1

Probar API1 → API2:

```bash
curl http://192.168.122.210:8081/api1/202404856/call-api2
```

Probar API1 → API3:

```bash
curl http://192.168.122.210:8081/api1/202404856/call-api3
```

## API2

Probar API2 → API1:

```bash
curl http://192.168.122.210:8082/api2/202404856/call-api1
```

Probar API2 → API3:

```bash
curl http://192.168.122.210:8082/api2/202404856/call-api3
```

## API3

Probar API3 → API1:

```bash
curl http://192.168.122.56:8083/api3/202404856/call-api1
```

Probar API3 → API2:

```bash
curl http://192.168.122.56:8083/api3/202404856/call-api2
```

---

# 15. Resultado esperado de las llamadas

Cuando una API consultada está funcionando correctamente, la respuesta debe indicar:

```json
{
  "apiname": "API#",
  "message": "The API# located on the VM# is working",
  "connection": true,
  "carnet": "202404856"
}
```

Si una API no está disponible:

```json
{
  "apiname": "API#",
  "message": "ERROR: The API# located on the VM# is not working",
  "connection": false,
  "carnet": "202404856"
}
```

---

# 16. Pruebas funcionales realizadas

Se verificaron los nueve endpoints principales del proyecto.

### VM1 - API1

```text
GET /health                         ✓
GET /api1/202404856/call-api2      ✓
GET /api1/202404856/call-api3      ✓
```

### VM1 - API2

```text
GET /health                         ✓
GET /api2/202404856/call-api1      ✓
GET /api2/202404856/call-api3      ✓
```

### VM2 - API3

```text
GET /health                         ✓
GET /api3/202404856/call-api1      ✓
GET /api3/202404856/call-api2      ✓
```

Las pruebas realizadas mostraron `connection: true` para las seis comunicaciones cruzadas.

---

# 17. Verificación final de la infraestructura

La infraestructura final debe quedar de la siguiente manera:

```text
VM1 - 192.168.122.210
Runtime: Containerd / nerdctl
    |
    +-- API1 :8081
    |
    +-- API2 :8082


VM2 - 192.168.122.56
Runtime: Podman
    |
    +-- API3 :8083


VM3 - 192.168.122.126
Runtime: Docker
    |
    +-- Zot Registry :5000
```

Comunicación REST:

```text
API1 <------> API2
 |             |
 |             |
 +------> API3 <------+
```

Registro de imágenes:

```text
API1 ──┐
API2 ──┼──> Zot Registry
API3 ──┘
```

---

# 18. Comandos rápidos de verificación

### VM1

```bash
sudo systemctl status containerd
nerdctl ps
curl http://192.168.122.210:8081/health
curl http://192.168.122.210:8082/health
```

### VM2

```bash
podman ps
curl http://192.168.122.56:8083/health
```

### VM3

```bash
docker ps
curl http://192.168.122.126:5000/v2/_catalog
```

### Verificar comunicación

```bash
curl http://192.168.122.210:8081/api1/202404856/call-api2
curl http://192.168.122.210:8081/api1/202404856/call-api3

curl http://192.168.122.210:8082/api2/202404856/call-api1
curl http://192.168.122.210:8082/api2/202404856/call-api3

curl http://192.168.122.56:8083/api3/202404856/call-api1
curl http://192.168.122.56:8083/api3/202404856/call-api2
```

---

# 19. Solución de problemas

## Docker no permite ejecutar comandos

Verificar:

```bash
sudo systemctl status docker
```

Si el usuario no tiene permisos sobre Docker:

```bash
sudo usermod -aG docker $USER
```

Después cerrar sesión y volver a ingresar.

---

## Containerd no está ejecutándose

```bash
sudo systemctl status containerd
sudo systemctl start containerd
```

---

## Podman no encuentra un contenedor

Verificar:

```bash
podman ps -a
```

---

## Zot no responde

Verificar:

```bash
docker ps -a
```

Revisar los logs:

```bash
docker logs zot
```

Comprobar el registro:

```bash
curl http://192.168.122.126:5000/v2/_catalog
```

---

## Una API no responde

Comprobar que el contenedor esté ejecutándose:

```bash
nerdctl ps
```

o:

```bash
podman ps
```

Según la VM.

También comprobar directamente el endpoint:

```bash
curl http://IP:PUERTO/health
```

---

# 20. Estado final esperado

Al finalizar la instalación, se debe contar con:

- Tres VMs funcionando.
- VM1 utilizando Containerd/nerdctl.
- VM2 utilizando Podman.
- VM3 utilizando Docker.
- Zot Registry disponible en `192.168.122.126:5000`.
- API1 disponible en `192.168.122.210:8081`.
- API2 disponible en `192.168.122.210:8082`.
- API3 disponible en `192.168.122.56:8083`.
- Las tres APIs desarrolladas en Go.
- Los endpoints `/health` funcionando.
- Comunicación REST entre las tres APIs.
- Imágenes `api1-202404856:v1`, `api2-202404856:v1` y `api3-202404856:v1`.
- Imágenes almacenadas en el registro privado Zot.
- Evidencia de pruebas funcionales.

