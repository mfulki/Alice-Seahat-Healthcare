package constant

type wrapCtx string

var (
	UserContext = "authorizedUser"
	TxContext   = wrapCtx("TxContext")
	RequestID   = "RequestID"
)
