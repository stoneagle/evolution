from urllib.parse import urlparse
from library import conf
import json
import time
import hashlib
import hmac
from urllib import parse
try:
    from urllib.request import urlopen, Request
except ImportError:
    from urllib2 import urlopen, Request


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
        path = self.host + url
        headers = self.gen_header('GET', url, '')
        request = Request(path, headers=headers)
        result = urlopen(request, timeout=10)
        data_str = result.read()
        data_str = data_str.decode('GBK')
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
        detail = verb + path + str(nonce) + data
        message = bytes(detail, 'utf-8')
        return hmac.new(bytes(secret, 'utf-8'), message, digestmod=hashlib.sha256).hexdigest()
