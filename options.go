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
	showReq     bool     // dumps requests, or not
	showRsp     bool     // dumps responses, or not
	showBody    bool     // dumps body data in requests or responses, or not
	showHeaders bool     // dumps headers in requests or responses, or not
	showCookies bool     // keeps `cookie` in dumped headers, or not
	showRaw     bool     // dumps body data in requests as the raw format, or makes the dumped data formatted and sorted
	cb          Callback // customized consumer callback function
}

type Option func(opt *options)

// WithShowReq is set to true, it will dump the header and body data of
// requests; when set to false, dumping will not occur.
func WithShowReq(show bool) Option {
	return func(opt *options) {
		opt.showReq = show
	}
}

// WithShowRsp is set to true, it will dump the header and body data of
// responses; when set to false, dumping will not occur.
func WithShowRsp(show bool) Option {
	return func(opt *options) {
		opt.showRsp = show
	}
}

// WithShowBody is set to true, it will dump the body data of requests or
// responses; when set to false, dumping will not occur.
//
// However, if both ShowReq and ShowRsp are set to false, the body data of
// requests or responses will not be dumped regardless of the WithShowBody
// setting.
func WithShowBody(show bool) Option {
	return func(opt *options) {
		opt.showBody = show
	}
}

// WithShowHeaders is set to true, it will dump the header of requests or
// responses; when set to false, dumping will not occur.
//
// However, if both ShowReq and ShowRsp are set to false, the header of
// requests or responses will not be dumped regardless of the WithShowHeaders
// settings.
func WithShowHeaders(show bool) Option {
	return func(opt *options) {
		opt.showHeaders = show
	}
}

// WithShowCookies is set to true, it will keep cookies in header for the
// dumping; when set to false, cookies will not show in the dumped headers.
func WithShowCookies(show bool) Option {
	return func(opt *options) {
		opt.showCookies = show
	}
}

// WithShowRaw is set to true, the raw format of the request body data will be
// preserved; when set to false, the dumped request body data will be formatted
// and sorted.
func WithShowRaw(show bool) Option {
	return func(opt *options) {
		opt.showRaw = show
	}
}

// WithCallback accepts your customized callback method to which the dumped
// string information will be passed, allowing you to consume the dumped
// information in your own way.
func WithCallback(cb Callback) Option {
	return func(opt *options) {
		opt.cb = cb
	}
}
