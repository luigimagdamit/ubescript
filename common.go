package main

// holds debug flags globally
var DEBUG_FLAG_GLOBAL = false
var DEBUG_TRACE_EXECUTION bool = true && DEBUG_FLAG_GLOBAL
var DEBUG_SCANNER_OUTPUT bool = true && DEBUG_FLAG_GLOBAL
var DEBUG_COMPILER_OUTPUT bool = true && DEBUG_FLAG_GLOBAL
var DEBUG_PRINT_CODE bool = true && DEBUG_FLAG_GLOBAL
var DEBUG_TABLE_CODE bool = false && DEBUG_FLAG_GLOBAL
var DEBUG_BYTE_READ bool = false && DEBUG_FLAG_GLOBAL
