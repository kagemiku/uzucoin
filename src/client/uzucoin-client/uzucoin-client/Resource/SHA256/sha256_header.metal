//
//  sha256_header.metal
//  uzucoin-client
//
//  Created by KAGE on 2019/07/02.
//  Copyright © 2019 kagemiku. All rights reserved.
//

#include <metal_stdlib>
using namespace metal;

#ifndef SHA256_H
#define SHA256_H

/// original code is written by Brad Conte (brad@bradconte.com)
/// repository: https://github.com/B-Con/crypto-algorithms
/// sha256.h, sha256.c is mixed in this file

constant int SHA256_BLOCK_SIZE = 32;            // SHA256 outputs a 32 byte digest
constant int NONCE_LENGTH = 10;
constant int PLAIN_TEXT_MAX_LENGTH = 200;

/**************************** DATA TYPES ****************************/
typedef unsigned char BYTE;             // 8-bit byte
typedef unsigned int  WORD;             // 32-bit word, change to "long" for 16-bit machines

typedef struct {
    BYTE data[64];
    WORD datalen;
    unsigned int bitlen;
    WORD state[8];
} SHA256_CTX;

/*********************** FUNCTION DECLARATIONS **********************/
static void sha256_init(thread SHA256_CTX *ctx);
static void sha256_update(thread SHA256_CTX *ctx, const BYTE data[], size_t len);
static void sha256_final(thread SHA256_CTX *ctx, BYTE hash[]);

/****************************** MACROS ******************************/
#define ROTLEFT(a,b) (((a) << (b)) | ((a) >> (32-(b))))
#define ROTRIGHT(a,b) (((a) >> (b)) | ((a) << (32-(b))))

#define CH(x,y,z) (((x) & (y)) ^ (~(x) & (z)))
#define MAJ(x,y,z) (((x) & (y)) ^ ((x) & (z)) ^ ((y) & (z)))
#define EP0(x) (ROTRIGHT(x,2) ^ ROTRIGHT(x,13) ^ ROTRIGHT(x,22))
#define EP1(x) (ROTRIGHT(x,6) ^ ROTRIGHT(x,11) ^ ROTRIGHT(x,25))
#define SIG0(x) (ROTRIGHT(x,7) ^ ROTRIGHT(x,18) ^ ((x) >> 3))
#define SIG1(x) (ROTRIGHT(x,17) ^ ROTRIGHT(x,19) ^ ((x) >> 10))

/**************************** VARIABLES *****************************/
constant WORD k[64] = {
    0x428a2f98,0x71374491,0xb5c0fbcf,0xe9b5dba5,0x3956c25b,0x59f111f1,0x923f82a4,0xab1c5ed5,
    0xd807aa98,0x12835b01,0x243185be,0x550c7dc3,0x72be5d74,0x80deb1fe,0x9bdc06a7,0xc19bf174,
    0xe49b69c1,0xefbe4786,0x0fc19dc6,0x240ca1cc,0x2de92c6f,0x4a7484aa,0x5cb0a9dc,0x76f988da,
    0x983e5152,0xa831c66d,0xb00327c8,0xbf597fc7,0xc6e00bf3,0xd5a79147,0x06ca6351,0x14292967,
    0x27b70a85,0x2e1b2138,0x4d2c6dfc,0x53380d13,0x650a7354,0x766a0abb,0x81c2c92e,0x92722c85,
    0xa2bfe8a1,0xa81a664b,0xc24b8b70,0xc76c51a3,0xd192e819,0xd6990624,0xf40e3585,0x106aa070,
    0x19a4c116,0x1e376c08,0x2748774c,0x34b0bcb5,0x391c0cb3,0x4ed8aa4a,0x5b9cca4f,0x682e6ff3,
    0x748f82ee,0x78a5636f,0x84c87814,0x8cc70208,0x90befffa,0xa4506ceb,0xbef9a3f7,0xc67178f2
};

/*********************** FUNCTION DEFINITIONS ***********************/
void sha256_transform(thread SHA256_CTX *ctx, const BYTE data[]);
void sha256_init(thread SHA256_CTX *ctx);
void sha256_update(thread SHA256_CTX *ctx, const BYTE data[], size_t len);
void sha256_final(thread SHA256_CTX *ctx, BYTE hash[]);
void sha256(thread const char* text, thread char* result);

#endif SHA256_H
