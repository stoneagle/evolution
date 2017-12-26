from urllib.parse import urlparse
from library import conf
import time
import hashlib
import hmac


def bitmex_header(method, url, body):
    expires = gen_bitmex_nonce()
    headers = {}
    headers['api-expires'] = str(expires)
    headers['api-key'] = conf.BITMEX_APIKEY
    headers['api-signature'] = gen_bitmex_sign(conf.BITMEX_APISECRET, method, url, expires, body)
    return headers


def gen_bitmex_nonce():
    return int(round(time.time() * 1000))


def gen_bitmex_sign(secret, verb, url, nonce, data):
    parsedURL = urlparse(url)
    path = parsedURL.path
    if parsedURL.query:
        path = path + '?' + parsedURL.query
    detail = verb + path + str(nonce) + data
    message = bytes(detail, 'utf-8')
    return hmac.new(bytes(secret, 'utf-8'), message, digestmod=hashlib.sha256).hexdigest()
