package main

// #cgo LDFLAGS: -lncurses
// #include <stdlib.h>
// #include <locale.h>
// #include <ncurses.h>
import "C"

import (
	"unsafe"
)

var (
	rows *int
	cols *int
)

const (
	LEGEND_WIDTH = 30
	LOG_HEIGHT   = 10
)

const (
	ACS_DEGREE = iota + 4194406
	ACS_PLMINUS
	ACS_BOARD
	ACS_LANTERN
	ACS_LRCORNER
	ACS_URCORNER
	ACS_LLCORNER
	ACS_ULCORNER
	ACS_PLUS
	ACS_S1
	ACS_S3
	ACS_HLINE
	ACS_S7
	ACS_S9
	ACS_LTEE
	ACS_RTEE
	ACS_BTEE
	ACS_TTEE
	ACS_VLINE
	ACS_LEQUAL
	ACS_GEQUAL
	ACS_PI
	ACS_NEQUAL
	ACS_STERLING
	ACS_BULLET
	ACS_LARROW    = 4194347
	ACS_RARROW    = 4194348
	ACS_DARROW    = 4194349
	ACS_UARROW    = 4194350
	ACS_BLOCK     = 4194352
	ACS_DIAMOND   = 4194400
	ACS_CKBOARD   = 4194401
	COLOR_BLACK   = C.COLOR_BLACK
	COLOR_BLUE    = C.COLOR_BLUE
	COLOR_GREEN   = C.COLOR_GREEN
	COLOR_CYAN    = C.COLOR_CYAN
	COLOR_RED     = C.COLOR_RED
	COLOR_MAGENTA = C.COLOR_MAGENTA
	COLOR_YELLOW  = C.COLOR_YELLOW
	COLOR_WHITE   = C.COLOR_WHITE
)

func ncurses_init() {
	C.initscr()
	//	C.setlocale(C.LC_ALL, C.CString(""))
	rows = (*int)(unsafe.Pointer(&C.LINES))
	cols = (*int)(unsafe.Pointer(&C.COLS))
	C.start_color()
	C.init_pair(C.short(1), C.short(C.COLOR_WHITE), C.short(C.COLOR_BLACK))
	C.init_pair(C.short(2), C.short(C.COLOR_RED), C.short(C.COLOR_BLACK))
	C.init_pair(C.short(3), C.short(C.COLOR_GREEN), C.short(C.COLOR_BLACK))
	C.init_pair(C.short(4), C.short(C.COLOR_MAGENTA), C.short(C.COLOR_BLACK))
	C.init_pair(C.short(5), C.short(C.COLOR_YELLOW), C.short(C.COLOR_BLACK))
	C.init_pair(C.short(6), C.short(C.COLOR_CYAN), C.short(C.COLOR_BLACK))
	C.curs_set(C.int(0))

	C.noecho()
}

func ncurses_draw_template() {
	ncurses_color(1)
	C.box(C.stdscr, C.chtype(ACS_VLINE), C.chtype(ACS_HLINE))
	C.mvvline(1, C.int(*cols-LEGEND_WIDTH), ACS_VLINE, C.int(*rows-(LOG_HEIGHT+1)))
	C.mvhline(C.int(*rows-LOG_HEIGHT), 1, ACS_HLINE, C.int(*cols-2))
	C.mvaddch(0, C.int(*cols-LEGEND_WIDTH), ACS_TTEE)
	C.mvaddch(C.int(*rows-LOG_HEIGHT), 0, ACS_LTEE)
	C.mvaddch(C.int(*rows-LOG_HEIGHT), C.int(*cols-1), ACS_RTEE)
	C.mvaddch(C.int(*rows-LOG_HEIGHT), C.int(*cols-LEGEND_WIDTH), ACS_BTEE)
	ncurses_move(1, *cols-((LEGEND_WIDTH/2)+6))
	ncurses_printw("LEGEND (K/D)")
	ncurses_move(*rows-LOG_HEIGHT, (*cols-(LEGEND_WIDTH+6))/2)
	ncurses_printw("GOBOTS ARENA")
}

// cgo can't do vararg, so we have to cook it
func ncurses_printw(s string) {
	for i := 0; i < len(s); i++ {
		C.addch(C.chtype(s[i]))
	}
}

func ncurses_move(y, x int) {
	C.move(C.int(y), C.int(x))
}

func ncurses_color(i int) {
	C.attron(C.COLOR_PAIR(C.int(i)))
}

func ncurses_refresh() {
	C.refresh()
}

func ncurses_endwin() {
	C.endwin()
}

func ncurses_gamefield() (y, x int) {
	return *rows - LOG_HEIGHT - 2, *cols - LEGEND_WIDTH - 2
}

func ncurses_clear() {
	C.erase()
}

func ncurses_draw_north() {
	C.addch(ACS_DARROW)
}

func ncurses_draw_south() {
	C.addch(ACS_UARROW)
}

func ncurses_draw_east() {
	C.addch(ACS_LARROW)
}

func ncurses_draw_west() {
	C.addch(ACS_RARROW)
}

func ncurses_getch() int {
	return int(C.getch())
}
