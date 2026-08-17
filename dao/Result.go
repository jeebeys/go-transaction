package dao

var (
	SUCCESS = Result{result: true}
	FAILURE = Result{result: false}
)

// Result 暂时无用
type Result struct {
	result bool
	object interface{}
}

func (r Result) Object(object interface{}) Result {
	r.object = object
	return r
}

func (r Result) Result() (bool, interface{}) {
	return r.result, r.object
}
