package errors

type Kind int

// Please add you errors to the end of this list ( check kind are not duplicated )
const (
	BadRequest                    Kind = 400
	Unauthorized                  Kind = 401
	PaymentRequired               Kind = 402
	Forbidden                     Kind = 403
	NotFound                      Kind = 404
	MethodNotAllowed              Kind = 405
	NotAcceptable                 Kind = 406
	ProxyAuthRequired             Kind = 407
	RequestTimeout                Kind = 408
	Conflict                      Kind = 409
	Gone                          Kind = 410
	LengthRequired                Kind = 411
	PreconditionFailed            Kind = 412
	RequestEntityTooLarge         Kind = 413
	RequestURITooLong             Kind = 414
	UnsupportedMediaType          Kind = 415
	RequestedRangeNotSatisfiable  Kind = 416
	ExpectationFailed             Kind = 417
	MisdirectedRequest            Kind = 421
	UnprocessableEntity           Kind = 422
	Locked                        Kind = 423
	FailedDependency              Kind = 424
	TooEarly                      Kind = 425
	UpgradeRequired               Kind = 426
	PreconditionRequired          Kind = 428
	TooManyRequests               Kind = 429
	RequestHeaderFieldsTooLarge   Kind = 431
	UnavailableForLegalReasons    Kind = 451
	InternalServerError           Kind = 500
	NotImplemented                Kind = 501
	BadGateway                    Kind = 502
	ServiceUnavailable            Kind = 503
	GatewayTimeout                Kind = 504
	HTTPVersionNotSupported       Kind = 505
	VariantAlsoNegotiates         Kind = 506
	InsufficientStorage           Kind = 507
	LoopDetected                  Kind = 508
	NotExtended                   Kind = 510
	NetworkAuthenticationRequired Kind = 511
)

func (k Kind) String() string {
	switch k {
	case BadRequest:
		return "bad request"
	case Unauthorized:
		return "unauthorized"
	case PaymentRequired:
		return "payment required"
	case Forbidden:
		return "forbidden"
	case NotFound:
		return "not found"
	case MethodNotAllowed:
		return "method not allowed"
	case NotAcceptable:
		return "not acceptable"
	case ProxyAuthRequired:
		return "proxy auth required"
	case RequestTimeout:
		return "request timeout"
	case Conflict:
		return "conflict"
	case Gone:
		return "gone"
	case LengthRequired:
		return "length required"
	case PreconditionFailed:
		return "precondition failed"
	case RequestEntityTooLarge:
		return "request entity too large"
	case RequestURITooLong:
		return "request URI too long"
	case UnsupportedMediaType:
		return "unsupported media type"
	case RequestedRangeNotSatisfiable:
		return "requested range not satisfiable"
	case ExpectationFailed:
		return "expectation failed"
	case MisdirectedRequest:
		return "misdirected request"
	case UnprocessableEntity:
		return "unprocessable entity"
	case Locked:
		return "locked"
	case FailedDependency:
		return "failed dependency"
	case TooEarly:
		return "too early"
	case UpgradeRequired:
		return "upgrade required"
	case PreconditionRequired:
		return "precondition required"
	case TooManyRequests:
		return "too many requests"
	case RequestHeaderFieldsTooLarge:
		return "request header fields too large"
	case UnavailableForLegalReasons:
		return "unavailable for legal reasons"
	case InternalServerError:
		return "internal server error"
	case NotImplemented:
		return "not implemented"
	case BadGateway:
		return "bad gateway"
	case ServiceUnavailable:
		return "service unavailable"
	case GatewayTimeout:
		return "gateway timeout"
	case HTTPVersionNotSupported:
		return "HTTP version not supported"
	case VariantAlsoNegotiates:
		return "variant also negotiates"
	case InsufficientStorage:
		return "insufficient storage"
	case LoopDetected:
		return "loop detected"
	case NotExtended:
		return "not extended"
	case NetworkAuthenticationRequired:
		return "network authentication required"
	default:
		return "unknown kind"
	}
}

func (k Kind) GetHttpStatus() int  {
	return int(k)
}

func AsKindContext(err error) (Kind, map[string]interface{}) {
	var target *WithKindContext
	if As(err, &target) {
		return target.Kind(), target.Context()
	}
	return InternalServerError, nil
}

func HasKind(err error, k Kind) bool {
	kind, _ := AsKindContext(err)
	if kind == k {
		return true
	}
	return false
}