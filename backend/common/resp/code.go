package resp

type ErrorCode int

const (
	ErrorOk           ErrorCode = 0
	ErrorParams       ErrorCode = 1
	ErrorFiles        ErrorCode = 2
	ErrorDatabase     ErrorCode = 3
	ErrorCache        ErrorCode = 4
	ErrorAuth         ErrorCode = 5
	ErrorEngine       ErrorCode = 6
	ErrorDataTransfer ErrorCode = 7
	ErrorSign         ErrorCode = 8
	ErrorApi          ErrorCode = 9
	ErrorServer       ErrorCode = 500
)

var (
	ErrorMessages map[ErrorCode]string = map[ErrorCode]string{
		ErrorParams:       "params error",
		ErrorFiles:        "file error",
		ErrorDatabase:     "database get error",
		ErrorCache:        "cache get error",
		ErrorAuth:         "auth error",
		ErrorEngine:       "engine error",
		ErrorDataTransfer: "data transfer error",
		ErrorSign:         "sign error",
		ErrorApi:          "api error",
		ErrorServer:       "server exception",
	}
)

type WsStatus int

const (
	WsMessage    WsStatus = 0
	WsConnect    WsStatus = 1
	WsDisconnect WsStatus = 2
	WsError      WsStatus = 3
)
