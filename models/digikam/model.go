package digikam

import (
	"fmt"
	"trimmer.io/go-xmp/xmp"
)

var (
	NsDk = xmp.NewNamespace("digiKam", "http://www.digikam.org/ns/1.0/", NewModel)
)

func init() {
	xmp.Register(NsDk, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &Digikam{}
}

func MakeModel(d *xmp.Document) (*Digikam, error) {
	m, err := d.MakeModel(NsDk)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*Digikam)
	return x, nil
}

func FindModel(d *xmp.Document) *Digikam {
	if m := d.FindModel(NsDk); m != nil {
		return m.(*Digikam)
	}
	return nil
}

type Digikam struct {
	TagsList               xmp.StringList `xmp:"digiKam:TagsList"`
	CaptionsAutorNames     xmp.AltString  `xmp:"digiKam:CaptionsAutorNames"`
	CaptionsDateTimeStamps xmp.AltString  `xmp:"digiKam:CaptionsDateTimeStamps"`
	ImageHistory           string         `xmp:"digiKam:ImageHistory"`
	LensCorrectionSettings string         `xmp:"digiKam:LensCorrectionSettings"`
	ColorLabel             string         `xmp:"digiKam:ColorLabel"`
	PickLabel              string         `xmp:"digiKam:PickLabel"`
}

func (x Digikam) Can(nsName string) bool {
	return NsDk.GetName() == nsName
}

func (x Digikam) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsDk}
}

func (x *Digikam) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *Digikam) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x Digikam) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *Digikam) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *Digikam) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsDk.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *Digikam) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsDk.GetName(), err)
	}
	return nil
}
