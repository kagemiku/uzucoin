//
//  ResolveNonce.metal
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2019/06/24.
//  Copyright © 2019 kagemiku. All rights reserved.
//

#include <metal_stdlib>
#import  "./MyUtility/myutility_header.metal"
#import  "./SHA256/sha256_header.metal"
#import  "./Loki/loki_header.metal"
using namespace metal;

static void genRandomStr(thread Loki* loki, int len, thread char* result) {
    for ( int i = 0; i < len; i++ ) {
        result[i] = int(loki->rand() * 1000) % 26 + 'a';
    }
    result[len] = '\0';
}

static bool checkHash(thread const char* hash) {
    constexpr int keyNum = 6;
    thread char keys[keyNum][20] = {
        "757a756b69",
        "757a75",
        "7a756b69",
        "75",
        "7a75",
        "6b69",
    };

    for ( int i = 0; i < keyNum; i++ ) {
        if ( strcontains(hash, keys[i]) ) {
            return true;
        }
    }

    return false;
}

kernel void ResolveNonce(device const char* timestamp,
                         device const char* prevHash,
                         device char* nonce,
                         uint index [[thread_position_in_grid]])
{
    Loki loki = Loki(1, 2, 3);

    char localTimestamp[60] = "";
    char localPrevHash[70] = "";
    char localNonce[NONCE_LENGTH + 1] = "";
    mydevstrcpy(localTimestamp, timestamp);
    mydevstrcpy(localPrevHash, prevHash);

    char nonceCand[NONCE_LENGTH + 1] = "";
    while ( true ) {
        genRandomStr(&loki, NONCE_LENGTH, nonceCand);
        char text[200] = "";
        mystrcat(text, localTimestamp);
        mystrcat(text, nonceCand);
        mystrcat(text, localPrevHash);


        char hash[SHA256_BLOCK_SIZE*2 + 1] = "";
        sha256(text, hash);
        if ( checkHash(hash) ) {
            mystrcpy(localNonce, nonceCand);
            break;
        }
    }

    mystrcpydev(nonce, localNonce);
}
