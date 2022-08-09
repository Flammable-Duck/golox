package interpreter

type RuntimeException struct {
    errors []string
}

func NewRuntimeException(msg string) RuntimeException{
    return RuntimeException{}.Add(msg)
}

func (r RuntimeException) Add(msg string) RuntimeException{
    return RuntimeException{errors: append(r.errors, msg)}
}

func (r RuntimeException) Error() string {
    var errstring string
    for _ , msg := range r.errors {
        errstring = msg + "\n" + errstring
    }
    return "Runtime exception: " + errstring
}
