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

// http://www.sno.phy.queensu.ca/~phil/exiftool/TagNames/QuickTime.html
// http://shadowofged.blogspot.ca/2008/06/itunes-content-ratings.html

package itunes

import (
	"bytes"
	"strconv"
)

type MediaType int

const (
	MediaTypeHomeVideo  MediaType = 0  // 0 = Movie (deprecated, use 9 instead)
	MediaTypeMusic      MediaType = 1  // 1 = Normal (Music)
	MediaTypeAudiobook  MediaType = 2  // 2 = Audiobook
	MediaTypeBookmark   MediaType = 5  // 5 = Whacked Bookmark
	MediaTypeMusicVideo MediaType = 6  // 6 = Music Video
	MediaTypeMovie      MediaType = 9  // 9 = Short Film / Movie
	MediaTypeTVShow     MediaType = 10 // 10 = TV Show
	MediaTypeBooklet    MediaType = 11 // 11 = Booklet
	MediaTypeRingtone   MediaType = 14 // 14 = Ringtone
	MediaTypePodcast    MediaType = 21 // 21 = Podcast
)

func (x MediaType) String() string {
	switch x {
	case MediaTypeHomeVideo:
		return "Home Video"
	case MediaTypeMusic:
		return "Music"
	case MediaTypeAudiobook:
		return "Audiobook"
	case MediaTypeBookmark:
		return "Whacked Bookmark"
	case MediaTypeMusicVideo:
		return "Music Video"
	case MediaTypeMovie:
		return "Movie"
	case MediaTypeTVShow:
		return "TV Show"
	case MediaTypeBooklet:
		return "Booklet"
	case MediaTypeRingtone:
		return "Ringtone"
	case MediaTypePodcast:
		return "Podcast"
	default:
		buf := bytes.Buffer{}
		buf.WriteByte('(')
		buf.WriteString(strconv.FormatInt(int64(x), 10))
		buf.WriteByte(')')
		return buf.String()
	}
}

type RatingCode int

const (
	RatingCodeNone        RatingCode = 0 // 0 = None
	RatingCodeExplicit    RatingCode = 1 // 1 = Explicit
	RatingCodeClean       RatingCode = 2 // 2 = Clean
	RatingCodeExplicitOld RatingCode = 4 // 4 = Explicit (old)
)

type PlayGapMode int

const (
	PlayGapInsertGap PlayGapMode = 0 // Insert Gap
	PlayGapNoGap     PlayGapMode = 1 // No Gap
)

type AppleStoreAccountType int

const (
	AppleStoreAccountTypeITunes AppleStoreAccountType = 0
	AppleStoreAccountTypeAOL    AppleStoreAccountType = 1
)

type LocationRole int

const (
	LocationRoleShooting  LocationRole = 0
	LocationRoleReal      LocationRole = 1
	LocationRoleFictional LocationRole = 2
)

type AppleStoreCountry int

const (
	AppleStoreUSA AppleStoreCountry = 143441 // United States
	AppleStoreFRA AppleStoreCountry = 143442 // France
	AppleStoreDEU AppleStoreCountry = 143443 // Germany
	AppleStoreGBR AppleStoreCountry = 143444 // United Kingdom
	AppleStoreAUT AppleStoreCountry = 143445 // Austria
	AppleStoreBEL AppleStoreCountry = 143446 // Belgium
	AppleStoreFIN AppleStoreCountry = 143447 // Finland
	AppleStoreGRC AppleStoreCountry = 143448 // Greece
	AppleStoreIRL AppleStoreCountry = 143449 // Ireland
	AppleStoreITA AppleStoreCountry = 143450 // Italy
	AppleStoreLUX AppleStoreCountry = 143451 // Luxembourg
	AppleStoreNLD AppleStoreCountry = 143452 // Netherlands
	AppleStorePRT AppleStoreCountry = 143453 // Portugal
	AppleStoreESP AppleStoreCountry = 143454 // Spain
	AppleStoreCAN AppleStoreCountry = 143455 // Canada
	AppleStoreSWE AppleStoreCountry = 143456 // Sweden
	AppleStoreNOR AppleStoreCountry = 143457 // Norway
	AppleStoreDNK AppleStoreCountry = 143458 // Denmark
	AppleStoreCHE AppleStoreCountry = 143459 // Switzerland
	AppleStoreAUS AppleStoreCountry = 143460 // Australia
	AppleStoreNZL AppleStoreCountry = 143461 // New Zealand
	AppleStoreJPN AppleStoreCountry = 143462 // Japan
	AppleStoreHKG AppleStoreCountry = 143463 // Hong Kong
	AppleStoreSGP AppleStoreCountry = 143464 // Singapore
	AppleStoreCHN AppleStoreCountry = 143465 // China
	AppleStoreKOR AppleStoreCountry = 143466 // Republic of Korea
	AppleStoreIND AppleStoreCountry = 143467 // India
	AppleStoreMEX AppleStoreCountry = 143468 // Mexico
	AppleStoreRUS AppleStoreCountry = 143469 // Russia
	AppleStoreTWN AppleStoreCountry = 143470 // Taiwan
	AppleStoreVNM AppleStoreCountry = 143471 // Vietnam
	AppleStoreZAF AppleStoreCountry = 143472 // South Africa
	AppleStoreMYS AppleStoreCountry = 143473 // Malaysia
	AppleStorePHL AppleStoreCountry = 143474 // Philippines
	AppleStoreTHA AppleStoreCountry = 143475 // Thailand
	AppleStoreIDN AppleStoreCountry = 143476 // Indonesia
	AppleStorePAK AppleStoreCountry = 143477 // Pakistan
	AppleStorePOL AppleStoreCountry = 143478 // Poland
	AppleStoreSAU AppleStoreCountry = 143479 // Saudi Arabia
	AppleStoreTUR AppleStoreCountry = 143480 // Turkey
	AppleStoreARE AppleStoreCountry = 143481 // United Arab Emirates
	AppleStoreHUN AppleStoreCountry = 143482 // Hungary
	AppleStoreCHL AppleStoreCountry = 143483 // Chile
	AppleStoreNPL AppleStoreCountry = 143484 // Nepal
	AppleStorePAN AppleStoreCountry = 143485 // Panama
	AppleStoreLKA AppleStoreCountry = 143486 // Sri Lanka
	AppleStoreROU AppleStoreCountry = 143487 // Romania
	AppleStoreCZE AppleStoreCountry = 143489 // Czech Republic
	AppleStoreISR AppleStoreCountry = 143491 // Israel
	AppleStoreUKR AppleStoreCountry = 143492 // Ukraine
	AppleStoreKWT AppleStoreCountry = 143493 // Kuwait
	AppleStoreHRV AppleStoreCountry = 143494 // Croatia
	AppleStoreCRI AppleStoreCountry = 143495 // Costa Rica
	AppleStoreSVK AppleStoreCountry = 143496 // Slovakia
	AppleStoreLBN AppleStoreCountry = 143497 // Lebanon
	AppleStoreQAT AppleStoreCountry = 143498 // Qatar
	AppleStoreSVN AppleStoreCountry = 143499 // Slovenia
	AppleStoreCOL AppleStoreCountry = 143501 // Colombia
	AppleStoreVEN AppleStoreCountry = 143502 // Venezuela
	AppleStoreBRA AppleStoreCountry = 143503 // Brazil
	AppleStoreGTM AppleStoreCountry = 143504 // Guatemala
	AppleStoreARG AppleStoreCountry = 143505 // Argentina
	AppleStoreSLV AppleStoreCountry = 143506 // El Salvador
	AppleStorePER AppleStoreCountry = 143507 // Peru
	AppleStoreDOM AppleStoreCountry = 143508 // Dominican Republic
	AppleStoreECU AppleStoreCountry = 143509 // Ecuador
	AppleStoreHND AppleStoreCountry = 143510 // Honduras
	AppleStoreJAM AppleStoreCountry = 143511 // Jamaica
	AppleStoreNIC AppleStoreCountry = 143512 // Nicaragua
	AppleStorePRY AppleStoreCountry = 143513 // Paraguay
	AppleStoreURY AppleStoreCountry = 143514 // Uruguay
	AppleStoreMAC AppleStoreCountry = 143515 // Macau
	AppleStoreEGY AppleStoreCountry = 143516 // Egypt
	AppleStoreKAZ AppleStoreCountry = 143517 // Kazakhstan
	AppleStoreEST AppleStoreCountry = 143518 // Estonia
	AppleStoreLVA AppleStoreCountry = 143519 // Latvia
	AppleStoreLTU AppleStoreCountry = 143520 // Lithuania
	AppleStoreMLT AppleStoreCountry = 143521 // Malta
	AppleStoreMDA AppleStoreCountry = 143523 // Moldova
	AppleStoreARM AppleStoreCountry = 143524 // Armenia
	AppleStoreBWA AppleStoreCountry = 143525 // Botswana
	AppleStoreBGR AppleStoreCountry = 143526 // Bulgaria
	AppleStoreJOR AppleStoreCountry = 143528 // Jordan
	AppleStoreKEN AppleStoreCountry = 143529 // Kenya
	AppleStoreMKD AppleStoreCountry = 143530 // Macedonia
	AppleStoreMDG AppleStoreCountry = 143531 // Madagascar
	AppleStoreMLI AppleStoreCountry = 143532 // Mali
	AppleStoreMUS AppleStoreCountry = 143533 // Mauritius
	AppleStoreNER AppleStoreCountry = 143534 // Niger
	AppleStoreSEN AppleStoreCountry = 143535 // Senegal
	AppleStoreTUN AppleStoreCountry = 143536 // Tunisia
	AppleStoreUGA AppleStoreCountry = 143537 // Uganda
	AppleStoreAIA AppleStoreCountry = 143538 // Anguilla
	AppleStoreBHS AppleStoreCountry = 143539 // Bahamas
	AppleStoreATG AppleStoreCountry = 143540 // Antigua and Barbuda
	AppleStoreBRB AppleStoreCountry = 143541 // Barbados
	AppleStoreBMU AppleStoreCountry = 143542 // Bermuda
	AppleStoreVGB AppleStoreCountry = 143543 // British Virgin Islands
	AppleStoreCYM AppleStoreCountry = 143544 // Cayman Islands
	AppleStoreDMA AppleStoreCountry = 143545 // Dominica
	AppleStoreGRD AppleStoreCountry = 143546 // Grenada
	AppleStoreMSR AppleStoreCountry = 143547 // Montserrat
	AppleStoreKNA AppleStoreCountry = 143548 // St. Kitts and Nevis
	AppleStoreLCA AppleStoreCountry = 143549 // St. Lucia
	AppleStoreVCT AppleStoreCountry = 143550 // St. Vincent and The Grenadines
	AppleStoreTTO AppleStoreCountry = 143551 // Trinidad and Tobago
	AppleStoreTCA AppleStoreCountry = 143552 // Turks and Caicos
	AppleStoreGUY AppleStoreCountry = 143553 // Guyana
	AppleStoreSUR AppleStoreCountry = 143554 // Suriname
	AppleStoreBLZ AppleStoreCountry = 143555 // Belize
	AppleStoreBOL AppleStoreCountry = 143556 // Bolivia
	AppleStoreCYP AppleStoreCountry = 143557 // Cyprus
	AppleStoreISL AppleStoreCountry = 143558 // Iceland
	AppleStoreBHR AppleStoreCountry = 143559 // Bahrain
	AppleStoreBRN AppleStoreCountry = 143560 // Brunei Darussalam
	AppleStoreNGA AppleStoreCountry = 143561 // Nigeria
	AppleStoreOMN AppleStoreCountry = 143562 // Oman
	AppleStoreDZA AppleStoreCountry = 143563 // Algeria
	AppleStoreAGO AppleStoreCountry = 143564 // Angola
	AppleStoreBLR AppleStoreCountry = 143565 // Belarus
	AppleStoreUZB AppleStoreCountry = 143566 // Uzbekistan
	AppleStoreAZE AppleStoreCountry = 143568 // Azerbaijan
	AppleStoreYEM AppleStoreCountry = 143571 // Yemen
	AppleStoreTZA AppleStoreCountry = 143572 // Tanzania
	AppleStoreGHA AppleStoreCountry = 143573 // Ghana
	AppleStoreALB AppleStoreCountry = 143575 // Albania
	AppleStoreBEN AppleStoreCountry = 143576 // Benin
	AppleStoreBTN AppleStoreCountry = 143577 // Bhutan
	AppleStoreBFA AppleStoreCountry = 143578 // Burkina Faso
	AppleStoreKHM AppleStoreCountry = 143579 // Cambodia
	AppleStoreCPV AppleStoreCountry = 143580 // Cape Verde
	AppleStoreTCD AppleStoreCountry = 143581 // Chad
	AppleStoreCOG AppleStoreCountry = 143582 // Republic of the Congo
	AppleStoreFJI AppleStoreCountry = 143583 // Fiji
	AppleStoreGMB AppleStoreCountry = 143584 // Gambia
	AppleStoreGNB AppleStoreCountry = 143585 // Guinea-Bissau
	AppleStoreKGZ AppleStoreCountry = 143586 // Kyrgyzstan
	AppleStoreLAO AppleStoreCountry = 143587 // Lao People's Democratic Republic
	AppleStoreLBR AppleStoreCountry = 143588 // Liberia
	AppleStoreMWI AppleStoreCountry = 143589 // Malawi
	AppleStoreMRT AppleStoreCountry = 143590 // Mauritania
	AppleStoreFSM AppleStoreCountry = 143591 // Federated States of Micronesia
	AppleStoreMNG AppleStoreCountry = 143592 // Mongolia
	AppleStoreMOZ AppleStoreCountry = 143593 // Mozambique
	AppleStoreNAM AppleStoreCountry = 143594 // Namibia
	AppleStorePLW AppleStoreCountry = 143595 // Palau
	AppleStorePNG AppleStoreCountry = 143597 // Papua New Guinea
	AppleStoreSTP AppleStoreCountry = 143598 // Sao Tome and Principe
	AppleStoreSYC AppleStoreCountry = 143599 // Seychelles
	AppleStoreSLE AppleStoreCountry = 143600 // Sierra Leone
	AppleStoreSLB AppleStoreCountry = 143601 // Solomon Islands
	AppleStoreSWZ AppleStoreCountry = 143602 // Swaziland
	AppleStoreTJK AppleStoreCountry = 143603 // Tajikistan
	AppleStoreTKM AppleStoreCountry = 143604 // Turkmenistan'AppleStore
	AppleStoreZWE AppleStoreCountry = 143605 // Zimbabwe
)

// iTunes Genre category, genre and subgenre
// https://affiliate.itunes.apple.com/resources/documentation/genre-mapping/
// https://itunes.apple.com/WebObjects/MZStoreServices.woa/ws/genres
type GenreID int

// ID3v1 Genre id
type GenreCode byte
