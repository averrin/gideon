package main

import (
	"encoding/hex"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// COLORS -- named colors mapping
var COLORS map[string]string
var once sync.Once

// GetColor -- get color by name or hex string
func GetColor(name string) sdl.Color {
	once.Do(initColors)
	if string(name[0]) != "#" {
		name = COLORS[name]
	}
	if name == "" {
		name = "#ffffff"
	}
	colorHex, _ := hex.DecodeString(name[1:])
	return sdl.Color{colorHex[0], colorHex[1], colorHex[2], 1}
}

func initColors() {
	COLORS = map[string]string{
		"sunset orange": "#FD5E53", "blue gray": "#6699CC", "forest green": "#6DAE81", "manatee": "#979AAA", "violet red": "#F75394", "sea green": "#9FE2BF", "mauvelous": "#EF98AA", "raw umber": "#714B23", "maroon": "#C8385A", "black": "#000000", "mountain meadow": "#30BA8F", "macaroni and cheese": "#FFBD88", "raw sienna": "#D68A59",
		"wild watermelon": "#FC6C85", "burnt orange": "#FF7F49", "teal blue": "#18A7B5", "antique brass": "#CD9575", "mango tango": "#FF8243", "olive green": "#BAB86C", "fern": "#71BC78", "outrageous orange": "#FF6E4A", "purple pizzazz": "#FE4EDA", "red orange": "#FF5349", "beaver": "#9F8170", "caribbean green": "#1CD3A2", "wild strawberry": "#FF43A4", "jazzberry jam": "#CA3767", "cerise": "#DD4492", "yellow orange": "#FFAE42", "blizzard blue": "#ACE5EE", "inchworm": "#B2EC5D", "laser lemon": "#FEFE22", "desert sand": "#EFCDB8", "electric lime": "#CEFF1D", "green": "#1CAC78", "fuchsia": "#C364C5", "cadet blue": "#B0B7C6", "bittersweet": "#FD7C6E", "blue violet": "#7366BD", "asparagus": "#87A96B", "tickle me pink": "#FC89AC", "goldenrod": "#FCD975", "magenta": "#F664AF", "periwinkle": "#C5D0E6", "thistle": "#EBC7DF", "outer space": "#414A4C", "orange": "#FF7538", "cotton candy": "#FFBCD9", "melon": "#FDBCB4", "aquamarine": "#78DBE2", "robin's egg blue": "#1FCECB", "sky blue": "#80DAEB", "orchid": "#E6A8D7",
		"denim": "#2B6CC4", "unmellow yellow": "#FFFF66", "banana mania": "#FAE7B5", "turquoise blue": "#77DDE7", "magic mint": "#AAF0D1", "blush": "#DE5D83", "tan": "#FAA76C", "dandelion": "#FDDB6D", "lavender": "#FCB4D5", "eggplant": "#6E5160", "blue bell": "#A2A2D0", "green blue": "#1164B4", "silver": "#CDC5C2", "green yellow": "#F0E891", "apricot": "#FDD9B5", "tropical rain forest": "#17806D", "screamin green": "#76FF7A", "violet blue": "#324AB2", "midnight blue": "#1A4876", "plum": "#8E4585", "royal purple": "#7851A9", "cornflower": "#9ACEEB", "red": "#EE204D", "neon carrot": "#FFA343", "shadow": "#8A795D", "wild blue yonder": "#A2ADD0", "almond": "#EFDECD", "pink flamingo": "#FC74FD", "razzle dazzle rose": "#FF48D0", "brick red": "#CB4154", "white": "#FFFFFF", "canary": "#FFFF99", "blue": "#1F75FE", "carnation pink": "#FFAACC", "cerulean": "#1DACD6", "vivid violet": "#8F509D", "pink sherbet": "#F78FA7", "copper": "#DD9475", "peach": "#FFCFAB", "tumbleweed": "#DEAA88", "timberwolf": "#DBD7D2", "navy blue": "#1974D2",
		"razzmatazz": "#E3256B", "granny smith apple": "#A8E4A0", "red violet": "#C0448F", "purple heart": "#7442C8", "pine green": "#158078", "wisteria": "#CDA4DE", "gray": "#95918C", "hot magenta": "#FF1DCE", "piggy pink": "#FDDDE6", "shamrock": "#45CEA2", "shocking pink": "#FB7EFD", "radical red": "#FF496C", "violet (purple)": "#926EAE", "maize": "#EDD19C", "mahogany": "#CD4A4C", "chestnut": "#BC5D58", "vivid tangerine": "#FFA089", "yellow": "#FCE883", "orange red": "#FF2B2B", "sepia": "#A5694F", "jungle green": "#3BB08F", "blue green": "#0D98BA", "gold": "#E7C697", "fuzzy wuzzy": "#CC6666", "pacific blue": "#1CA9C9", "atomic tangerine": "#FFA474", "scarlet": "#FC2847", "brown": "#B4674D", "yellow green": "#C5E384", "spring green": "#ECEABE", "mulberry": "#C54B8C", "purple mountai's majesty": "#9D81BA", "salmon": "#FF9BAA", "indigo": "#5D76CB", "sunglow": "#FFCF48", "burnt sienna": "#EA7E5D", "orange yellow": "#F8D568", "lemon yellow": "#FFF44F"}
}
