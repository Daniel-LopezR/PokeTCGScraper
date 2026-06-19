#include "stdio.h"
#include "string.h"
#include "stdint.h"
#include "scraper/generated/_cgo_export.h"
#include "../raylib/raylib.h"

#define RAYGUI_IMPLEMENTATION
#include "../raylib/raygui.h"

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

typedef struct ScraperWork ScraperWork;
struct ScraperWork {
  char* Browser;
  uint32_t  WantedListLen;
  char** WantedList;
};

typedef struct Card Card;
struct Card {
  char* Name;
  char* Url;
  char* Condition;
  char* Language;
  char* Description;
  int Quantity;
  float Price;
};

typedef struct Seller Seller;
struct Seller {
  char* Name;
  char* Region;
  char* Category;
  char* Url;
  Card* CardsAvailable;
};

int _main() {
  const int screenWidth = 800;
  const int screenHeight = 600;
  global_app_state.current_scene = MAIN;
  global_app_state.scene_name = "";
  InitWindow(screenWidth, screenHeight, "PokeTCGScraper");
  SetTargetFPS(60);
  
  bool showMessageBox = false;
  bool secretView = false;
  
  while (!WindowShouldClose()) {
	DrawFPS(screenWidth - 100, 20);

	switch(global_app_state.current_scene){
	case MAIN:
	  sprintf(global_app_state.scene_name, "Curent Scene is Main(%d)", global_app_state.current_scene);
	  break;
	case SCRAPER:
	  sprintf(global_app_state.scene_name, "Curent Scene is Scraper(%d)", global_app_state.current_scene);
	  break;
	case CONFIG:
	  sprintf(global_app_state.scene_name, "Curent Scene is Config(%d)", global_app_state.current_scene);
	  break;
	}
	BeginDrawing();
    ClearBackground(GetColor(GuiGetStyle(DEFAULT, BACKGROUND_COLOR)));
    DrawText(global_app_state.scene_name, 20, 20, 20, BLACK);
	
	if (GuiButton((Rectangle){ 20, 60, 120, 30 }, "Go to Main")){
	  global_app_state.current_scene = MAIN;
	  //showMessageBox = true;
	}
	if (GuiButton((Rectangle){ 20, 120, 120, 30 }, "Go to Scraper")){
	  global_app_state.current_scene = SCRAPER;
	  //showMessageBox = true;
	}
	if (GuiButton((Rectangle){ 20, 180, 120, 30 }, "Go to Config")){
	  global_app_state.current_scene = CONFIG;
	  //showMessageBox = true;
	}
	
	GuiTextBox((Rectangle){ 20, 300, 120, 30 }, "Search for an Item", 300, true);
	
	GuiPanel((Rectangle){ 20, 360, 120, 30 }, "Test");
	
    EndDrawing();
  }
  CloseWindow();
  return 0;
}
