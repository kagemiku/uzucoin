//
//  ResolveNonce.metal
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2019/06/24.
//  Copyright © 2019 kagemiku. All rights reserved.
//

#include <metal_stdlib>
using namespace metal;

kernel void ResolveNonce(device const char* timestamp,
                         device const char* prevHash,
                         device char* nonce,
                         uint index [[thread_position_in_grid]])
{
    int idx = 0, nonceIdx = 0;
    while ( timestamp[idx] != '\0' ) {
        nonce[nonceIdx] = timestamp[idx];
        ++idx;
        ++nonceIdx;
    }

    idx = 0;
    while ( prevHash[idx] != '\0' ) {
        nonce[nonceIdx] = prevHash[idx];
        ++idx;
        ++nonceIdx;
    }
    nonce[nonceIdx] = '\0';
}
