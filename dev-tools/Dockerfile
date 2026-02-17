FROM python:3.14-slim

ARG PORT=3002
ENV PORT=${PORT}

WORKDIR /app
COPY server.py .

CMD ["python", "-u", "server.py"]
