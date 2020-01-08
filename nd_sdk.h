#ifndef _ND_SDK_H
#define _ND_LIB_H
#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>


typedef void* ndBTHandle;
typedef void* ndIPCallOutHandle;

/*Agent*/
void nd_init();
void nd_free();

/*Method*/
void nd_method_entry(ndBTHandle bt_handle, char *method);
void nd_method_exit(ndBTHandle bt_handle, char *method);

/*BT*/
ndBTHandle nd_bt_begin(char* bt_name, char* correlation_header);
int nd_bt_end(ndBTHandle bt);
void nd_bt_store(ndBTHandle bt_handle,char* unique_bt_id);
ndBTHandle nd_bt_get(char* unique_bt_id);
int nd_bt_add_error(ndBTHandle bt_handle, int err_level, char *message, int mark_bt_as_error);

/*IP Callout*/
ndIPCallOutHandle nd_ip_db_callout_begin(ndBTHandle bt_handle, char *db_host , char *db_query);
int nd_ip_db_callout_end(ndBTHandle bt_handle, ndIPCallOutHandle ip_handle);
ndIPCallOutHandle nd_ip_http_callout_begin(ndBTHandle bt_handle, char *http_host , char *url );
int nd_ip_http_callout_end(ndBTHandle bt_handle, ndIPCallOutHandle ip_handle);

/* Logger */
#define FILE_PATH_SIZE 		1024 + 1
#define MAX_LOG_FILE_SIZE	10 * 1024 * 1024

#define LOG_HEADER 	        "#Date      TimeStamp|PID   |level|function_name:line|format\n"

typedef struct NDSDK_LOG
{
  char 			file_name[2][FILE_PATH_SIZE];
  int 			threadIdx;
  int			cur_file_size;
  int			max_file_size;
  void 			(*log)(int level, const char *func_name, int line, char *format, ...);
  FILE			*file_fp;
  void 			*sdk_lock;
}ndSDKLog;

enum {
  LL_DEBUG=1, 	// 1 
  LL_INFO, 	// 2
  LL_WARNING, 	// 3
  LL_CRITICAL, 	// 4
  LL_ERROR 	// 5
};

#define DATE() { \
 char time_stamp[20 + 1] = ""; \
 nd_sdk_log_time(time_stamp, sizeof(time_stamp));\
}

#define ND_SDK_LOG(LEVEL, ...) { \
  if(nd_sdk_log && nd_sdk_log->log) \
    nd_sdk_log->log(LEVEL, __func__, __LINE__,  __VA_ARGS__); \
}



#ifdef __cplusplus
}
#endif
#endif
