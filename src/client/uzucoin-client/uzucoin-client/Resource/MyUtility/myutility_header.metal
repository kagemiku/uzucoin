//
//  myutility_header.metal
//  uzucoin-client
//
//  Created by KAGE on 2019/07/02.
//  Copyright © 2019 kagemiku. All rights reserved.
//

#include <metal_stdlib>
#include "../SHA256/sha256_header.metal"
using namespace metal;

#ifndef MY_UTILITY
#define MY_UTILITY

size_t mystrlen(thread const char* text);
size_t mydevstrlen(device const char* text);
bool strcontains(thread const char* base, thread const char* findstr);
void mystrcat(thread char* dest, thread const char* src);
void mystrcpy(thread char* dest, thread const char* src);
void mydevstrcpy(thread char* dest, device const char* src);
void mystrcpydev(device char* dest, thread const char* src);
char bytetochr(BYTE buf);
void bytetostr(BYTE buf, thread char *result);

#endif MY_UTILITY
