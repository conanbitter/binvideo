#include "bvf_decode.h"
#include <stdlib.h>

BVF_File* bvf_open(const char* filename) {
    BVF_File* result = malloc(sizeof(BVF_File));
    result->file = fopen(filename, "rb");

    uint16_t magic = 0;
    fread(&magic, 2, 1, result->file);
    if (magic != 0x5642) {
        printf("Wrong file format.");
        free(result);
        return NULL;
    }
    uint8_t version = 0;
    fread(&version, 1, 1, result->file);
    if (version != 1) {
        printf("Wrong file format version.");
        free(result);
        return NULL;
    }
    result->format_version = version;

    uint16_t width, height, length;
    float frame_time;
    fread(&width, 2, 1, result->file);
    fread(&height, 2, 1, result->file);
    fread(&length, 2, 1, result->file);
    fread(&frame_time, 4, 1, result->file);
    result->width = width;
    result->height = height;
    result->length = length;
    result->frame_time = 1.0f / frame_time;

    result->frames_offset = ftell(result->file);

    result->frame = malloc(result->width * result->height);
    result->buffer = NULL;
    result->buffer_size = 0;

    result->current_frame = -1;

    return result;
}

void bvf_close(BVF_File** file) {
    fclose((*file)->file);
    free((*file)->frame);
    if ((*file)->buffer != NULL) {
        free((*file)->buffer);
    }
    free(*file);
    *file = NULL;
}

uint8_t* bvf_next_frame(BVF_File* file) {
    file->current_frame++;
    if (file->current_frame >= file->length) return file->frame;

    uint32_t data_length;
    int is_clean;
    fread(&data_length, 4, 1, file->file);
    if ((data_length && 1 << 31) > 0) {
        is_clean = 1;
        data_length &= ~(1 << 31);
    }

    if (data_length > file->buffer_size) {
        free(file->buffer);
        malloc(data_length);
    }
    fread(file->buffer, data_length, 1, file->file);
    return file->frame;
}