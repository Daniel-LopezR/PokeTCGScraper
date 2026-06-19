#include "stdio.h"
#include "pipe/pipe_win32.h"
#include "stdbool.h"
#include "stdint.h"
#include <stdint.h>
#include <string.h>
#include "../raylib/raylib.h"

#define RAYGUI_IMPLEMENTATION
#include "../raylib/raygui.h"


#define da_append(xs, x) { \
    if (xs.count >= xs.capacity) { \
	    if (!xs.capacity) xs.capacity = 1; \
	    xs.capacity *= 2; \
	    xs.items = realloc(xs.items, xs.capacity*sizeof(*xs.items)); \
	    printf("New size: %zu\n", xs.capacity);	\
    } \
    xs.items[xs.count++] = x; \
}


#define global_variable static

typedef enum {
  MAIN = 1,
  SCRAPER,
  CONFIG
} AppScene;

typedef struct GlobalAppState GlobalAppState;
struct GlobalAppState {
  AppScene current_scene;
  char* scene_name;
};

global_variable GlobalAppState global_app_state;

int main() {
    CreatePipes();
    CreateChildProcess();
     
    // Read Data Pipe
    // ReadFromPipe();
    //bool run = true;
    //do {
    //     char command[CMD_BUFSIZE];
    //     printf("Commands to execute: ");
    //     scanf_s("%s", command);
    //     if(strcmp(command, "EXIT") == 0 || strcmp(command, "exit") == 0){ 
    //         run = false;
    //   } else {
    //       WriteToCmdPipe(command);
    //       ReadFromPipe();
    //   }
    //   
    //} while(run);


    const int screenWidth = 800;
    const int screenHeight = 600;
    global_app_state.current_scene = MAIN;
    global_app_state.scene_name = "";
    InitWindow(screenWidth, screenHeight, "PokeTCGScraper");
    SetTargetFPS(60);
  
    while (!WindowShouldClose()) {
        // This is non blocking but runs every frame, which is bad
        // Should run maybe every second only when it-s scraping
        ReadFromPipe();
        DrawFPS(screenWidth - 100, 20);

	    switch(global_app_state.current_scene){
	    case MAIN:
	        sprintf(global_app_state.scene_name, "Current Scene is Main(%d)", global_app_state.current_scene);
	        break;
	    case SCRAPER:
	        sprintf(global_app_state.scene_name, "Current Scene is Scraper(%d)", global_app_state.current_scene);
	    break;
	    case CONFIG:
	        sprintf(global_app_state.scene_name, "Current Scene is Config(%d)", global_app_state.current_scene);
	        break;
	    }
	    BeginDrawing();
        ClearBackground(GetColor(GuiGetStyle(DEFAULT, BACKGROUND_COLOR)));
        DrawText(global_app_state.scene_name, 20, 20, 20, BLACK);
	
	    if (GuiButton((Rectangle){ 20, 60, 120, 30 }, "Go to Main")){
	        global_app_state.current_scene = MAIN;
	    }
	    if (GuiButton((Rectangle){ 20, 120, 120, 30 }, "Go to Scraper")){
	        global_app_state.current_scene = SCRAPER;
            char command[BUFSIZE] = "SCRAP";
            WriteToCmdPipe(command);
	    }
	    if (GuiButton((Rectangle){ 20, 180, 120, 30 }, "Go to Config")){
	        global_app_state.current_scene = CONFIG;
	    }
        EndDrawing();
    }
    CloseWindow();
    return 0;
}
