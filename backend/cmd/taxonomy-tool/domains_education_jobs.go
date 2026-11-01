package main

func getEducationJobsDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 64. Root: Education, Vocational Training & Academic Tutoring
		{
			ID:            "root-education-training",
			Slug:          "education-vocational-training-tutoring",
			Code:          "EDU_TRAINING",
			Level:         1,
			CanonicalName: "Education, Vocational Training & Academic Tutoring",
			Definition:    "Technical and Vocational Education and Training (TVET Levels I–V), language instruction, software engineering bootcamps, and home academic tutoring.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "GraduationCap",
			Translations: map[string]string{
				"am": "ትምህርት፣ የሙያ ስልጠና (TVET) እና የማጠናከሪያ ትምህርት",
				"or": "Barnoota, Leenjii Ogummaa fi Qophii",
				"ti": "ትምህርቲ፣ ናይ ሞያ ስልጠናን መተዓረዪ ትምህርትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-tvet-vocational-courses",
					Slug:          "tvet-technical-vocational-training",
					Code:          "EDU_TVET",
					Level:         2,
					CanonicalName: "Technical & Vocational Training (TVET Levels I–V)",
					Definition:    "Accredited occupational certificate courses in automotive mechanics, electrical installation, construction masonry, and garment manufacturing.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የቴክኒክ እና ሙያ ስልጠና ኮርሶች (ደረጃ 1-5)",
						"or": "Leenjii Teeknikaa fi Ogummaa (TVET)",
						"ti": "ናይ ቴክኒክን ሞያን ስልጠናታት (TVET)",
					},
				},
			},
		},

		// 65. Root: Books, Academic Literature & Publishing
		{
			ID:            "root-books-publishing",
			Slug:          "books-academic-literature-publishing",
			Code:          "EDU_BOOKS",
			Level:         1,
			CanonicalName: "Books, Academic Literature & Publishing",
			Definition:    "University textbooks, national curriculum secondary school texts, Amharic/Afan Oromo/Tigrinya literature, and religious historical books.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "BookOpen",
			Translations: map[string]string{
				"am": "መጻሕፍት፣ የትምህርት ህትመቶች እና ስነ-ጽሑፍ",
				"or": "Kitaabota, Barreeffamoota Barnootaa fi Og-barruu",
				"ti": "መጻሕፍቲ፣ ናይ ትምህርቲ ሕትመታትን ስነ-ጽሑፍን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-academic-textbooks",
					Slug:          "university-school-academic-textbooks",
					Code:          "BOOK_TEXTBOOKS",
					Level:         2,
					CanonicalName: "School & University Academic Textbooks",
					Definition:    "Medicine, engineering, law, economics, and natural science reference volumes and curriculum books.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የዩኒቨርሲቲ እና የትምህርት ቤት መማሪያ መጻሕፍት",
						"or": "Kitaabota Barnootaa Yuunvarsiitii",
						"ti": "ናይ ዩኒቨርሲቲን ቤት ትምህርትን መጻሕፍቲ",
					},
				},
			},
		},

		// 66. Root: Professional Occupations & Recruitment
		{
			ID:            "root-jobs-recruitment",
			Slug:          "professional-occupations-careers-recruitment",
			Code:          "CAREER_OCCUPATIONS",
			Level:         1,
			CanonicalName: "Professional Occupations, Careers & Recruitment",
			Definition:    "Standardized occupational classifications for employment vacancies, professional candidate profiles, and skilled craftsmanship.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Users",
			Translations: map[string]string{
				"am": "ክፍት የስራ ቦታዎች፣ ሙያዎች እና የቅጥር እድሎች",
				"or": "Hojiiwwan, Ogummaa fi Qacarraa Hojjettootaa",
				"ti": "ክፍቲ ናይ ስራሕ ቦታታት፣ ሞያታትን ምልመላን",
			},
			Children: []TaxonomyNode{
				// Super Family 66.1: Engineering & Technical Professions
				{
					ID:            "sf-engineering-professions",
					Slug:          "civil-mechanical-software-engineering-roles",
					Code:          "JOB_ENGINEERING",
					Level:         2,
					CanonicalName: "Engineering, Architecture & Technical Professions",
					Definition:    "Civil engineers, electrical engineers, structural draughtsmen, and software engineers.",
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "የምህንድስና እና የቴክኒክ ባለሙያዎች",
						"or": "Ogeeyyii Injinariingii fi Teeknikaa",
						"ti": "ናይ መሃንድስነትን ቴክኒክን ሰብ ሞያ",
					},
					Children: []TaxonomyNode{
						// Level 3: Civil Engineers (ISCO 2142)
						{
							ID:            "cat-job-civil-engineer",
							Slug:          "civil-structural-engineers",
							Code:          "JOB_CIVIL_ENG",
							Level:         3,
							CanonicalName: "Civil & Structural Engineers (ሲቪል መሃንዲስ)",
							Definition:    "Design, supervision, and site project management for building, road, and water infrastructure works.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "ISCO_08", Code: "2142", Name: "Civil Engineers"},
							},
							Translations: map[string]string{
								"am": "የሲቪል እና የህንፃ ምህንድስና ባለሙያዎች",
								"or": "Injinara Siiviilii fi Ijaarsaa",
								"ti": "ናይ ሲቪልን ህንጻን መሃንዲሳት",
							},
						},
						// Level 3: Software Engineers & Developers (ISCO 2512)
						{
							ID:            "cat-job-software-engineer",
							Slug:          "software-engineers-developers",
							Code:          "JOB_SOFTWARE_ENG",
							Level:         3,
							CanonicalName: "Software Engineers, Full-Stack & App Developers",
							Definition:    "Backend, frontend, mobile app development, and cloud systems engineering.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "ISCO_08", Code: "2512", Name: "Software Developers"},
							},
							Translations: map[string]string{
								"am": "የሶፍትዌር ምህንድስና እና አፕሊኬሽን አበልጻጊዎች",
								"or": "Injinara Sooftiweerii fi Dizaayinaa",
								"ti": "ናይ ሶፍትዌር መሃንዲሳትን ኣልማዕትን",
							},
						},
					},
				},
				// Super Family 66.2: Healthcare & Clinical Professions
				{
					ID:            "sf-healthcare-professions",
					Slug:          "medical-doctors-nursing-healthcare-roles",
					Code:          "JOB_HEALTHCARE",
					Level:         2,
					CanonicalName: "Physicians, Nursing Staff & Allied Health Practitioners",
					Definition:    "General practitioners, registered clinic nurses, pharmacists, and laboratory technologists.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "ISCO_08", Code: "2211", Name: "Generalist Medical Practitioners"},
					},
					Translations: map[string]string{
						"am": "የህክምና ዶክተሮች፣ ነርሶች እና የጤና ባለሙያዎች",
						"or": "Ogeeyyii Fayyaa, Doktoorota fi Narsii",
						"ti": "ናይ ሕክምና ሓካይም፣ ነርሳትን ሰብ ሞያ ጥዕናን",
					},
				},
			},
		},
	}
}
