#ifndef _C_TICK_STORE_H_
#define _C_TICK_STORE_H_ 
#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif 

#define BLOCK_SZ_32      32
#define BLOCK_SZ_64      64
#define BLOCK_SZ_128     128
#define BLOCK_SZ_256     256
#define BLOCK_SZ_512     512
#define BLOCK_SZ_1024    1024
#define BLOCK_SZ_2048    2048
#define BLOCK_SZ_4096    4096
#define BLOCK_SZ_8192    8192
#define BLOCK_SZ_16384   16384
#define BLOCK_SZ_32768   32768
#define BLOCK_SZ_65536   65536
#define BLOCK_SZ_131072  131072
#define BLOCK_SZ_262144  262144
#define BLOCK_SZ_524288  524288


// 支持的类型
typedef char BLOCK_32[BLOCK_SZ_32];
typedef char BLOCK_64[BLOCK_SZ_64];
typedef char BLOCK_128[BLOCK_SZ_128];
typedef char BLOCK_256[BLOCK_SZ_256];
typedef char BLOCK_512[BLOCK_SZ_512];
typedef char BLOCK_1024[BLOCK_SZ_1024];
typedef char BLOCK_2048[BLOCK_SZ_2048];
typedef char BLOCK_4096[BLOCK_SZ_4096];
typedef char BLOCK_8192[BLOCK_SZ_8192];
typedef char BLOCK_16384[BLOCK_SZ_16384];
typedef char BLOCK_32768[BLOCK_SZ_32768];
typedef char BLOCK_65536[BLOCK_SZ_65536];
typedef char BLOCK_131072[BLOCK_SZ_131072];
typedef char BLOCK_262144[BLOCK_SZ_262144];
typedef char BLOCK_524288[BLOCK_SZ_524288];

/// 打开模式 (对应 C++ OpenMode)
#define TICK_STORE_OPEN            0
#define TICK_STORE_CREATE          1
#define TICK_STORE_OPEN_OR_CREATE  2

// C 语言接口
typedef struct c_tick_store {
    void* impl;  // 内部实现指针，实际类型为 TickStore<Record>*
    int32_t type;
    size_t record_size;
} c_tick_store;

typedef c_tick_store* TickStoreHandle;

extern TickStoreHandle tick_store_new(int32_t type);
extern void tick_store_free(TickStoreHandle handle);

extern int32_t tick_store_open(TickStoreHandle handle, const char* path, int32_t mode, size_t init_file_size);
extern void tick_store_close(TickStoreHandle handle);
extern int32_t tick_store_add_code(TickStoreHandle handle, const char* code, uint32_t max_records);
extern int32_t tick_store_push_back(TickStoreHandle handle, const char* code, const void* record);
extern size_t tick_store_size(TickStoreHandle handle, const char* code);
extern const void* tick_store_at(TickStoreHandle handle, const char* code, int32_t index);
extern int32_t tick_store_get_value(TickStoreHandle handle, const char* code, int32_t index, void* out_record);
extern const int32_t tick_store_set_at(TickStoreHandle handle, const char* code, int32_t index, const void* record) ;
extern void tick_store_flush(TickStoreHandle handle);
extern const char** tick_store_get_all_codes(TickStoreHandle handle, size_t* out_count);
extern int32_t tick_store_is_ring_buffer(TickStoreHandle handle, const char* code);
extern void tick_store_free_code_list(const char** codes, size_t count);


#ifdef __cplusplus
}
#endif


#endif