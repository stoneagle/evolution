package resp

type ErrorCode int

const (
	ErrorOk          ErrorCode = 0
	ErrorParams      ErrorCode = 1
	ErrorFiles       ErrorCode = 2
	ErrorMysql       ErrorCode = 3
	ErrorRedis       ErrorCode = 4
	ErrorAuth        ErrorCode = 5
	ErrorEngine      ErrorCode = 6
	ErrorDataService ErrorCode = 7
	ErrorServer      ErrorCode = 500
)

type WsStatus int

const (
	WsMessage    WsStatus = 0
	WsConnect    WsStatus = 1
	WsDisconnect WsStatus = 2
	WsError      WsStatus = 3
)
