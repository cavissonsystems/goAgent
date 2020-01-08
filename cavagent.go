package nd

/*
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#cgo LDFLAGS: -L/root/go-projects/src/go-agent/. -lAprCommonCode -lm
#include "nd_sdk.h"

uintptr_t bt_handle_to_int(ndBTHandle bt_handle) {
    return (uintptr_t)bt_handle;
}
ndBTHandle bt_int_to_handle(uintptr_t bt) {
	return (ndBTHandle)bt;
}
uintptr_t ip_handle_to_int(ndIPCallOutHandle ip_handle) {
    return (uintptr_t)ip_handle;
}
ndIPCallOutHandle ip_int_to_handle(uintptr_t bt) {
	return (ndIPCallOutHandle)bt;
}
*/
import "C"

import (
	"unsafe"
        "net/http"
        "context"
)

func Method_entry(bt uint64, method string) {
	method_c := C.CString(method)
	defer C.free(unsafe.Pointer(method_c))
	C.nd_method_entry(C.bt_int_to_handle(C.uintptr_t(bt)), method_c)
}

func Method_exit(bt uint64, method string) {
	method_c := C.CString(method)
	defer C.free(unsafe.Pointer(method_c))
	C.nd_method_exit(C.bt_int_to_handle(C.uintptr_t(bt)), method_c)
}

func Sdk_init() {
	C.nd_init()
}

func Sdk_free() {
	C.nd_free()
}

func Updated_context(ctx context.Context,bt uint64)(context.Context){
	return context.WithValue(ctx,"CavissonTx",bt)
}

func RequestWithContext(ctx context.Context,req *http.Request)(*http.Request){
	reqCopy := req.WithContext(ctx)
	return reqCopy
}

func Start_transacation(name string, req *http.Request)(*http.Request){
	bt := BT_begin(name,"")
	ctx :=req.Context()
	new_ctx :=Updated_context(ctx,bt)
	req =RequestWithContext(new_ctx, req)
	return req
}

func Current_Transaction(ctx context.Context) (uint64) {
	return ctx.Value("CavissonTx").(uint64)
}

func BT_begin(bt_name string, correlation_header string) uint64 {
        bt_name_c := C.CString(bt_name)
        correlation_header_c := C.CString(correlation_header)
	defer C.free(unsafe.Pointer(bt_name_c))
	defer C.free(unsafe.Pointer(correlation_header_c))
	bt := C.nd_bt_begin(bt_name_c, correlation_header_c)
        return uint64(C.bt_handle_to_int(bt))
}

func BT_end(bt uint64) int {
	rc := C.nd_bt_end(C.bt_int_to_handle(C.uintptr_t(bt)))
	return int(rc)
}

func BT_store(bt uint64, unique_bt_id string) {
	unique_bt_id_c := C.CString(unique_bt_id)
	defer C.free(unsafe.Pointer(unique_bt_id_c))
	C.nd_bt_store(C.bt_int_to_handle(C.uintptr_t(bt)), unique_bt_id_c)
}

func BT_get(unique_bt_id string) uint64 {
	unique_bt_id_c := C.CString(unique_bt_id)
	defer C.free(unsafe.Pointer(unique_bt_id_c))

	bt := C.nd_bt_get(unique_bt_id_c)

	return uint64(C.bt_handle_to_int(bt))
}

func BT_add_error(bt uint64, err_level int, message string, mark_bt_as_error int) int {
	rc := C.nd_bt_add_error(C.bt_int_to_handle(C.uintptr_t(bt)), C.int(err_level), C.CString(message), C.int(mark_bt_as_error))
	return int(rc)
}

func IP_db_callout_begin(bt uint64, db_host string, db_query string) uint64 {
	db_host_c := C.CString(db_host)
	db_query_c := C.CString(db_query)

	defer C.free(unsafe.Pointer(db_host_c))
	defer C.free(unsafe.Pointer(db_query_c))

	ip_handle := C.nd_ip_db_callout_begin(C.bt_int_to_handle(C.uintptr_t(bt)), db_host_c, db_query_c)
	return uint64(C.ip_handle_to_int(ip_handle))
}

func IP_db_callout_end(bt uint64, ip_handle uint64) int {
	rc := C.nd_ip_db_callout_end(C.bt_int_to_handle(C.uintptr_t(bt)), C.ip_int_to_handle(C.uintptr_t(ip_handle)))
	return int(rc)
}

func IP_http_callout_begin(bt uint64, http_host string, url string) uint64 {
	http_host_c := C.CString(http_host)
	url_c := C.CString(url)

	defer C.free(unsafe.Pointer(http_host_c))
	defer C.free(unsafe.Pointer(url_c))

	http_handle := C.nd_ip_http_callout_begin(C.bt_int_to_handle(C.uintptr_t(bt)), http_host_c, url_c)
	return uint64(C.ip_handle_to_int(http_handle))
}

func IP_http_callout_end(bt uint64, ip_handle uint64) int {
	rc := C.nd_ip_http_callout_end(C.bt_int_to_handle(C.uintptr_t(bt)), C.ip_int_to_handle(C.uintptr_t(ip_handle)))
	return int(rc)
}
