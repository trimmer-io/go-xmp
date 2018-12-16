// from http://www.exiv2.org/tags-xmp-digiKam.html

package digikam

type ColorLabel string

const (
	NoColorLabel ColorLabel = "0"
	LabelRed     ColorLabel = "1"
	LabelOrange  ColorLabel = "2"
	LabelYellow  ColorLabel = "3"
	LabelGreen   ColorLabel = "4"
	LabelBlue    ColorLabel = "5"
	LabelMagenta ColorLabel = "6"
	LabelGray    ColorLabel = "7"
	LabelBlack   ColorLabel = "8"
	LabelWhite   ColorLabel = "9"
)

type PickLabel string

const (
	NoPickLabel           PickLabel = "0"
	ItemRejected          PickLabel = "1"
	ItemPendingValidation           = "2"
	ItemAccepted                    = "3"
)
