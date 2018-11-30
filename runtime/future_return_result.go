package runtime

import (
	"github.com/1-bi/servicebus"
	"github.com/1-bi/uerrors"
	"github.com/vmihailenco/msgpack"
	"reflect"
)

/**
 * --------------------------------- *
 * share object : base baseFutureReturnResult
 * --------------------------------- *
 */
type baseFutureReturnResult struct {
	state   int8
	resErrs map[string]uerrors.CodeError
	resRes  map[string][]byte
}

/**
 * define,  ALL_COMPLETE ,  ANY_ERRORS , ALL_ERRORS
 */
func (this *baseFutureReturnResult) State() int8 {
	return this.state
}

/**
 *  return all error from service event running
 */
func (this *baseFutureReturnResult) Errors(procName string) uerrors.CodeError {
	return this.resErrs[procName]
}

/**
 * rturn all return Results
 */
func (this *baseFutureReturnResult) ReturnResults(procName string, inReturn interface{}) uerrors.CodeError {

	// ---- validate type ---
	if reflect.TypeOf(inReturn).Kind() != reflect.Ptr {
		return servicebus.Err000003
	}

	resobj := this.resRes[procName]
	msgpack.Unmarshal(resobj, inReturn)
	return nil

}

/**
 * get the first error directly
 */
func (this *baseFutureReturnResult) Error() uerrors.CodeError {
	// --- get the first error --
	for _, errObj := range this.resErrs {
		return errObj
	}
	return nil
}

/**
 * get the first result directly
 */
func (this *baseFutureReturnResult) ReturnResult(inReturn interface{}) uerrors.CodeError {

	// ---- validate type ---
	if reflect.TypeOf(inReturn).Kind() != reflect.Ptr {
		return servicebus.Err000003
	}

	var objRef []byte
	for _, tmpObjRef := range this.resRes {
		objRef = tmpObjRef
		break
	}

	msgpack.Unmarshal(objRef, inReturn)

	return nil
}
