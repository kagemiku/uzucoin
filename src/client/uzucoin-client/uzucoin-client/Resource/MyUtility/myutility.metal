//
//  myutility.metal
//  uzucoin-client
//
//  Created by KAGE on 2019/07/02.
//  Copyright © 2019 kagemiku. All rights reserved.
//

#include <metal_stdlib>
#include "myutility_header.metal"
using namespace metal;

size_t mystrlen(thread const char* text) {
    size_t len = 0;
    while ( text[len] != '\0' ) {
        ++len;
    }

    return len;
}

size_t mydevstrlen(device const char* text) {
    size_t len = 0;
    while ( text[len] != '\0' ) {
        ++len;
    }

    return len;
}

bool strcontains(thread const char* base, thread const char* findstr) {
    size_t baseLen = mystrlen(base);
    size_t findstrLen = mystrlen(findstr);

    if ( baseLen == findstrLen || findstrLen == 0) { return true; }
    if ( baseLen == 0 ) { return false; }

    for ( int i = 0; i <= baseLen - findstrLen; i++ ) {
        bool contained = true;
        for ( int j = 0; j < findstrLen; j++ ) {
            if ( base[i+j] != findstr[j] ) {
                contained = false;
                break;
            }
        }

        if ( contained ) {
            return true;
        }
    }

    return false;
}

void mystrcat(thread char* dest, thread const char* src) {
    size_t destLen = mystrlen(dest);
    size_t srcLen = mystrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[destLen + i] = src[i];
    }
    dest[destLen + srcLen] = '\0';
}

void mystrcpy(thread char* dest, thread const char* src) {
    size_t srcLen = mystrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[i] = src[i];
    }
    dest[srcLen] = '\0';
}

void mydevstrcpy(thread char* dest, device const char* src) {
    size_t srcLen = mydevstrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[i] = src[i];
    }
    dest[srcLen] = '\0';
}

void mystrcpydev(device char* dest, thread const char* src) {
    size_t srcLen = mystrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[i] = src[i];
    }
    dest[srcLen] = '\0';
}

char bytetochr(BYTE buf) {
    if ( 0 <= buf && buf <= 9 ) {
        return buf + '0';
    } else {
        return buf + 'a' - 10;
    }
}

void bytetostr(BYTE buf, thread char *result) {
    BYTE upper = (buf & 0xf0) >> 4;
    BYTE lower = buf & 0x0f;

    result[0] = bytetochr(upper);
    result[1] = bytetochr(lower);
}
