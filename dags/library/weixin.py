from library import conf
import wxpy


class WXBot(object):
    # 初始化的bot机器人
    BotObj = None

    def __init__(self):
        cache_path = conf.WEIXIN_BOT_CACHE_PATH + "/qr_cache"
        self.BotObj = wxpy.Bot(cache_path=cache_path, console_qr=1)
        return

    def sendHelper(self, msg):
        self.BotObj.file_helper.send(msg)
        return
