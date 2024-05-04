#ifndef _BVFDECODE_H
#define _BVFDECODE_H

#include <stdint.h>
#include <stdio.h>

typedef uint8_t Block[16];

typedef struct BVF_File {
    // Header
    int format_version;
    int width;
    int height;
    int length;
    // Other data
    FILE* file;
    float frame_time;
    uint8_t* frame;
    int current_frame;
    long frames_offset;

    uint8_t* buffer;
    size_t buffer_size;
    Block* blocks;
    Block* last_blocks;
    size_t block_data_size;
    int blocks_width;
    int blocks_height;
    int* curve;
} BVF_File;

BVF_File* bvf_open(const char* filename);
void bvf_close(BVF_File** file);
uint8_t* bvf_next_frame(BVF_File* file);
// char* bvf_prev_frame(BVF_File* file);
//  char* bvf_seek(BVF_File* file, float seconds, int relative, int precise);
//  char* bvf_seek(BVF_File* file, int frames, int relative, int precise);

#endif