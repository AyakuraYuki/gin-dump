package gin_dump

type Callback func(dumpStr string)

const (
	defaultShowReq     = true
	defaultShowRsp     = true
	defaultShowBody    = true
	defaultShowHeaders = true
	defaultShowCookies = true
	defaultShowRaw     = false
)

type options struct {
	showReq     bool
	showRsp     bool
	showBody    bool
	showHeaders bool
	showCookies bool
	showRaw     bool
	cb          Callback
}

type Option func(opt *options)

func WithShowReq(show bool) Option     { return func(opt *options) { opt.showReq = show } }
func WithShowRsp(show bool) Option     { return func(opt *options) { opt.showRsp = show } }
func WithShowBody(show bool) Option    { return func(opt *options) { opt.showBody = show } }
func WithShowHeaders(show bool) Option { return func(opt *options) { opt.showHeaders = show } }
func WithShowCookies(show bool) Option { return func(opt *options) { opt.showCookies = show } }
func WithShowRaw(show bool) Option     { return func(opt *options) { opt.showRaw = show } }
func WithCallback(cb Callback) Option  { return func(opt *options) { opt.cb = cb } }
