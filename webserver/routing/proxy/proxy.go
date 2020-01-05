package proxy

import (
	// routes to use for proxy
	_ "cant/webserver/routing/filter"
	_ "cant/webserver/routing/livelog"
	_ "cant/webserver/routing/pgn"
	_ "cant/webserver/routing/spn"
	_ "cant/webserver/routing/stats"
	_ "cant/webserver/routing/xstatic"
)
