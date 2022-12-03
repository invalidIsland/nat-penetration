package define

const (
	KeepAliveStr  = "KeepAlive\n"
	NewConnection = "New Connection\n"

	LocalServerAddr   = ":8088"
	ControlServerAddr = ":8000"
	TunnelServerAddr  = ":8010"
	UserRequestAddr   = ":8020"

	BufSize = 1 << 10
)
