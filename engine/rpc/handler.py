from rpc.engine import ttypes
from source.ts import classify
from library import conf
from controller import obtain
import json
import h5py


class EngineServiceHandler:

    def getType(self, assetType):
        ret = ttypes.Response()
        if assetType == ttypes.AssetType.Stock:
            ret.data = json.dumps({"ashare": {"tushare": ["industry", "concept", "hot"]}})
        elif assetType == ttypes.AssetType.Exchange:
            ret.data = json.dumps({"bitcoin": {"bitmex": ["default"]}})
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "asset type not exist"
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

    def getClassify(self, assetType, ctype, source, sub):
        ret = ttypes.Response()
        if assetType == ttypes.AssetType.Exchange:
            if ctype == "bitmex":
                if source == "default":
                    ret.desc = "developed"
                else:
                    ret.code = ttypes.ResponseState.StateErrorBusiness
                    ret.desc = "source not exist"
            else:
                ret.code = ttypes.ResponseState.StateErrorBusiness
                ret.desc = "ctype not exist"
        elif assetType == ttypes.AssetType.Stock:
            if ctype == "ashare":
                if source == "tushare":
                    ret.data = classify.get_json(sub).to_json(orient='records')
                else:
                    ret.code = ttypes.ResponseState.StateErrorBusiness
                    ret.desc = "source not exist"
            else:
                ret.code = ttypes.ResponseState.StateErrorBusiness
                ret.desc = "ctype not exist"
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "source not exist"
        return ret

    def getItem(self, assetType, ctype, source, tag, name):
        ret = ttypes.Response()
        if assetType == ttypes.AssetType.Stock:
            if ctype == "ashare":
                if source == "tushare":
                    ret.data = classify.get_detail(tag, name, retry_count=1, pause=conf.REQUEST_BLANK).to_json(orient='records')
                else:
                    ret.code = ttypes.ResponseState.StateErrorBusiness
                    ret.desc = "source not exist"
            else:
                ret.code = ttypes.ResponseState.StateErrorBusiness
                ret.desc = "ctype not exist"
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "source not exist"
        return ret

    def getItemPoint(self, assetType, ctype, source, code):
        ret = ttypes.Response()
        if assetType == ttypes.AssetType.Stock:
            if ctype == "ashare":
                if source == "tushare":
                    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
                    obtain.code_share(f, code)
                    f.close()
                    ret.data = "true";
                else:
                    ret.code = ttypes.ResponseState.StateErrorBusiness
                    ret.desc = "source not exist"
            else:
                ret.code = ttypes.ResponseState.StateErrorBusiness
                ret.desc = "ctype not exist"
        else:
            ret.code = ttypes.ResponseState.StateErrorBusiness
            ret.desc = "source not exist"
        return ret
