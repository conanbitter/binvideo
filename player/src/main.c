#include <stdio.h>
#define SDL_MAIN_HANDLED
#include <SDL.h>
#include "bvf_decode.h"

SDL_Window *window;
SDL_Renderer *renderer;
SDL_Texture *screen;
SDL_Rect screen_rect = {0, 0, 0, 0};
BVF_File *video;

int screen_width = 100, screen_height = 100, screen_scale = 1;
Uint32 color_black;
Uint32 color_white;

void convert_frame(uint8_t *data, int width, int height) {
    Uint32 *pixels;
    int pitch;
    SDL_LockTexture(screen, NULL, &pixels, &pitch);
    pitch /= sizeof(Uint32);

    for (int y = 0; y < height; y++) {
        for (int x = 0; x < width; x++) {
            pixels[y * pitch + x] = data[y * width + x] ? color_white : color_black;
        }
    }

    SDL_UnlockTexture(screen);
}

void adjust_frame(int width, int height) {
    float videoAR = (float)video->width / (float)video->height;
    float windowAR = (float)width / (float)height;
    if (videoAR < windowAR) {
        screen_rect.h = height;
        screen_rect.y = 0;
        screen_rect.w = height * videoAR;
        screen_rect.x = (width - screen_rect.w) / 2;
    } else {
        screen_rect.w = width;
        screen_rect.x = 0;
        screen_rect.h = width / videoAR;
        screen_rect.y = (height - screen_rect.h) / 2;
    }
}

void set_scale(int scale) {
    SDL_SetWindowSize(window, video->width * scale, video->height * scale);
    SDL_SetWindowPosition(window, SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED);
    screen_rect.x = 0;
    screen_rect.y = 0;
    screen_rect.w = video->width * scale;
    screen_rect.h = video->height * scale;
}

void main(int argc, char *argv[]) {
    /*if (argc < 2) {
        printf("Use player <filename>");
        return;
    }*/

    video = bvf_open("../../data/test.bvf");  // argv[1]
    if (video == NULL) {
        return;
    }

    printf("frame size: %dx%d\nframes total:%d\nfps: %f\n",
           video->width,
           video->height,
           video->length,
           1.0f / video->frame_time);

    screen_width = video->width;
    screen_height = video->height;

    SDL_Init(SDL_INIT_VIDEO);
    window = SDL_CreateWindow("RVF", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, screen_width, screen_height, SDL_WINDOW_RESIZABLE);
    renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC);
    screen = SDL_CreateTexture(renderer, SDL_PIXELFORMAT_RGBA8888, SDL_TEXTUREACCESS_STREAMING, video->width, video->height);

    Uint32 format;
    SDL_QueryTexture(screen, &format, NULL, NULL, NULL);
    SDL_PixelFormat *pixel_format = SDL_AllocFormat(format);
    color_black = SDL_MapRGBA(pixel_format, 0, 0, 0, 255);
    color_white = SDL_MapRGBA(pixel_format, 255, 255, 255, 255);

    screen_rect.w = screen_width;
    screen_rect.h = screen_height;

    SDL_Event event;
    int working = 1;

    uint8_t *data = bvf_next_frame(video);
    if (data == NULL) {
        working = 0;
    }
    convert_frame(data, video->width, video->height);

    float elapsed = 0;
    float freq = SDL_GetPerformanceFrequency();
    Uint64 last_timer = SDL_GetPerformanceCounter();

    while (working) {
        while (SDL_PollEvent(&event)) {
            switch (event.type) {
                case SDL_QUIT:
                    working = 0;
                    break;
                case SDL_WINDOWEVENT:
                    if (event.window.event == SDL_WINDOWEVENT_RESIZED) {
                        int window_width = event.window.data1;
                        int window_height = event.window.data2;
                        adjust_frame(window_width, window_height);
                    }
                    break;
                case SDL_KEYDOWN:
                    switch (event.key.keysym.sym) {
                        case SDLK_ESCAPE:
                            working = 0;
                            break;
                        case SDLK_1:
                            set_scale(1);
                            break;
                        case SDLK_2:
                            set_scale(2);
                            break;
                        case SDLK_3:
                            set_scale(3);
                            break;
                        case SDLK_4:
                            set_scale(4);
                            break;
                        case SDLK_5:
                            set_scale(5);
                            break;
                    }
                    break;
            }
        }
        SDL_RenderCopy(renderer, screen, NULL, &screen_rect);
        SDL_RenderPresent(renderer);

        elapsed += (float)(SDL_GetPerformanceCounter() - last_timer) / freq;
        last_timer = SDL_GetPerformanceCounter();

        if (elapsed > video->frame_time) {
            uint8_t *data = bvf_next_frame(video);
            if (data == NULL) {
                working = 0;
            }
            convert_frame(data, video->width, video->height);
            elapsed -= video->frame_time;
        }

        SDL_Delay(5);
    }

    SDL_DestroyRenderer(renderer);
    SDL_DestroyWindow(window);
    SDL_Quit();
    bvf_close(&video);
}