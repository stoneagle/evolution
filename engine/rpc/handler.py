from rpc.engine import ttypes
from source.ts import classify
import json


class EngineServiceHandler:

    def getType(self, resource):
        ret = ttypes.Response()
        if resource == ttypes.QuantType.Stock:
            ret.data = json.dumps({"ashare": ["tushare"]})
        elif resource == ttypes.QuantType.Exchange:
            ret.data = json.dumps({"bitcoin": ["bitmex"]})
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "source not exist"
        return ret

    def getStrategy(self, stype):
        ret = ttypes.Response()
        if stype == "bitmex":
            ret.data = json.dumps(["reverse"])
        elif stype == "ashare":
            ret.data = json.dumps(["reverse"])
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "stype not exist"
        return ret

    def getClassify(self, stype, source):
        ret = ttypes.Response()
        if stype == "bitmex":
            if source == "bitmex":
                ret.desc = "developed"
            else:
                ret.code = ttypes.ResponseState.StateErrorBusiness
                ret.desc = "source not exist"
        elif stype == "ashare":
            if source == "tushare":
                ret.data = classify.get_json().to_json(orient='records')
                ret.desc = "developed"
            else:
                ret.code = ttypes.ResponseState.StateErrorBusiness
                ret.desc = "source not exist"
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "stype not exist"
        return ret
