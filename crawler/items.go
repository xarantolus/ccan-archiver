package crawler

import (
	"time"
)

// This file fontains all additional items that aren't listed on ccan.de. This includes Freeware keys, linux and mac versions and US versions (only german games are linked)

// Warning: All timestamps are estimates from listed date at the website (e.g. Published Year). They might be wrong
var additionalItems = []CCANItem{
	// Freeware Key for Clonk Endeavour; Copy this file to the directory that clonk.exe is installed in
	CCANItem{
		Name:          "Freeware",
		Date:          time.Date(2004, 01, 01, 0, 0, 0, 0, time.Local),
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Key",
		Engine:        "CE",
		DownloadLink:  "http://www.clonkx.de/endeavour/Freeware.c4k",
	},

	// Clonk Planet 'US' Version - only the German one is linked
	CCANItem{
		Name:          "Clonk Planet US",
		Date:          time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local), // Published in 2000
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Engine",
		Engine:        "CP",
		DownloadLink:  "http://www.clonkx.de/planet/cp465us_free.exe",
	},

	// Freeware Key/Instructions for Clonk Planet
	CCANItem{
		Name:          "Freeware Key Clonk Planet DE",
		Date:          time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Key",
		Engine:        "CP",
		DownloadLink:  "http://www.clonkx.de/planet/cp_freeware_de.txt",
	},
	CCANItem{
		Name:          "Freeware Key Clonk Planet US",
		Date:          time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Key",
		Engine:        "CP",
		DownloadLink:  "http://www.clonkx.de/planet/cp_freeware_us.txt",
	},

	// For Clonk Rage, only the Windows version is linked. The following items link the linux tar and mac zip so these versions will also be archived:
	CCANItem{
		Name:          "Clonk Rage Linux",
		Date:          time.Date(2014, 5, 4, 23, 25, 52, 0, time.Local), // Last-Modified: Sun, 04 May 2014 23:25:52 GMT - time.Local is wrong but it's close enough i guess?
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Engine",
		Engine:        "CR",
		DownloadLink:  "http://www.clonkx.de/rage/cr_full_linux.tar.bz2",
	},
	CCANItem{
		Name:          "Clonk Rage Mac",
		Date:          time.Date(2014, 5, 4, 23, 27, 0, 0, time.Local), // Last-Modified: Sun, 04 May 2014 23:27:00 GMT - time.Local is wrong but it's close enough i guess?
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Engine",
		Engine:        "CR",
		DownloadLink:  "http://www.clonkx.de/rage/cr_full_mac.zip",
	},

	// Clonk 3 Radikal - US version
	CCANItem{
		Name:          "Clonk 3 Radikal US",
		Date:          time.Date(1996, 0, 0, 0, 0, 0, 0, time.Local), // Published in 1996
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Engine",
		Engine:        "C3",
		DownloadLink:  "http://www.clonkx.de/classics/clonk34us.zip",
	},

	// Clonk 4 - US version
	CCANItem{
		Name:          "Clonk US",                                    // The german Clonk 4 entry is called "Clonk.zip"
		Date:          time.Date(1996, 0, 0, 0, 0, 0, 0, time.Local), // Published in 1996
		DownloadCount: 1,
		Author:        "Redwolf Design",
		Votes:         0,
		Category:      "Engine",
		Engine:        "C4.25",
		DownloadLink:  "http://www.clonkx.de/classics/clonk407us.zip",
	},
}
