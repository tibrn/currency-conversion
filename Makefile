version: "3.7"

networks:
  backend:
    name: backend
    driver: bridge

services:
  currency_convertor_web:
    build:
      context: .
      dockerfile: Dockerfile.Test
    env_file: .env.test
    networks:
      - backend

  currency_convertor_jobs:
    image: nats:alpine
    restart: always
    networks:
      - backend

  redis:
    image: redis:alpine
    restart: always
    networks:
      - backend

