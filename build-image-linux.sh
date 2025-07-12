#!/bin/bash

# Variables
IMAGE_NAME="get-go"
IMAGE_TAG="lastest"
REGISTRY="medfriend1710/medfriend"
REMOTE_TAG="getway-go"
DOCKERFILE_NAME="Dockerfile"  # Nombre del Dockerfile
DOCKERFILE_DIR="."           # Directorio donde está el Dockerfile

# Construir la imagen usando el Dockerfile especificado
echo "🚀 Construyendo imagen Docker: ${IMAGE_NAME}:${IMAGE_TAG} con ${DOCKERFILE_NAME} ..."
docker build -f ${DOCKERFILE_DIR}/${DOCKERFILE_NAME} -t ${IMAGE_NAME}:${IMAGE_TAG} ${DOCKERFILE_DIR}

# Verificar si la construcción fue exitosa
if [ $? -eq 0 ]; then
    echo "✅ Imagen ${IMAGE_NAME}:${IMAGE_TAG} creada correctamente."
else
    echo "❌ Error al construir la imagen Docker."
    exit 1
fi

# Etiquetar la imagen para el registry
echo "🏷️ Etiquetando la imagen como ${REGISTRY}:${REMOTE_TAG} ..."
docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}:${REMOTE_TAG}

# Hacer push al registry
echo "📤 Haciendo push al registry: ${REGISTRY}:${REMOTE_TAG} ..."
docker push ${REGISTRY}:${REMOTE_TAG}

# Verificar si el push fue exitoso
if [ $? -eq 0 ]; then
    echo "✅ Imagen subida correctamente a ${REGISTRY}:${REMOTE_TAG}."
else
    echo "❌ Error al subir la imagen al registry."
    exit 1
fi
