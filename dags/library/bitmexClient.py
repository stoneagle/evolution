from library import conf, console
import json
import time
import hashlib
import hmac
import requests
from urllib.parse import urlparse
from urllib import parse


class Client(object):
    host = None
    route = None

    def __init__(self, route):
        self.host = conf.BITMEX_HOST
        self.route = route
        return

    def get(self, params):
        data = bytes(parse.urlencode(params), encoding='utf8')
        url = self.route + '?' + data.decode('utf8')
        return self._operate("GET", url, '')

    def post(self, content):
        data = json.dumps(content)
        return self._operate("POST", self.route, data)

    def put(self, content):
        data = json.dumps(content)
        return self._operate("PUT", self.route, data)

    def delete(self, content):
        data = json.dumps(content)
        return self._operate("DELETE", self.route, data)

    def _operate(self, method, route, data):
        headers = self.gen_header(method, route, data)
        path = self.host + route
        headers["content-type"] = "application/json"
        if method == "GET":
            r = requests.get(path, data=data, headers=headers)
        elif method == "POST":
            r = requests.post(path, data=data, headers=headers)
        elif method == "PUT":
            r = requests.put(path, data=data, headers=headers)
        elif method == "DELETE":
            r = requests.delete(path, data=data, headers=headers)
        data_str = r.text
        remain = r.headers["X-RateLimit-Remaining"]
        if int(remain) < 50:
            console.write_msg("bitmex remain num " + remain + " < 50")
        data_json = json.loads(data_str)
        return data_json

    def gen_header(self, method, url, body):
        expires = self.gen_nonce()
        headers = {}
        headers['api-expires'] = str(expires)
        headers['api-key'] = conf.BITMEX_APIKEY
        headers['api-signature'] = self.gen_sign(conf.BITMEX_APISECRET, method, url, expires, body)
        return headers

    def gen_nonce(self):
        return int(round(time.time() * 1000))

    def gen_sign(self, secret, verb, url, nonce, data):
        parsedURL = urlparse(url)
        path = parsedURL.path
        if parsedURL.query:
            path = path + '?' + parsedURL.query
        # 清除data中的空格
        detail = verb + path + str(nonce) + data
        message = bytes(detail, 'utf-8')
        return hmac.new(bytes(secret, 'utf-8'), message, digestmod=hashlib.sha256).hexdigest()
