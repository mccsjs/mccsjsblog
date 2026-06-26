import http.server
import socketserver
import urllib.request
import json

class ProxyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path.startswith("/api/"):
            # 代理 API 请求到后端
            try:
                proxy_url = f"http://localhost:8080{self.path}"
                req = urllib.request.Request(proxy_url)
                with urllib.request.urlopen(req) as response:
                    self.send_response(response.status)
                    for header, value in response.getheaders():
                        self.send_header(header, value)
                    self.end_headers()
                    self.wfile.write(response.read())
            except Exception as e:
                self.send_error(500, str(e))
        else:
            # 静态文件请求
            super().do_GET()

if __name__ == "__main__":
    with socketserver.TCPServer(("0.0.0.0", 3100), ProxyHTTPRequestHandler) as httpd:
        print(f"Server running on port 3100")
        httpd.serve_forever()
