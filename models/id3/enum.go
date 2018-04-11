// Copyright (c) 2017-2018 Alexander Eichhorn
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package id3

type MarkerType byte

const (
	MarkerTypePadding        MarkerType = 0x00 // padding (has no meaning)
	MarkerTypeEOS            MarkerType = 0x01 // end of initial silence
	MarkerTypeIntroStart     MarkerType = 0x02 // intro start
	MarkerTypeMainStart      MarkerType = 0x03 // main part start
	MarkerTypeOutroStart     MarkerType = 0x04 // outro start
	MarkerTypeOutroEnd       MarkerType = 0x05 // outro end
	MarkerTypeVerseStart     MarkerType = 0x06 // verse start
	MarkerTypeRefrainStart   MarkerType = 0x07 // refrain start
	MarkerTypeInterludeStart MarkerType = 0x08 // interlude start
	MarkerTypeThemeStart     MarkerType = 0x09 // theme start
	MarkerTypeVariationStart MarkerType = 0x0A // variation start
	MarkerTypeKeyChange      MarkerType = 0x0B // key change
	MarkerTypeTimeChange     MarkerType = 0x0C // time change
	MarkerTypeTemporaryNoise MarkerType = 0x0D // momentary unwanted noise (Snap, Crackle & Pop)
	MarkerTypeNoiseStart     MarkerType = 0x0E // sustained noise
	MarkerTypeNoiseEnd       MarkerType = 0x0F // sustained noise end
	MarkerTypeIntroEnd       MarkerType = 0x10 // intro end
	MarkerTypeMainEnd        MarkerType = 0x11 // main part end
	MarkerTypeVerseEnd       MarkerType = 0x12 // verse end
	MarkerTypeRefrainEnd     MarkerType = 0x13 // refrain end
	MarkerTypeThemeEnde      MarkerType = 0x14 // theme end
	MarkerTypeProfanityStart MarkerType = 0x15 // profanity
	MarkerTypeProfanityEnd   MarkerType = 0x16 // profanity end
	MarkerTypeAudioEnd       MarkerType = 0xFD // audio end (start of silence)
	MarkerTypeFileEnd        MarkerType = 0xFE // audio file ends
)

type UnitType byte

const (
	UnitTypeMs    UnitType = 0x1 // millisec
	UnitTypeFrame UnitType = 0x2 // MPEG frames
)

type PositionType byte

const (
	PositionTypeFrame PositionType = 0x1 // MPEG frames
	PositionTypeMs    PositionType = 0x2 // millisec
)

const (
	GenreBlues                 = 0
	GenreClassicRock           = 1
	GenreCountry               = 2
	GenreDance                 = 3
	GenreDisco                 = 4
	GenreFunk                  = 5
	GenreGrunge                = 6
	GenreHipHop                = 7
	GenreJazz                  = 8
	GenreMetal                 = 9
	GenreNewAge                = 10
	GenreOldies                = 11
	GenreOther                 = 12
	GenrePop                   = 13
	GenreRaB                   = 14
	GenreRap                   = 15
	GenreReggae                = 16
	GenreRock                  = 17
	GenreTechno                = 18
	GenreIndustrial            = 19
	GenreAlternative           = 20
	GenreSka                   = 21
	GenreDeathMetal            = 22
	GenrePranks                = 23
	GenreSoundtrack            = 24
	GenreEuroTechno            = 25
	GenreAmbient               = 26
	GenreTripHop               = 27
	GenreVocal                 = 28
	GenreJazzFunk              = 29
	GenreFusion                = 30
	GenreTrance                = 31
	GenreClassical             = 32
	GenreInstrumental          = 33
	GenreAcid                  = 34
	GenreHouse                 = 35
	GenreGame                  = 36
	GenreSoundClip             = 37
	GenreGospel                = 38
	GenreNoise                 = 39
	GenreAltRock               = 40
	GenreBass                  = 41
	GenreSoul                  = 42
	GenrePunk                  = 43
	GenreSpace                 = 44
	GenreMeditative            = 45
	GenreInstrumentalPop       = 46
	GenreInstrumentalRock      = 47
	GenreEthnic                = 48
	GenreGothic                = 49
	GenreDarkwave              = 50
	GenreTechnoIndustrial      = 51
	GenreElectronic            = 52
	GenrePopFolk               = 53
	GenreEurodance             = 54
	GenreDream                 = 55
	GenreSouthernRock          = 56
	GenreComedy                = 57
	GenreCult                  = 58
	GenreGangstaRap            = 59
	GenreTop40                 = 60
	GenreChristianRap          = 61
	GenrePopFunk               = 62
	GenreJungle                = 63
	GenreNativeAmerican        = 64
	GenreCabaret               = 65
	GenreNewWave               = 66
	GenrePsychedelic           = 67 /* sic, the misspelling is used in the specification */
	GenreRave                  = 68
	GenreShowtunes             = 69
	GenreTrailer               = 70
	GenreLoFi                  = 71
	GenreTribal                = 72
	GenreAcidPunk              = 73
	GenreAcidJazz              = 74
	GenrePolka                 = 75
	GenreRetro                 = 76
	GenreMusical               = 77
	GenreRockNRoll             = 78
	GenreHardRock              = 79
	GenreFolk                  = 80
	GenreFolkRock              = 81
	GenreNationalFolk          = 82
	GenreSwing                 = 83
	GenreFastFusion            = 84
	GenreBebob                 = 85
	GenreLatin                 = 86
	GenreRevival               = 87
	GenreCeltic                = 88
	GenreBluegrass             = 89
	GenreAvantgarde            = 90
	GenreGothicRock            = 91
	GenreProgressiveRock       = 92
	GenrePsychedelicRock       = 93
	GenreSymphonicRock         = 94
	GenreSlowRock              = 95
	GenreBigBand               = 96
	GenreChorus                = 97
	GenreEasyListening         = 98
	GenreAcoustic              = 99
	GenreHumour                = 100
	GenreSpeech                = 101
	GenreChanson               = 102
	GenreOpera                 = 103
	GenreChamberMusic          = 104
	GenreSonata                = 105
	GenreSymphony              = 106
	GenreBootyBass             = 107
	GenrePrimus                = 108
	GenrePornGroove            = 109
	GenreSatire                = 110
	GenreSlowJam               = 111
	GenreClub                  = 112
	GenreTango                 = 113
	GenreSamba                 = 114
	GenreFolklore              = 115
	GenreBallad                = 116
	GenrePowerBallad           = 117
	GenreRhythmicSoul          = 118
	GenreFreestyle             = 119
	GenreDuet                  = 120
	GenrePunkRock              = 121
	GenreDrumSolo              = 122
	GenreAcapella              = 123
	GenreEuroHouse             = 124
	GenreDanceHall             = 125
	GenreGoa                   = 126
	GenreDrumNBass             = 127
	GenreClubHouse             = 128
	GenreHardcore              = 129
	GenreTerror                = 130
	GenreIndie                 = 131
	GenreBritPop               = 132
	GenreAfroPunk              = 133
	GenrePolskPunk             = 134
	GenreBeat                  = 135
	GenreChristianGangsta      = 136
	GenreHeavyMetal            = 137
	GenreBlackMetal            = 138
	GenreCrossover             = 139
	GenreContemporaryChristian = 140
	GenreChristianRock         = 141
	GenreMerengue              = 142
	GenreSalsa                 = 143
	GenreThrashMetal           = 144
	GenreAnime                 = 145
	GenreJPop                  = 146
	GenreSynthPop              = 147
	// ref http://alicja.homelinux.com/~mats/text/Music/MP3/ID3/Genres.txt
	GenreAbstract          = 148
	GenreArtRock           = 149
	GenreBaroque           = 150
	GenreBhangra           = 151
	GenreBigBeat           = 152
	GenreBreakbeat         = 153
	GenreChillout          = 154
	GenreDowntempo         = 155
	GenreDub               = 156
	GenreEBM               = 157
	GenreEclectic          = 158
	GenreElectro           = 159
	GenreElectroclash      = 160
	GenreEmo               = 161
	GenreExperimental      = 162
	GenreGarage            = 163
	GenreGlobal            = 164
	GenreIDM               = 165
	GenreIllbient          = 166
	GenreIndustroGoth      = 167
	GenreJamBand           = 168
	GenreKrautrock         = 169
	GenreLeftfield         = 170
	GenreLounge            = 171
	GenreMathRock          = 172
	GenreNewRomantic       = 173
	GenreNuBreakz          = 174
	GenrePostPunk          = 175
	GenrePostRock          = 176
	GenrePsytrance         = 177
	GenreShoegaze          = 178
	GenreSpaceRock         = 179
	GenreTropRock          = 180
	GenreWorldMusic        = 181
	GenreNeoclassical      = 182
	GenreAudiobook         = 183
	GenreAudioTheatre      = 184
	GenreNeueDeutscheWelle = 185
	GenrePodcast           = 186
	GenreIndieRock         = 187
	GenreGFunk             = 188
	GenreDubstep           = 189
	GenreGarageRock        = 190
	GenrePsybient          = 191
	GenreNone              = 255
)

var GenreMap map[GenreV1]string = map[GenreV1]string{
	GenreBlues:                 "Blues",
	GenreClassicRock:           "Classic Rock",
	GenreCountry:               "Country",
	GenreDance:                 "Dance",
	GenreDisco:                 "Disco",
	GenreFunk:                  "Funk",
	GenreGrunge:                "Grunge",
	GenreHipHop:                "Hip-Hop",
	GenreJazz:                  "Jazz",
	GenreMetal:                 "Metal",
	GenreNewAge:                "New Age",
	GenreOldies:                "Oldies",
	GenreOther:                 "Other",
	GenrePop:                   "Pop",
	GenreRaB:                   "R&B",
	GenreRap:                   "Rap",
	GenreReggae:                "Reggae",
	GenreRock:                  "Rock",
	GenreTechno:                "Techno",
	GenreIndustrial:            "Industrial",
	GenreAlternative:           "Alternative",
	GenreSka:                   "Ska",
	GenreDeathMetal:            "Death Metal",
	GenrePranks:                "Pranks",
	GenreSoundtrack:            "Soundtrack",
	GenreEuroTechno:            "Euro-Techno",
	GenreAmbient:               "Ambient",
	GenreTripHop:               "Trip-Hop",
	GenreVocal:                 "Vocal",
	GenreJazzFunk:              "Jazz+Funk",
	GenreFusion:                "Fusion",
	GenreTrance:                "Trance",
	GenreClassical:             "Classical",
	GenreInstrumental:          "Instrumental",
	GenreAcid:                  "Acid",
	GenreHouse:                 "House",
	GenreGame:                  "Game",
	GenreSoundClip:             "Sound Clip",
	GenreGospel:                "Gospel",
	GenreNoise:                 "Noise",
	GenreAltRock:               "Alt. Rock",
	GenreBass:                  "Bass",
	GenreSoul:                  "Soul",
	GenrePunk:                  "Punk",
	GenreSpace:                 "Space",
	GenreMeditative:            "Meditative",
	GenreInstrumentalPop:       "Instrumental Pop",
	GenreInstrumentalRock:      "Instrumental Rock",
	GenreEthnic:                "Ethnic",
	GenreGothic:                "Gothic",
	GenreDarkwave:              "Darkwave",
	GenreTechnoIndustrial:      "Techno-Industrial",
	GenreElectronic:            "Electronic",
	GenrePopFolk:               "Pop-Folk",
	GenreEurodance:             "Eurodance",
	GenreDream:                 "Dream",
	GenreSouthernRock:          "Southern Rock",
	GenreComedy:                "Comedy",
	GenreCult:                  "Cult",
	GenreGangstaRap:            "Gangsta Rap",
	GenreTop40:                 "Top 40",
	GenreChristianRap:          "Christian Rap",
	GenrePopFunk:               "Pop/Funk",
	GenreJungle:                "Jungle",
	GenreNativeAmerican:        "Native American",
	GenreCabaret:               "Cabaret",
	GenreNewWave:               "New Wave",
	GenrePsychedelic:           "Psychedelic",
	GenreRave:                  "Rave",
	GenreShowtunes:             "Showtunes",
	GenreTrailer:               "Trailer",
	GenreLoFi:                  "Lo-Fi",
	GenreTribal:                "Tribal",
	GenreAcidPunk:              "Acid Punk",
	GenreAcidJazz:              "Acid Jazz",
	GenrePolka:                 "Polka",
	GenreRetro:                 "Retro",
	GenreMusical:               "Musical",
	GenreRockNRoll:             "Rock & Roll",
	GenreHardRock:              "Hard Rock",
	GenreFolk:                  "Folk",
	GenreFolkRock:              "Folk-Rock",
	GenreNationalFolk:          "National Folk",
	GenreSwing:                 "Swing",
	GenreFastFusion:            "Fast-Fusion",
	GenreBebob:                 "Bebob",
	GenreLatin:                 "Latin",
	GenreRevival:               "Revival",
	GenreCeltic:                "Celtic",
	GenreBluegrass:             "Bluegrass",
	GenreAvantgarde:            "Avantgarde",
	GenreGothicRock:            "Gothic Rock",
	GenreProgressiveRock:       "Progressive Rock",
	GenrePsychedelicRock:       "Psychedelic Rock",
	GenreSymphonicRock:         "Symphonic Rock",
	GenreSlowRock:              "Slow Rock",
	GenreBigBand:               "Big Band",
	GenreChorus:                "Chorus",
	GenreEasyListening:         "Easy Listening",
	GenreAcoustic:              "Acoustic",
	GenreHumour:                "Humour",
	GenreSpeech:                "Speech",
	GenreChanson:               "Chanson",
	GenreOpera:                 "Opera",
	GenreChamberMusic:          "Chamber Music",
	GenreSonata:                "Sonata",
	GenreSymphony:              "Symphony",
	GenreBootyBass:             "Booty Bass",
	GenrePrimus:                "Primus",
	GenrePornGroove:            "Porn Groove",
	GenreSatire:                "Satire",
	GenreSlowJam:               "Slow Jam",
	GenreClub:                  "Club",
	GenreTango:                 "Tango",
	GenreSamba:                 "Samba",
	GenreFolklore:              "Folklore",
	GenreBallad:                "Ballad",
	GenrePowerBallad:           "Power Ballad",
	GenreRhythmicSoul:          "Rhythmic Soul",
	GenreFreestyle:             "Freestyle",
	GenreDuet:                  "Duet",
	GenrePunkRock:              "Punk Rock",
	GenreDrumSolo:              "Drum Solo",
	GenreAcapella:              "A Capella",
	GenreEuroHouse:             "Euro-House",
	GenreDanceHall:             "Dance Hall",
	GenreGoa:                   "Goa",
	GenreDrumNBass:             "Drum & Bass",
	GenreClubHouse:             "Club-House",
	GenreHardcore:              "Hardcore",
	GenreTerror:                "Terror",
	GenreIndie:                 "Indie",
	GenreBritPop:               "BritPop",
	GenreAfroPunk:              "Afro Punk",
	GenrePolskPunk:             "Polsk Punk",
	GenreBeat:                  "Beat",
	GenreChristianGangsta:      "Christian Gangsta",
	GenreHeavyMetal:            "Heavy Metal",
	GenreBlackMetal:            "Black Metal",
	GenreCrossover:             "Crossover",
	GenreContemporaryChristian: "Contemporary Christian",
	GenreChristianRock:         "Christian Rock",
	GenreMerengue:              "Merengue",
	GenreSalsa:                 "Salsa",
	GenreThrashMetal:           "Thrash Metal",
	GenreAnime:                 "Anime",
	GenreJPop:                  "JPop",
	GenreSynthPop:              "SynthPop",
	GenreAbstract:              "Abstract",
	GenreArtRock:               "Art Rock",
	GenreBaroque:               "Baroque",
	GenreBhangra:               "Bhangra",
	GenreBigBeat:               "Big Beat",
	GenreBreakbeat:             "Breakbeat",
	GenreChillout:              "Chillout",
	GenreDowntempo:             "Downtempo",
	GenreDub:                   "Dub",
	GenreEBM:                   "EBM",
	GenreEclectic:              "Eclectic",
	GenreElectro:               "Electro",
	GenreElectroclash:          "Electroclash",
	GenreEmo:                   "Emo",
	GenreExperimental:          "Experimental",
	GenreGarage:                "Garage",
	GenreGlobal:                "Global",
	GenreIDM:                   "IDM",
	GenreIllbient:              "Illbient",
	GenreIndustroGoth:          "Industro-Goth",
	GenreJamBand:               "Jam Band",
	GenreKrautrock:             "Krautrock",
	GenreLeftfield:             "Leftfield",
	GenreLounge:                "Lounge",
	GenreMathRock:              "Math Rock",
	GenreNewRomantic:           "New Romantic",
	GenreNuBreakz:              "Nu-Breakz",
	GenrePostPunk:              "Post-Punk",
	GenrePostRock:              "Post-Rock",
	GenrePsytrance:             "Psytrance",
	GenreShoegaze:              "Shoegaze",
	GenreSpaceRock:             "Space Rock",
	GenreTropRock:              "Trop Rock",
	GenreWorldMusic:            "World Music",
	GenreNeoclassical:          "Neoclassical",
	GenreAudiobook:             "Audiobook",
	GenreAudioTheatre:          "Audio Theatre",
	GenreNeueDeutscheWelle:     "Neue Deutsche Welle",
	GenrePodcast:               "Podcast",
	GenreIndieRock:             "Indie Rock",
	GenreGFunk:                 "G-Funk",
	GenreDubstep:               "Dubstep",
	GenreGarageRock:            "Garage Rock",
	GenrePsybient:              "Psybient",
	GenreNone:                  "None",
}

type LyricsType byte

const (
	LyricsTypeOther      LyricsType = 0x00 // other
	LyricsTypeLyrics     LyricsType = 0x01 // lyrics
	LyricsTypeTranscript LyricsType = 0x02 // text transcription
	LyricsTypePart       LyricsType = 0x03 // movement/part name (e.g. "Adagio")
	LyricsTypeEvent      LyricsType = 0x04 // events (e.g. "Don Quijote enters the stage")
	LyricsTypeChord      LyricsType = 0x05 // chord (e.g. "Bb F Fsus")
	LyricsTypeTrivia     LyricsType = 0x06 // trivia/'pop up' information
	LyricsTypeWebUrl     LyricsType = 0x07 // URLs to webpages
	LyricsTypeImageUrl   LyricsType = 0x08 // URLs to images
)

type ChannelType byte

const (
	ChannelTypeOther       ChannelType = 0x00 // Other
	ChannelTypeMaster      ChannelType = 0x01 // Master volume
	ChannelTypeFrontRight  ChannelType = 0x02 // Front right
	ChannelTypeFrontLeft   ChannelType = 0x03 // Front left
	ChannelTypeBackRight   ChannelType = 0x04 // Back right
	ChannelTypeBackLeft    ChannelType = 0x05 // Back left
	ChannelTypeFrontCentre ChannelType = 0x06 // Front centre
	ChannelTypeBackCentre  ChannelType = 0x07 // Back centre
	ChannelTypeSubwoofer   ChannelType = 0x08 // Subwoofer
)

type EqualizationMethod byte

const (
	// No interpolation is made. A jump from one adjustment level to
	// another occurs in the middle between two adjustment points.
	EqualizationMethodBand EqualizationMethod = 0x00
	// Interpolation between adjustment points is linear.
	EqualizationMethodLinear EqualizationMethod = 0x01
)

type PictureType byte

const (
	PictureTypeOther         PictureType = 0x00 // Other
	PictureTypeFileIcon32    PictureType = 0x01 // 32x32 pixels 'file icon' (PNG only)
	PictureTypeFileIcon      PictureType = 0x02 // Other file icon
	PictureTypeFrontCover    PictureType = 0x03 // Cover (front)
	PictureTypeBackCover     PictureType = 0x04 // Cover (back)
	PictureTypeLeaflet       PictureType = 0x05 // Leaflet page
	PictureTypeMedia         PictureType = 0x06 // Media (e.g. label side of CD)
	PictureTypePerformer     PictureType = 0x07 // Lead artist/lead performer/soloist
	PictureTypeArtist        PictureType = 0x08 // Artist/performer
	PictureTypeConductor     PictureType = 0x09 // Conductor
	PictureTypeBand          PictureType = 0x0A // Band/Orchestra
	PictureTypeComposer      PictureType = 0x0B // Composer
	PictureTypeWriter        PictureType = 0x0C // Lyricist/text writer
	PictureTypeLocation      PictureType = 0x0D // Recording Location
	PictureTypeRecording     PictureType = 0x0E // During recording
	PictureTypePerformance   PictureType = 0x0F // During performance
	PictureTypeSnapshot      PictureType = 0x10 // Movie/video screen capture
	PictureTypeFinish        PictureType = 0x11 // A bright coloured fish
	PictureTypeIllustration  PictureType = 0x12 // Illustration
	PictureTypeArtistLogo    PictureType = 0x13 // Band/artist logotype
	PictureTypePublisherLogo PictureType = 0x14 // Publisher/Studio logotype
)

type DeliveryMethod byte

const (
	DeliveryMethodOther        DeliveryMethod = 0x00 // Other
	DeliveryMethodStdCD        DeliveryMethod = 0x01 // Standard CD album with other songs
	DeliveryMethodCompressedCD DeliveryMethod = 0x02 // Compressed audio on CD
	DeliveryMethodDownload     DeliveryMethod = 0x03 // File over the Internet
	DeliveryMethodStream       DeliveryMethod = 0x04 // Stream over the Internet
	DeliveryMethodNoteSheets   DeliveryMethod = 0x05 // As note sheets
	DeliveryMethodBook         DeliveryMethod = 0x06 // As note sheets in a book with other sheets
	DeliveryMethodOtherMedia   DeliveryMethod = 0x07 // Music on other media
	DeliveryMethodMerchandise  DeliveryMethod = 0x08 // Non-musical merchandise
)
