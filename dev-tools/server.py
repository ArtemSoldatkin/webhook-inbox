import os
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import parse_qs, urlparse


class Handler(BaseHTTPRequestHandler):

    def _log_message(self, method: str) -> None:
        parsed_url = urlparse(self.path)
        query_params = parse_qs(parsed_url.query)
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        print(
            f"{method} {self.path} Headers: {self.headers} Query Params: {query_params} Body: {body!r}"
        )

    def _set_response(self, status_code: int, text: str) -> None:
        self.send_response(status_code)
        self.send_header("Content-Type", "text/plain")
        self.end_headers()
        self.wfile.write(text.encode("utf-8"))

    def do_GET(self) -> None:
        self._log_message("GET")
        if self.path == "/4xx":
            return self._set_response(status_code=400, text="Bad Request")
        if self.path == "/5xx":
            return self._set_response(status_code=500, text="Internal Server Error")
        return self._set_response(status_code=200, text="OK")

    def do_POST(self) -> None:
        self._log_message("POST")
        if self.path == "/4xx":
            return self._set_response(status_code=400, text="Bad Request")
        if self.path == "/5xx":
            return self._set_response(status_code=500, text="Internal Server Error")
        return self._set_response(status_code=200, text="OK")

    def do_PUT(self) -> None:
        self._log_message("PUT")
        if self.path == "/4xx":
            return self._set_response(status_code=400, text="Bad Request")
        if self.path == "/5xx":
            return self._set_response(status_code=500, text="Internal Server Error")
        return self._set_response(status_code=200, text="OK")

    def do_PATCH(self) -> None:
        self._log_message("PATCH")
        if self.path == "/4xx":
            return self._set_response(status_code=400, text="Bad Request")
        if self.path == "/5xx":
            return self._set_response(status_code=500, text="Internal Server Error")
        return self._set_response(status_code=200, text="OK")

    def do_DELETE(self) -> None:
        self._log_message("DELETE")
        if self.path == "/4xx":
            return self._set_response(status_code=400, text="Bad Request")
        if self.path == "/5xx":
            return self._set_response(status_code=500, text="Internal Server Error")
        return self._set_response(status_code=200, text="OK")


if __name__ == "__main__":
    port = int(os.getenv("PORT", 3002))
    HTTPServer(("", port), Handler).serve_forever()
