/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ansi

import (
	"fmt"
)

// Color256 represents an `xterm-256color` ANSI color code.
type Color256 int

// For W3C color list, see:
// - https://jonasjacek.github.io/colors/
// - https://en.wikipedia.org/wiki/X11_color_names#Clashes_between_web_and_X11_colors_in_the_CSS_color_scheme

const (
	// Color256Black is an `xterm-256color` representing `Black` (#000000).
	Color256Black Color256 = 0
	// Color256Maroon is an `xterm-256color` representing `Maroon` (#800000).
	Color256Maroon Color256 = 1
	// Color256Green is an `xterm-256color` representing `Green` (#008000).
	Color256Green Color256 = 2
	// Color256Olive is an `xterm-256color` representing `Olive` (#808000).
	Color256Olive Color256 = 3
	// Color256Navy is an `xterm-256color` representing `Navy` (#000080).
	Color256Navy Color256 = 4
	// Color256Purple is an `xterm-256color` representing `Purple` (#800080).
	Color256Purple Color256 = 5
	// Color256Teal is an `xterm-256color` representing `Teal` (#008080).
	Color256Teal Color256 = 6
	// Color256Silver is an `xterm-256color` representing `Silver` (#c0c0c0).
	Color256Silver Color256 = 7
	// Color256Grey is an `xterm-256color` representing `Grey` (#808080).
	Color256Grey Color256 = 8
	// Color256Red is an `xterm-256color` representing `Red` (#ff0000).
	Color256Red Color256 = 9
	// Color256Lime is an `xterm-256color` representing `Lime` (#00ff00).
	Color256Lime Color256 = 10
	// Color256Yellow is an `xterm-256color` representing `Yellow` (#ffff00).
	Color256Yellow Color256 = 11
	// Color256Blue is an `xterm-256color` representing `Blue` (#0000ff).
	Color256Blue Color256 = 12
	// Color256Fuchsia is an `xterm-256color` representing `Fuchsia` (#ff00ff).
	Color256Fuchsia Color256 = 13
	// Color256Aqua is an `xterm-256color` representing `Aqua` (#00ffff).
	Color256Aqua Color256 = 14
	// Color256White is an `xterm-256color` representing `White` (#ffffff).
	Color256White Color256 = 15
	// Color256Grey0 is an `xterm-256color` representing `Grey0` (#000000).
	Color256Grey0 Color256 = 16
	// Color256NavyBlue is an `xterm-256color` representing `NavyBlue` (#00005f).
	Color256NavyBlue Color256 = 17
	// Color256DarkBlue is an `xterm-256color` representing `DarkBlue` (#000087).
	Color256DarkBlue Color256 = 18
	// Color256Blue3 is an `xterm-256color` representing `Blue3` (#0000af).
	Color256Blue3 Color256 = 19
	// Color256Blue3Alt2 is an `xterm-256color` representing `Blue3` (#0000d7).
	// The `Alt2` suffix was added because the name `Blue3` describes
	// multiple colors in the W3C color list.
	Color256Blue3Alt2 Color256 = 20
	// Color256Blue1 is an `xterm-256color` representing `Blue1` (#0000ff).
	Color256Blue1 Color256 = 21
	// Color256DarkGreen is an `xterm-256color` representing `DarkGreen` (#005f00).
	Color256DarkGreen Color256 = 22
	// Color256DeepSkyBlue4 is an `xterm-256color` representing `DeepSkyBlue4` (#005f5f).
	Color256DeepSkyBlue4 Color256 = 23
	// Color256DeepSkyBlue4Alt2 is an `xterm-256color` representing `DeepSkyBlue4` (#005f87).
	// The `Alt2` suffix was added because the name `DeepSkyBlue4` describes
	// multiple colors in the W3C color list.
	Color256DeepSkyBlue4Alt2 Color256 = 24
	// Color256DeepSkyBlue4Alt3 is an `xterm-256color` representing `DeepSkyBlue4` (#005faf).
	// The `Alt3` suffix was added because the name `DeepSkyBlue4` describes
	// multiple colors in the W3C color list.
	Color256DeepSkyBlue4Alt3 Color256 = 25
	// Color256DodgerBlue3 is an `xterm-256color` representing `DodgerBlue3` (#005fd7).
	Color256DodgerBlue3 Color256 = 26
	// Color256DodgerBlue2 is an `xterm-256color` representing `DodgerBlue2` (#005fff).
	Color256DodgerBlue2 Color256 = 27
	// Color256Green4 is an `xterm-256color` representing `Green4` (#008700).
	Color256Green4 Color256 = 28
	// Color256SpringGreen4 is an `xterm-256color` representing `SpringGreen4` (#00875f).
	Color256SpringGreen4 Color256 = 29
	// Color256Turquoise4 is an `xterm-256color` representing `Turquoise4` (#008787).
	Color256Turquoise4 Color256 = 30
	// Color256DeepSkyBlue3 is an `xterm-256color` representing `DeepSkyBlue3` (#0087af).
	Color256DeepSkyBlue3 Color256 = 31
	// Color256DeepSkyBlue3Alt2 is an `xterm-256color` representing `DeepSkyBlue3` (#0087d7).
	// The `Alt2` suffix was added because the name `DeepSkyBlue3` describes
	// multiple colors in the W3C color list.
	Color256DeepSkyBlue3Alt2 Color256 = 32
	// Color256DodgerBlue1 is an `xterm-256color` representing `DodgerBlue1` (#0087ff).
	Color256DodgerBlue1 Color256 = 33
	// Color256Green3 is an `xterm-256color` representing `Green3` (#00af00).
	Color256Green3 Color256 = 34
	// Color256SpringGreen3 is an `xterm-256color` representing `SpringGreen3` (#00af5f).
	Color256SpringGreen3 Color256 = 35
	// Color256DarkCyan is an `xterm-256color` representing `DarkCyan` (#00af87).
	Color256DarkCyan Color256 = 36
	// Color256LightSeaGreen is an `xterm-256color` representing `LightSeaGreen` (#00afaf).
	Color256LightSeaGreen Color256 = 37
	// Color256DeepSkyBlue2 is an `xterm-256color` representing `DeepSkyBlue2` (#00afd7).
	Color256DeepSkyBlue2 Color256 = 38
	// Color256DeepSkyBlue1 is an `xterm-256color` representing `DeepSkyBlue1` (#00afff).
	Color256DeepSkyBlue1 Color256 = 39
	// Color256Green3Alt2 is an `xterm-256color` representing `Green3` (#00d700).
	// The `Alt2` suffix was added because the name `Green3` describes
	// multiple colors in the W3C color list.
	Color256Green3Alt2 Color256 = 40
	// Color256SpringGreen3Alt2 is an `xterm-256color` representing `SpringGreen3` (#00d75f).
	// The `Alt2` suffix was added because the name `SpringGreen3` describes
	// multiple colors in the W3C color list.
	Color256SpringGreen3Alt2 Color256 = 41
	// Color256SpringGreen2 is an `xterm-256color` representing `SpringGreen2` (#00d787).
	Color256SpringGreen2 Color256 = 42
	// Color256Cyan3 is an `xterm-256color` representing `Cyan3` (#00d7af).
	Color256Cyan3 Color256 = 43
	// Color256DarkTurquoise is an `xterm-256color` representing `DarkTurquoise` (#00d7d7).
	Color256DarkTurquoise Color256 = 44
	// Color256Turquoise2 is an `xterm-256color` representing `Turquoise2` (#00d7ff).
	Color256Turquoise2 Color256 = 45
	// Color256Green1 is an `xterm-256color` representing `Green1` (#00ff00).
	Color256Green1 Color256 = 46
	// Color256SpringGreen2Alt2 is an `xterm-256color` representing `SpringGreen2` (#00ff5f).
	// The `Alt2` suffix was added because the name `SpringGreen2` describes
	// multiple colors in the W3C color list.
	Color256SpringGreen2Alt2 Color256 = 47
	// Color256SpringGreen1 is an `xterm-256color` representing `SpringGreen1` (#00ff87).
	Color256SpringGreen1 Color256 = 48
	// Color256MediumSpringGreen is an `xterm-256color` representing `MediumSpringGreen` (#00ffaf).
	Color256MediumSpringGreen Color256 = 49
	// Color256Cyan2 is an `xterm-256color` representing `Cyan2` (#00ffd7).
	Color256Cyan2 Color256 = 50
	// Color256Cyan1 is an `xterm-256color` representing `Cyan1` (#00ffff).
	Color256Cyan1 Color256 = 51
	// Color256DarkRed is an `xterm-256color` representing `DarkRed` (#5f0000).
	Color256DarkRed Color256 = 52
	// Color256DeepPink4 is an `xterm-256color` representing `DeepPink4` (#5f005f).
	Color256DeepPink4 Color256 = 53
	// Color256Purple4 is an `xterm-256color` representing `Purple4` (#5f0087).
	Color256Purple4 Color256 = 54
	// Color256Purple4Alt2 is an `xterm-256color` representing `Purple4` (#5f00af).
	// The `Alt2` suffix was added because the name `Purple4` describes
	// multiple colors in the W3C color list.
	Color256Purple4Alt2 Color256 = 55
	// Color256Purple3 is an `xterm-256color` representing `Purple3` (#5f00d7).
	Color256Purple3 Color256 = 56
	// Color256BlueViolet is an `xterm-256color` representing `BlueViolet` (#5f00ff).
	Color256BlueViolet Color256 = 57
	// Color256Orange4 is an `xterm-256color` representing `Orange4` (#5f5f00).
	Color256Orange4 Color256 = 58
	// Color256Grey37 is an `xterm-256color` representing `Grey37` (#5f5f5f).
	Color256Grey37 Color256 = 59
	// Color256MediumPurple4 is an `xterm-256color` representing `MediumPurple4` (#5f5f87).
	Color256MediumPurple4 Color256 = 60
	// Color256SlateBlue3 is an `xterm-256color` representing `SlateBlue3` (#5f5faf).
	Color256SlateBlue3 Color256 = 61
	// Color256SlateBlue3Alt2 is an `xterm-256color` representing `SlateBlue3` (#5f5fd7).
	// The `Alt2` suffix was added because the name `SlateBlue3` describes
	// multiple colors in the W3C color list.
	Color256SlateBlue3Alt2 Color256 = 62
	// Color256RoyalBlue1 is an `xterm-256color` representing `RoyalBlue1` (#5f5fff).
	Color256RoyalBlue1 Color256 = 63
	// Color256Chartreuse4 is an `xterm-256color` representing `Chartreuse4` (#5f8700).
	Color256Chartreuse4 Color256 = 64
	// Color256DarkSeaGreen4 is an `xterm-256color` representing `DarkSeaGreen4` (#5f875f).
	Color256DarkSeaGreen4 Color256 = 65
	// Color256PaleTurquoise4 is an `xterm-256color` representing `PaleTurquoise4` (#5f8787).
	Color256PaleTurquoise4 Color256 = 66
	// Color256SteelBlue is an `xterm-256color` representing `SteelBlue` (#5f87af).
	Color256SteelBlue Color256 = 67
	// Color256SteelBlue3 is an `xterm-256color` representing `SteelBlue3` (#5f87d7).
	Color256SteelBlue3 Color256 = 68
	// Color256CornflowerBlue is an `xterm-256color` representing `CornflowerBlue` (#5f87ff).
	Color256CornflowerBlue Color256 = 69
	// Color256Chartreuse3 is an `xterm-256color` representing `Chartreuse3` (#5faf00).
	Color256Chartreuse3 Color256 = 70
	// Color256DarkSeaGreen4Alt2 is an `xterm-256color` representing `DarkSeaGreen4` (#5faf5f).
	// The `Alt2` suffix was added because the name `DarkSeaGreen4` describes
	// multiple colors in the W3C color list.
	Color256DarkSeaGreen4Alt2 Color256 = 71
	// Color256CadetBlue is an `xterm-256color` representing `CadetBlue` (#5faf87).
	Color256CadetBlue Color256 = 72
	// Color256CadetBlueAlt2 is an `xterm-256color` representing `CadetBlue` (#5fafaf).
	// The `Alt2` suffix was added because the name `CadetBlue` describes
	// multiple colors in the W3C color list.
	Color256CadetBlueAlt2 Color256 = 73
	// Color256SkyBlue3 is an `xterm-256color` representing `SkyBlue3` (#5fafd7).
	Color256SkyBlue3 Color256 = 74
	// Color256SteelBlue1 is an `xterm-256color` representing `SteelBlue1` (#5fafff).
	Color256SteelBlue1 Color256 = 75
	// Color256Chartreuse3Alt2 is an `xterm-256color` representing `Chartreuse3` (#5fd700).
	// The `Alt2` suffix was added because the name `Chartreuse3` describes
	// multiple colors in the W3C color list.
	Color256Chartreuse3Alt2 Color256 = 76
	// Color256PaleGreen3 is an `xterm-256color` representing `PaleGreen3` (#5fd75f).
	Color256PaleGreen3 Color256 = 77
	// Color256SeaGreen3 is an `xterm-256color` representing `SeaGreen3` (#5fd787).
	Color256SeaGreen3 Color256 = 78
	// Color256Aquamarine3 is an `xterm-256color` representing `Aquamarine3` (#5fd7af).
	Color256Aquamarine3 Color256 = 79
	// Color256MediumTurquoise is an `xterm-256color` representing `MediumTurquoise` (#5fd7d7).
	Color256MediumTurquoise Color256 = 80
	// Color256SteelBlue1Alt2 is an `xterm-256color` representing `SteelBlue1` (#5fd7ff).
	// The `Alt2` suffix was added because the name `SteelBlue1` describes
	// multiple colors in the W3C color list.
	Color256SteelBlue1Alt2 Color256 = 81
	// Color256Chartreuse2 is an `xterm-256color` representing `Chartreuse2` (#5fff00).
	Color256Chartreuse2 Color256 = 82
	// Color256SeaGreen2 is an `xterm-256color` representing `SeaGreen2` (#5fff5f).
	Color256SeaGreen2 Color256 = 83
	// Color256SeaGreen1 is an `xterm-256color` representing `SeaGreen1` (#5fff87).
	Color256SeaGreen1 Color256 = 84
	// Color256SeaGreen1Alt2 is an `xterm-256color` representing `SeaGreen1` (#5fffaf).
	// The `Alt2` suffix was added because the name `SeaGreen1` describes
	// multiple colors in the W3C color list.
	Color256SeaGreen1Alt2 Color256 = 85
	// Color256Aquamarine1 is an `xterm-256color` representing `Aquamarine1` (#5fffd7).
	Color256Aquamarine1 Color256 = 86
	// Color256DarkSlateGray2 is an `xterm-256color` representing `DarkSlateGray2` (#5fffff).
	Color256DarkSlateGray2 Color256 = 87
	// Color256DarkRedAlt2 is an `xterm-256color` representing `DarkRed` (#870000).
	// The `Alt2` suffix was added because the name `DarkRed` describes
	// multiple colors in the W3C color list.
	Color256DarkRedAlt2 Color256 = 88
	// Color256DeepPink4Alt2 is an `xterm-256color` representing `DeepPink4` (#87005f).
	// The `Alt2` suffix was added because the name `DeepPink4` describes
	// multiple colors in the W3C color list.
	Color256DeepPink4Alt2 Color256 = 89
	// Color256DarkMagenta is an `xterm-256color` representing `DarkMagenta` (#870087).
	Color256DarkMagenta Color256 = 90
	// Color256DarkMagentaAlt2 is an `xterm-256color` representing `DarkMagenta` (#8700af).
	// The `Alt2` suffix was added because the name `DarkMagenta` describes
	// multiple colors in the W3C color list.
	Color256DarkMagentaAlt2 Color256 = 91
	// Color256DarkViolet is an `xterm-256color` representing `DarkViolet` (#8700d7).
	Color256DarkViolet Color256 = 92
	// Color256PurpleAlt2 is an `xterm-256color` representing `Purple` (#8700ff).
	// The `Alt2` suffix was added because the name `Purple` describes
	// multiple colors in the W3C color list.
	Color256PurpleAlt2 Color256 = 93
	// Color256Orange4Alt2 is an `xterm-256color` representing `Orange4` (#875f00).
	// The `Alt2` suffix was added because the name `Orange4` describes
	// multiple colors in the W3C color list.
	Color256Orange4Alt2 Color256 = 94
	// Color256LightPink4 is an `xterm-256color` representing `LightPink4` (#875f5f).
	Color256LightPink4 Color256 = 95
	// Color256Plum4 is an `xterm-256color` representing `Plum4` (#875f87).
	Color256Plum4 Color256 = 96
	// Color256MediumPurple3 is an `xterm-256color` representing `MediumPurple3` (#875faf).
	Color256MediumPurple3 Color256 = 97
	// Color256MediumPurple3Alt2 is an `xterm-256color` representing `MediumPurple3` (#875fd7).
	// The `Alt2` suffix was added because the name `MediumPurple3` describes
	// multiple colors in the W3C color list.
	Color256MediumPurple3Alt2 Color256 = 98
	// Color256SlateBlue1 is an `xterm-256color` representing `SlateBlue1` (#875fff).
	Color256SlateBlue1 Color256 = 99
	// Color256Yellow4 is an `xterm-256color` representing `Yellow4` (#878700).
	Color256Yellow4 Color256 = 100
	// Color256Wheat4 is an `xterm-256color` representing `Wheat4` (#87875f).
	Color256Wheat4 Color256 = 101
	// Color256Grey53 is an `xterm-256color` representing `Grey53` (#878787).
	Color256Grey53 Color256 = 102
	// Color256LightSlateGrey is an `xterm-256color` representing `LightSlateGrey` (#8787af).
	Color256LightSlateGrey Color256 = 103
	// Color256MediumPurple is an `xterm-256color` representing `MediumPurple` (#8787d7).
	Color256MediumPurple Color256 = 104
	// Color256LightSlateBlue is an `xterm-256color` representing `LightSlateBlue` (#8787ff).
	Color256LightSlateBlue Color256 = 105
	// Color256Yellow4Alt2 is an `xterm-256color` representing `Yellow4` (#87af00).
	// The `Alt2` suffix was added because the name `Yellow4` describes
	// multiple colors in the W3C color list.
	Color256Yellow4Alt2 Color256 = 106
	// Color256DarkOliveGreen3 is an `xterm-256color` representing `DarkOliveGreen3` (#87af5f).
	Color256DarkOliveGreen3 Color256 = 107
	// Color256DarkSeaGreen is an `xterm-256color` representing `DarkSeaGreen` (#87af87).
	Color256DarkSeaGreen Color256 = 108
	// Color256LightSkyBlue3 is an `xterm-256color` representing `LightSkyBlue3` (#87afaf).
	Color256LightSkyBlue3 Color256 = 109
	// Color256LightSkyBlue3Alt2 is an `xterm-256color` representing `LightSkyBlue3` (#87afd7).
	// The `Alt2` suffix was added because the name `LightSkyBlue3` describes
	// multiple colors in the W3C color list.
	Color256LightSkyBlue3Alt2 Color256 = 110
	// Color256SkyBlue2 is an `xterm-256color` representing `SkyBlue2` (#87afff).
	Color256SkyBlue2 Color256 = 111
	// Color256Chartreuse2Alt2 is an `xterm-256color` representing `Chartreuse2` (#87d700).
	// The `Alt2` suffix was added because the name `Chartreuse2` describes
	// multiple colors in the W3C color list.
	Color256Chartreuse2Alt2 Color256 = 112
	// Color256DarkOliveGreen3Alt2 is an `xterm-256color` representing `DarkOliveGreen3` (#87d75f).
	// The `Alt2` suffix was added because the name `DarkOliveGreen3` describes
	// multiple colors in the W3C color list.
	Color256DarkOliveGreen3Alt2 Color256 = 113
	// Color256PaleGreen3Alt2 is an `xterm-256color` representing `PaleGreen3` (#87d787).
	// The `Alt2` suffix was added because the name `PaleGreen3` describes
	// multiple colors in the W3C color list.
	Color256PaleGreen3Alt2 Color256 = 114
	// Color256DarkSeaGreen3 is an `xterm-256color` representing `DarkSeaGreen3` (#87d7af).
	Color256DarkSeaGreen3 Color256 = 115
	// Color256DarkSlateGray3 is an `xterm-256color` representing `DarkSlateGray3` (#87d7d7).
	Color256DarkSlateGray3 Color256 = 116
	// Color256SkyBlue1 is an `xterm-256color` representing `SkyBlue1` (#87d7ff).
	Color256SkyBlue1 Color256 = 117
	// Color256Chartreuse1 is an `xterm-256color` representing `Chartreuse1` (#87ff00).
	Color256Chartreuse1 Color256 = 118
	// Color256LightGreen is an `xterm-256color` representing `LightGreen` (#87ff5f).
	Color256LightGreen Color256 = 119
	// Color256LightGreenAlt2 is an `xterm-256color` representing `LightGreen` (#87ff87).
	// The `Alt2` suffix was added because the name `LightGreen` describes
	// multiple colors in the W3C color list.
	Color256LightGreenAlt2 Color256 = 120
	// Color256PaleGreen1 is an `xterm-256color` representing `PaleGreen1` (#87ffaf).
	Color256PaleGreen1 Color256 = 121
	// Color256Aquamarine1Alt2 is an `xterm-256color` representing `Aquamarine1` (#87ffd7).
	// The `Alt2` suffix was added because the name `Aquamarine1` describes
	// multiple colors in the W3C color list.
	Color256Aquamarine1Alt2 Color256 = 122
	// Color256DarkSlateGray1 is an `xterm-256color` representing `DarkSlateGray1` (#87ffff).
	Color256DarkSlateGray1 Color256 = 123
	// Color256Red3 is an `xterm-256color` representing `Red3` (#af0000).
	Color256Red3 Color256 = 124
	// Color256DeepPink4Alt3 is an `xterm-256color` representing `DeepPink4` (#af005f).
	// The `Alt3` suffix was added because the name `DeepPink4` describes
	// multiple colors in the W3C color list.
	Color256DeepPink4Alt3 Color256 = 125
	// Color256MediumVioletRed is an `xterm-256color` representing `MediumVioletRed` (#af0087).
	Color256MediumVioletRed Color256 = 126
	// Color256Magenta3 is an `xterm-256color` representing `Magenta3` (#af00af).
	Color256Magenta3 Color256 = 127
	// Color256DarkVioletAlt2 is an `xterm-256color` representing `DarkViolet` (#af00d7).
	// The `Alt2` suffix was added because the name `DarkViolet` describes
	// multiple colors in the W3C color list.
	Color256DarkVioletAlt2 Color256 = 128
	// Color256PurpleAlt3 is an `xterm-256color` representing `Purple` (#af00ff).
	// The `Alt3` suffix was added because the name `Purple` describes
	// multiple colors in the W3C color list.
	Color256PurpleAlt3 Color256 = 129
	// Color256DarkOrange3 is an `xterm-256color` representing `DarkOrange3` (#af5f00).
	Color256DarkOrange3 Color256 = 130
	// Color256IndianRed is an `xterm-256color` representing `IndianRed` (#af5f5f).
	Color256IndianRed Color256 = 131
	// Color256HotPink3 is an `xterm-256color` representing `HotPink3` (#af5f87).
	Color256HotPink3 Color256 = 132
	// Color256MediumOrchid3 is an `xterm-256color` representing `MediumOrchid3` (#af5faf).
	Color256MediumOrchid3 Color256 = 133
	// Color256MediumOrchid is an `xterm-256color` representing `MediumOrchid` (#af5fd7).
	Color256MediumOrchid Color256 = 134
	// Color256MediumPurple2 is an `xterm-256color` representing `MediumPurple2` (#af5fff).
	Color256MediumPurple2 Color256 = 135
	// Color256DarkGoldenrod is an `xterm-256color` representing `DarkGoldenrod` (#af8700).
	Color256DarkGoldenrod Color256 = 136
	// Color256LightSalmon3 is an `xterm-256color` representing `LightSalmon3` (#af875f).
	Color256LightSalmon3 Color256 = 137
	// Color256RosyBrown is an `xterm-256color` representing `RosyBrown` (#af8787).
	Color256RosyBrown Color256 = 138
	// Color256Grey63 is an `xterm-256color` representing `Grey63` (#af87af).
	Color256Grey63 Color256 = 139
	// Color256MediumPurple2Alt2 is an `xterm-256color` representing `MediumPurple2` (#af87d7).
	// The `Alt2` suffix was added because the name `MediumPurple2` describes
	// multiple colors in the W3C color list.
	Color256MediumPurple2Alt2 Color256 = 140
	// Color256MediumPurple1 is an `xterm-256color` representing `MediumPurple1` (#af87ff).
	Color256MediumPurple1 Color256 = 141
	// Color256Gold3 is an `xterm-256color` representing `Gold3` (#afaf00).
	Color256Gold3 Color256 = 142
	// Color256DarkKhaki is an `xterm-256color` representing `DarkKhaki` (#afaf5f).
	Color256DarkKhaki Color256 = 143
	// Color256NavajoWhite3 is an `xterm-256color` representing `NavajoWhite3` (#afaf87).
	Color256NavajoWhite3 Color256 = 144
	// Color256Grey69 is an `xterm-256color` representing `Grey69` (#afafaf).
	Color256Grey69 Color256 = 145
	// Color256LightSteelBlue3 is an `xterm-256color` representing `LightSteelBlue3` (#afafd7).
	Color256LightSteelBlue3 Color256 = 146
	// Color256LightSteelBlue is an `xterm-256color` representing `LightSteelBlue` (#afafff).
	Color256LightSteelBlue Color256 = 147
	// Color256Yellow3 is an `xterm-256color` representing `Yellow3` (#afd700).
	Color256Yellow3 Color256 = 148
	// Color256DarkOliveGreen3Alt3 is an `xterm-256color` representing `DarkOliveGreen3` (#afd75f).
	// The `Alt3` suffix was added because the name `DarkOliveGreen3` describes
	// multiple colors in the W3C color list.
	Color256DarkOliveGreen3Alt3 Color256 = 149
	// Color256DarkSeaGreen3Alt2 is an `xterm-256color` representing `DarkSeaGreen3` (#afd787).
	// The `Alt2` suffix was added because the name `DarkSeaGreen3` describes
	// multiple colors in the W3C color list.
	Color256DarkSeaGreen3Alt2 Color256 = 150
	// Color256DarkSeaGreen2 is an `xterm-256color` representing `DarkSeaGreen2` (#afd7af).
	Color256DarkSeaGreen2 Color256 = 151
	// Color256LightCyan3 is an `xterm-256color` representing `LightCyan3` (#afd7d7).
	Color256LightCyan3 Color256 = 152
	// Color256LightSkyBlue1 is an `xterm-256color` representing `LightSkyBlue1` (#afd7ff).
	Color256LightSkyBlue1 Color256 = 153
	// Color256GreenYellow is an `xterm-256color` representing `GreenYellow` (#afff00).
	Color256GreenYellow Color256 = 154
	// Color256DarkOliveGreen2 is an `xterm-256color` representing `DarkOliveGreen2` (#afff5f).
	Color256DarkOliveGreen2 Color256 = 155
	// Color256PaleGreen1Alt2 is an `xterm-256color` representing `PaleGreen1` (#afff87).
	// The `Alt2` suffix was added because the name `PaleGreen1` describes
	// multiple colors in the W3C color list.
	Color256PaleGreen1Alt2 Color256 = 156
	// Color256DarkSeaGreen2Alt2 is an `xterm-256color` representing `DarkSeaGreen2` (#afffaf).
	// The `Alt2` suffix was added because the name `DarkSeaGreen2` describes
	// multiple colors in the W3C color list.
	Color256DarkSeaGreen2Alt2 Color256 = 157
	// Color256DarkSeaGreen1 is an `xterm-256color` representing `DarkSeaGreen1` (#afffd7).
	Color256DarkSeaGreen1 Color256 = 158
	// Color256PaleTurquoise1 is an `xterm-256color` representing `PaleTurquoise1` (#afffff).
	Color256PaleTurquoise1 Color256 = 159
	// Color256Red3Alt2 is an `xterm-256color` representing `Red3` (#d70000).
	// The `Alt2` suffix was added because the name `Red3` describes
	// multiple colors in the W3C color list.
	Color256Red3Alt2 Color256 = 160
	// Color256DeepPink3 is an `xterm-256color` representing `DeepPink3` (#d7005f).
	Color256DeepPink3 Color256 = 161
	// Color256DeepPink3Alt2 is an `xterm-256color` representing `DeepPink3` (#d70087).
	// The `Alt2` suffix was added because the name `DeepPink3` describes
	// multiple colors in the W3C color list.
	Color256DeepPink3Alt2 Color256 = 162
	// Color256Magenta3Alt2 is an `xterm-256color` representing `Magenta3` (#d700af).
	// The `Alt2` suffix was added because the name `Magenta3` describes
	// multiple colors in the W3C color list.
	Color256Magenta3Alt2 Color256 = 163
	// Color256Magenta3Alt3 is an `xterm-256color` representing `Magenta3` (#d700d7).
	// The `Alt3` suffix was added because the name `Magenta3` describes
	// multiple colors in the W3C color list.
	Color256Magenta3Alt3 Color256 = 164
	// Color256Magenta2 is an `xterm-256color` representing `Magenta2` (#d700ff).
	Color256Magenta2 Color256 = 165
	// Color256DarkOrange3Alt2 is an `xterm-256color` representing `DarkOrange3` (#d75f00).
	// The `Alt2` suffix was added because the name `DarkOrange3` describes
	// multiple colors in the W3C color list.
	Color256DarkOrange3Alt2 Color256 = 166
	// Color256IndianRedAlt2 is an `xterm-256color` representing `IndianRed` (#d75f5f).
	// The `Alt2` suffix was added because the name `IndianRed` describes
	// multiple colors in the W3C color list.
	Color256IndianRedAlt2 Color256 = 167
	// Color256HotPink3Alt2 is an `xterm-256color` representing `HotPink3` (#d75f87).
	// The `Alt2` suffix was added because the name `HotPink3` describes
	// multiple colors in the W3C color list.
	Color256HotPink3Alt2 Color256 = 168
	// Color256HotPink2 is an `xterm-256color` representing `HotPink2` (#d75faf).
	Color256HotPink2 Color256 = 169
	// Color256Orchid is an `xterm-256color` representing `Orchid` (#d75fd7).
	Color256Orchid Color256 = 170
	// Color256MediumOrchid1 is an `xterm-256color` representing `MediumOrchid1` (#d75fff).
	Color256MediumOrchid1 Color256 = 171
	// Color256Orange3 is an `xterm-256color` representing `Orange3` (#d78700).
	Color256Orange3 Color256 = 172
	// Color256LightSalmon3Alt2 is an `xterm-256color` representing `LightSalmon3` (#d7875f).
	// The `Alt2` suffix was added because the name `LightSalmon3` describes
	// multiple colors in the W3C color list.
	Color256LightSalmon3Alt2 Color256 = 173
	// Color256LightPink3 is an `xterm-256color` representing `LightPink3` (#d78787).
	Color256LightPink3 Color256 = 174
	// Color256Pink3 is an `xterm-256color` representing `Pink3` (#d787af).
	Color256Pink3 Color256 = 175
	// Color256Plum3 is an `xterm-256color` representing `Plum3` (#d787d7).
	Color256Plum3 Color256 = 176
	// Color256Violet is an `xterm-256color` representing `Violet` (#d787ff).
	Color256Violet Color256 = 177
	// Color256Gold3Alt2 is an `xterm-256color` representing `Gold3` (#d7af00).
	// The `Alt2` suffix was added because the name `Gold3` describes
	// multiple colors in the W3C color list.
	Color256Gold3Alt2 Color256 = 178
	// Color256LightGoldenrod3 is an `xterm-256color` representing `LightGoldenrod3` (#d7af5f).
	Color256LightGoldenrod3 Color256 = 179
	// Color256Tan is an `xterm-256color` representing `Tan` (#d7af87).
	Color256Tan Color256 = 180
	// Color256MistyRose3 is an `xterm-256color` representing `MistyRose3` (#d7afaf).
	Color256MistyRose3 Color256 = 181
	// Color256Thistle3 is an `xterm-256color` representing `Thistle3` (#d7afd7).
	Color256Thistle3 Color256 = 182
	// Color256Plum2 is an `xterm-256color` representing `Plum2` (#d7afff).
	Color256Plum2 Color256 = 183
	// Color256Yellow3Alt2 is an `xterm-256color` representing `Yellow3` (#d7d700).
	// The `Alt2` suffix was added because the name `Yellow3` describes
	// multiple colors in the W3C color list.
	Color256Yellow3Alt2 Color256 = 184
	// Color256Khaki3 is an `xterm-256color` representing `Khaki3` (#d7d75f).
	Color256Khaki3 Color256 = 185
	// Color256LightGoldenrod2 is an `xterm-256color` representing `LightGoldenrod2` (#d7d787).
	Color256LightGoldenrod2 Color256 = 186
	// Color256LightYellow3 is an `xterm-256color` representing `LightYellow3` (#d7d7af).
	Color256LightYellow3 Color256 = 187
	// Color256Grey84 is an `xterm-256color` representing `Grey84` (#d7d7d7).
	Color256Grey84 Color256 = 188
	// Color256LightSteelBlue1 is an `xterm-256color` representing `LightSteelBlue1` (#d7d7ff).
	Color256LightSteelBlue1 Color256 = 189
	// Color256Yellow2 is an `xterm-256color` representing `Yellow2` (#d7ff00).
	Color256Yellow2 Color256 = 190
	// Color256DarkOliveGreen1 is an `xterm-256color` representing `DarkOliveGreen1` (#d7ff5f).
	Color256DarkOliveGreen1 Color256 = 191
	// Color256DarkOliveGreen1Alt2 is an `xterm-256color` representing `DarkOliveGreen1` (#d7ff87).
	// The `Alt2` suffix was added because the name `DarkOliveGreen1` describes
	// multiple colors in the W3C color list.
	Color256DarkOliveGreen1Alt2 Color256 = 192
	// Color256DarkSeaGreen1Alt2 is an `xterm-256color` representing `DarkSeaGreen1` (#d7ffaf).
	// The `Alt2` suffix was added because the name `DarkSeaGreen1` describes
	// multiple colors in the W3C color list.
	Color256DarkSeaGreen1Alt2 Color256 = 193
	// Color256Honeydew2 is an `xterm-256color` representing `Honeydew2` (#d7ffd7).
	Color256Honeydew2 Color256 = 194
	// Color256LightCyan1 is an `xterm-256color` representing `LightCyan1` (#d7ffff).
	Color256LightCyan1 Color256 = 195
	// Color256Red1 is an `xterm-256color` representing `Red1` (#ff0000).
	Color256Red1 Color256 = 196
	// Color256DeepPink2 is an `xterm-256color` representing `DeepPink2` (#ff005f).
	Color256DeepPink2 Color256 = 197
	// Color256DeepPink1 is an `xterm-256color` representing `DeepPink1` (#ff0087).
	Color256DeepPink1 Color256 = 198
	// Color256DeepPink1Alt2 is an `xterm-256color` representing `DeepPink1` (#ff00af).
	// The `Alt2` suffix was added because the name `DeepPink1` describes
	// multiple colors in the W3C color list.
	Color256DeepPink1Alt2 Color256 = 199
	// Color256Magenta2Alt2 is an `xterm-256color` representing `Magenta2` (#ff00d7).
	// The `Alt2` suffix was added because the name `Magenta2` describes
	// multiple colors in the W3C color list.
	Color256Magenta2Alt2 Color256 = 200
	// Color256Magenta1 is an `xterm-256color` representing `Magenta1` (#ff00ff).
	Color256Magenta1 Color256 = 201
	// Color256OrangeRed1 is an `xterm-256color` representing `OrangeRed1` (#ff5f00).
	Color256OrangeRed1 Color256 = 202
	// Color256IndianRed1 is an `xterm-256color` representing `IndianRed1` (#ff5f5f).
	Color256IndianRed1 Color256 = 203
	// Color256IndianRed1Alt2 is an `xterm-256color` representing `IndianRed1` (#ff5f87).
	// The `Alt2` suffix was added because the name `IndianRed1` describes
	// multiple colors in the W3C color list.
	Color256IndianRed1Alt2 Color256 = 204
	// Color256HotPink is an `xterm-256color` representing `HotPink` (#ff5faf).
	Color256HotPink Color256 = 205
	// Color256HotPinkAlt2 is an `xterm-256color` representing `HotPink` (#ff5fd7).
	// The `Alt2` suffix was added because the name `HotPink` describes
	// multiple colors in the W3C color list.
	Color256HotPinkAlt2 Color256 = 206
	// Color256MediumOrchid1Alt2 is an `xterm-256color` representing `MediumOrchid1` (#ff5fff).
	// The `Alt2` suffix was added because the name `MediumOrchid1` describes
	// multiple colors in the W3C color list.
	Color256MediumOrchid1Alt2 Color256 = 207
	// Color256DarkOrange is an `xterm-256color` representing `DarkOrange` (#ff8700).
	Color256DarkOrange Color256 = 208
	// Color256Salmon1 is an `xterm-256color` representing `Salmon1` (#ff875f).
	Color256Salmon1 Color256 = 209
	// Color256LightCoral is an `xterm-256color` representing `LightCoral` (#ff8787).
	Color256LightCoral Color256 = 210
	// Color256PaleVioletRed1 is an `xterm-256color` representing `PaleVioletRed1` (#ff87af).
	Color256PaleVioletRed1 Color256 = 211
	// Color256Orchid2 is an `xterm-256color` representing `Orchid2` (#ff87d7).
	Color256Orchid2 Color256 = 212
	// Color256Orchid1 is an `xterm-256color` representing `Orchid1` (#ff87ff).
	Color256Orchid1 Color256 = 213
	// Color256Orange1 is an `xterm-256color` representing `Orange1` (#ffaf00).
	Color256Orange1 Color256 = 214
	// Color256SandyBrown is an `xterm-256color` representing `SandyBrown` (#ffaf5f).
	Color256SandyBrown Color256 = 215
	// Color256LightSalmon1 is an `xterm-256color` representing `LightSalmon1` (#ffaf87).
	Color256LightSalmon1 Color256 = 216
	// Color256LightPink1 is an `xterm-256color` representing `LightPink1` (#ffafaf).
	Color256LightPink1 Color256 = 217
	// Color256Pink1 is an `xterm-256color` representing `Pink1` (#ffafd7).
	Color256Pink1 Color256 = 218
	// Color256Plum1 is an `xterm-256color` representing `Plum1` (#ffafff).
	Color256Plum1 Color256 = 219
	// Color256Gold1 is an `xterm-256color` representing `Gold1` (#ffd700).
	Color256Gold1 Color256 = 220
	// Color256LightGoldenrod2Alt2 is an `xterm-256color` representing `LightGoldenrod2` (#ffd75f).
	// The `Alt2` suffix was added because the name `LightGoldenrod2` describes
	// multiple colors in the W3C color list.
	Color256LightGoldenrod2Alt2 Color256 = 221
	// Color256LightGoldenrod2Alt3 is an `xterm-256color` representing `LightGoldenrod2` (#ffd787).
	// The `Alt3` suffix was added because the name `LightGoldenrod2` describes
	// multiple colors in the W3C color list.
	Color256LightGoldenrod2Alt3 Color256 = 222
	// Color256NavajoWhite1 is an `xterm-256color` representing `NavajoWhite1` (#ffd7af).
	Color256NavajoWhite1 Color256 = 223
	// Color256MistyRose1 is an `xterm-256color` representing `MistyRose1` (#ffd7d7).
	Color256MistyRose1 Color256 = 224
	// Color256Thistle1 is an `xterm-256color` representing `Thistle1` (#ffd7ff).
	Color256Thistle1 Color256 = 225
	// Color256Yellow1 is an `xterm-256color` representing `Yellow1` (#ffff00).
	Color256Yellow1 Color256 = 226
	// Color256LightGoldenrod1 is an `xterm-256color` representing `LightGoldenrod1` (#ffff5f).
	Color256LightGoldenrod1 Color256 = 227
	// Color256Khaki1 is an `xterm-256color` representing `Khaki1` (#ffff87).
	Color256Khaki1 Color256 = 228
	// Color256Wheat1 is an `xterm-256color` representing `Wheat1` (#ffffaf).
	Color256Wheat1 Color256 = 229
	// Color256Cornsilk1 is an `xterm-256color` representing `Cornsilk1` (#ffffd7).
	Color256Cornsilk1 Color256 = 230
	// Color256Grey100 is an `xterm-256color` representing `Grey100` (#ffffff).
	Color256Grey100 Color256 = 231
	// Color256Grey3 is an `xterm-256color` representing `Grey3` (#080808).
	Color256Grey3 Color256 = 232
	// Color256Grey7 is an `xterm-256color` representing `Grey7` (#121212).
	Color256Grey7 Color256 = 233
	// Color256Grey11 is an `xterm-256color` representing `Grey11` (#1c1c1c).
	Color256Grey11 Color256 = 234
	// Color256Grey15 is an `xterm-256color` representing `Grey15` (#262626).
	Color256Grey15 Color256 = 235
	// Color256Grey19 is an `xterm-256color` representing `Grey19` (#303030).
	Color256Grey19 Color256 = 236
	// Color256Grey23 is an `xterm-256color` representing `Grey23` (#3a3a3a).
	Color256Grey23 Color256 = 237
	// Color256Grey27 is an `xterm-256color` representing `Grey27` (#444444).
	Color256Grey27 Color256 = 238
	// Color256Grey30 is an `xterm-256color` representing `Grey30` (#4e4e4e).
	Color256Grey30 Color256 = 239
	// Color256Grey35 is an `xterm-256color` representing `Grey35` (#585858).
	Color256Grey35 Color256 = 240
	// Color256Grey39 is an `xterm-256color` representing `Grey39` (#626262).
	Color256Grey39 Color256 = 241
	// Color256Grey42 is an `xterm-256color` representing `Grey42` (#6c6c6c).
	Color256Grey42 Color256 = 242
	// Color256Grey46 is an `xterm-256color` representing `Grey46` (#767676).
	Color256Grey46 Color256 = 243
	// Color256Grey50 is an `xterm-256color` representing `Grey50` (#808080).
	Color256Grey50 Color256 = 244
	// Color256Grey54 is an `xterm-256color` representing `Grey54` (#8a8a8a).
	Color256Grey54 Color256 = 245
	// Color256Grey58 is an `xterm-256color` representing `Grey58` (#949494).
	Color256Grey58 Color256 = 246
	// Color256Grey62 is an `xterm-256color` representing `Grey62` (#9e9e9e).
	Color256Grey62 Color256 = 247
	// Color256Grey66 is an `xterm-256color` representing `Grey66` (#a8a8a8).
	Color256Grey66 Color256 = 248
	// Color256Grey70 is an `xterm-256color` representing `Grey70` (#b2b2b2).
	Color256Grey70 Color256 = 249
	// Color256Grey74 is an `xterm-256color` representing `Grey74` (#bcbcbc).
	Color256Grey74 Color256 = 250
	// Color256Grey78 is an `xterm-256color` representing `Grey78` (#c6c6c6).
	Color256Grey78 Color256 = 251
	// Color256Grey82 is an `xterm-256color` representing `Grey82` (#d0d0d0).
	Color256Grey82 Color256 = 252
	// Color256Grey85 is an `xterm-256color` representing `Grey85` (#dadada).
	Color256Grey85 Color256 = 253
	// Color256Grey89 is an `xterm-256color` representing `Grey89` (#e4e4e4).
	Color256Grey89 Color256 = 254
	// Color256Grey93 is an `xterm-256color` representing `Grey93` (#eeeeee).
	Color256Grey93 Color256 = 255
)

// Apply applies a color256 to a given string.
func (c Color256) Apply(text string) string {
	return fmt.Sprintf("\033[38;5;%dm%s%s", c, text, ColorReset)
}
