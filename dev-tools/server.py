import os
from http.server import BaseHTTPRequestHandler, HTTPServer


class Handler(BaseHTTPRequestHandler):
    def do_GET(self) -> None:
        print(f"GET {self.path} Headers: {self.headers}")
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK")

    def do_POST(self) -> None:
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        print(f"POST {self.path} Headers: {self.headers} Body: {body!r}")
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK")

    def do_PUT(self) -> None:
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        print(f"PUT {self.path} Headers: {self.headers} Body: {body!r}")
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK")

    def do_PATCH(self) -> None:
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        print(f"PATCH {self.path} Headers: {self.headers} Body: {body!r}")
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK")

    def do_DELETE(self) -> None:
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        print(f"DELETE {self.path} Headers: {self.headers} Body: {body!r}")
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK")


if __name__ == "__main__":
    port = int(os.getenv("PORT", 3002))
    HTTPServer(("", port), Handler).serve_forever()
