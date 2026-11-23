#!/usr/bin/env python3
"""Receive checksum-verified iPhone reports using a pinned digest or a one-time token."""
import argparse
import hashlib
import hmac
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
from pathlib import Path
import re
import threading


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--host', required=True)
    parser.add_argument('--port', type=int, default=8766)
    parser.add_argument('--out', type=Path, default=Path('experiments/iphone/raw'))
    authorization = parser.add_mutually_exclusive_group(required=True)
    authorization.add_argument('--expected-sha256')
    authorization.add_argument('--token-file', type=Path)
    parser.add_argument('--reports', type=int, default=1)
    parser.add_argument('--allow-interrupted', action='store_true', help='Diagnostic import only; analysis still rejects interrupted runs')
    args = parser.parse_args()
    if args.reports < 1: parser.error('--reports must be positive')
    token = args.token_file.read_text().strip() if args.token_file else None
    if token is not None and len(token) < 16: parser.error('token must contain at least 16 characters')
    args.out.mkdir(parents=True, exist_ok=True)
    received = set()
    lock = threading.Lock()

    class Handler(BaseHTTPRequestHandler):
        def do_GET(self):
            self.send_response(200); self.end_headers(); self.wfile.write(b'{"ready":true}')

        def do_POST(self):
            if self.path != '/upload':
                self.send_error(404); return
            if token is not None and not hmac.compare_digest(self.headers.get('Authorization', ''), 'Bearer '+token):
                self.send_error(403); return
            try: count = int(self.headers.get('Content-Length', '0'))
            except ValueError:
                self.send_error(400); return
            if not 0 < count <= 4*1024*1024:
                self.send_error(413); return
            raw = self.rfile.read(count)
            digest = hashlib.sha256(raw).hexdigest()
            expected = args.expected_sha256 or self.headers.get('X-PVote-SHA256', '')
            if not hmac.compare_digest(digest, expected):
                self.send_error(400, 'SHA-256 mismatch'); return
            try: report = json.loads(raw)
            except (ValueError, UnicodeDecodeError):
                self.send_error(400, 'Invalid JSON'); return
            if not isinstance(report, dict):
                self.send_error(400, 'Invalid report'); return
            run = report.get('runId', '')
            if not re.fullmatch(r'[A-Za-z0-9_.:-]{1,100}', run):
                self.send_error(400, 'Invalid run ID'); return
            meta = report.get('metadata', {})
            statuses = {'complete', 'interrupted'} if args.allow_interrupted else {'complete'}
            if meta.get('isSimulator') is not False or meta.get('buildConfiguration') != 'Release' or report.get('status') not in statuses:
                self.send_error(400, 'Requires a completed physical-device Release report'); return
            destination = args.out/f'benchmark-{run}.json'
            with lock:
                if destination.exists() and hashlib.sha256(destination.read_bytes()).hexdigest() != digest:
                    self.send_error(409, 'Refusing to overwrite a different report'); return
                temporary = destination.with_suffix('.json.partial')
                temporary.write_bytes(raw); temporary.replace(destination)
                received.add(run)
                done = len(received) >= args.reports
            body = json.dumps({'saved': destination.name, 'sha256': digest, 'bytes': len(raw)}).encode()
            self.send_response(200); self.send_header('Content-Type', 'application/json')
            self.send_header('Content-Length', str(len(body))); self.end_headers(); self.wfile.write(body)
            print(f'SAVED {destination} samples={len(report["samples"])} sha256={digest}', flush=True)
            if done: threading.Thread(target=self.server.shutdown, daemon=True).start()

    server = ThreadingHTTPServer((args.host, args.port), Handler)
    print(f'READY http://{args.host}:{args.port}/upload', flush=True)
    try: server.serve_forever()
    finally: server.server_close()


if __name__ == '__main__': main()
