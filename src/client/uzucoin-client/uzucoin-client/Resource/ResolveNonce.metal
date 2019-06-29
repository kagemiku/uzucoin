//
//  ResolveNonce.metal
//  uzucoin-client
//
//  Created by Akira Fukunaga on 2019/06/24.
//  Copyright © 2019 kagemiku. All rights reserved.
//

#include <metal_stdlib>
#import  "./Loki/loki_header.metal"
using namespace metal;

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


/*********************** UTILITIES ***********************/

static size_t mystrlen(thread const char* text) {
    size_t len = 0;
    while ( text[len] != '\0' ) {
        ++len;
    }

    return len;
}

static size_t mydevstrlen(device const char* text) {
    size_t len = 0;
    while ( text[len] != '\0' ) {
        ++len;
    }

    return len;
}

static bool strcontains(thread const char* base, thread const char* findstr) {
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

static void mystrcat(thread char* dest, thread const char* src) {
    size_t destLen = mystrlen(dest);
    size_t srcLen = mystrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[destLen + i] = src[i];
    }
    dest[destLen + srcLen] = '\0';
}

static void mystrcpy(thread char* dest, thread const char* src) {
    size_t srcLen = mystrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[i] = src[i];
    }
    dest[srcLen] = '\0';
}

static void mydevstrcpy(thread char* dest, device const char* src) {
    size_t srcLen = mydevstrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[i] = src[i];
    }
    dest[srcLen] = '\0';
}

static void mystrcpydev(device char* dest, thread const char* src) {
    size_t srcLen = mystrlen(src);
    for ( int i = 0; i < srcLen; i++ ) {
        dest[i] = src[i];
    }
    dest[srcLen] = '\0';
}

static char bytetochr(BYTE buf) {
    if ( 0 <= buf && buf <= 9 ) {
        return buf + '0';
    } else {
        return buf + 'a' - 10;
    }
}

static void bytetostr(BYTE buf, thread char *result) {
    BYTE upper = (buf & 0xf0) >> 4;
    BYTE lower = buf & 0x0f;

    result[0] = bytetochr(upper);
    result[1] = bytetochr(lower);
}


/*********************** FUNCTION DEFINITIONS ***********************/
static void sha256_transform(thread SHA256_CTX *ctx, const BYTE data[])
{
    WORD a, b, c, d, e, f, g, h, i, j, t1, t2, m[64];

    for (i = 0, j = 0; i < 16; ++i, j += 4)
        m[i] = (data[j] << 24) | (data[j + 1] << 16) | (data[j + 2] << 8) | (data[j + 3]);
    for ( ; i < 64; ++i)
        m[i] = SIG1(m[i - 2]) + m[i - 7] + SIG0(m[i - 15]) + m[i - 16];

    a = ctx->state[0];
    b = ctx->state[1];
    c = ctx->state[2];
    d = ctx->state[3];
    e = ctx->state[4];
    f = ctx->state[5];
    g = ctx->state[6];
    h = ctx->state[7];

    for (i = 0; i < 64; ++i) {
        t1 = h + EP1(e) + CH(e,f,g) + k[i] + m[i];
        t2 = EP0(a) + MAJ(a,b,c);
        h = g;
        g = f;
        f = e;
        e = d + t1;
        d = c;
        c = b;
        b = a;
        a = t1 + t2;
    }

    ctx->state[0] += a;
    ctx->state[1] += b;
    ctx->state[2] += c;
    ctx->state[3] += d;
    ctx->state[4] += e;
    ctx->state[5] += f;
    ctx->state[6] += g;
    ctx->state[7] += h;
}

static void sha256_init(thread SHA256_CTX *ctx)
{
    ctx->datalen = 0;
    ctx->bitlen = 0;
    ctx->state[0] = 0x6a09e667;
    ctx->state[1] = 0xbb67ae85;
    ctx->state[2] = 0x3c6ef372;
    ctx->state[3] = 0xa54ff53a;
    ctx->state[4] = 0x510e527f;
    ctx->state[5] = 0x9b05688c;
    ctx->state[6] = 0x1f83d9ab;
    ctx->state[7] = 0x5be0cd19;
}

static void sha256_update(thread SHA256_CTX *ctx, const BYTE data[], size_t len)
{
    WORD i;

    for (i = 0; i < len; ++i) {
        ctx->data[ctx->datalen] = data[i];
        ctx->datalen++;
        if (ctx->datalen == 64) {
            sha256_transform(ctx, ctx->data);
            ctx->bitlen += 512;
            ctx->datalen = 0;
        }
    }
}

static void sha256_final(thread SHA256_CTX *ctx, BYTE hash[])
{
    WORD i;

    i = ctx->datalen;

    // Pad whatever data is left in the buffer.
    if (ctx->datalen < 56) {
        ctx->data[i++] = 0x80;
        while (i < 56)
            ctx->data[i++] = 0x00;
    }
    else {
        ctx->data[i++] = 0x80;
        while (i < 64)
            ctx->data[i++] = 0x00;
        sha256_transform(ctx, ctx->data);
        int idx = 0;
        while (idx < 56) {
            ctx->data[idx++] = 0x00;
        }
    }

    // Append to the padding the total message's length in bits and transform.
    ctx->bitlen += ctx->datalen * 8;
    ctx->data[63] = ctx->bitlen;
    ctx->data[62] = ctx->bitlen >> 8;
    ctx->data[61] = ctx->bitlen >> 16;
    ctx->data[60] = ctx->bitlen >> 24;
    ctx->data[59] = ctx->bitlen >> 32;
    ctx->data[58] = ctx->bitlen >> 40;
    ctx->data[57] = ctx->bitlen >> 48;
    ctx->data[56] = ctx->bitlen >> 56;
    sha256_transform(ctx, ctx->data);

    // Since this implementation uses little endian byte ordering and SHA uses big endian,
    // reverse all the bytes when copying the final state to the output hash.
    for (i = 0; i < 4; ++i) {
        hash[i]      = (ctx->state[0] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 4]  = (ctx->state[1] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 8]  = (ctx->state[2] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 12] = (ctx->state[3] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 16] = (ctx->state[4] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 20] = (ctx->state[5] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 24] = (ctx->state[6] >> (24 - i * 8)) & 0x000000ff;
        hash[i + 28] = (ctx->state[7] >> (24 - i * 8)) & 0x000000ff;
    }
}

static void sha256(thread const char* text, thread char* result) {
    BYTE textByte[PLAIN_TEXT_MAX_LENGTH];
    for ( int i = 0; i < mystrlen(text); i++ ) {
        textByte[i] = text[i];
    }
    textByte[mystrlen(text)] = '\0';

    BYTE buf[SHA256_BLOCK_SIZE];
    SHA256_CTX ctx;

    sha256_init(&ctx);
    sha256_update(&ctx, textByte, mystrlen(text));
    sha256_final(&ctx, buf);

    for ( int i = 0; i < SHA256_BLOCK_SIZE*2; i+=2 ) {
        bytetostr(buf[i/2], &result[i]);
    }
    result[SHA256_BLOCK_SIZE*2] = '\0';
}

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
