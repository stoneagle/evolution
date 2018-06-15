namespace go engine 
namespace py engine 

enum AssetType {
    Stock = 1,
    Exchange = 2
}

enum ResponseState {
    StateOk = 0,
    StateErrorBusiness = 1,
    StateErrorEmpty = 2
}

struct Response {
    1: ResponseState code
    2: string desc
    3: string data
}

service EngineService {
    Response getType(1: AssetType assetType);
    Response getStrategy(1: string stype);
    Response getClassify(1: AssetType assetType, 2: string ctype, 3: string source, 4: string sub);
}
