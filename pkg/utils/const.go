package utils

const (
	LocalID     string = "localID"
	LocalUser   string = "localUser"
	LocalFile   string = "localFile"
	LocalDTO    string = "localDTO"
	LocalFilter string = "localFilter"

	ParamID   string = "id"
	ParamMail string = "email"

	PGAuth                string = "Auth"
	PGProfile             string = "Profile"
	PGAuthProfile                = PGAuth + "." + PGProfile
	PGEmployee            string = "Employee"
	PGDepartment          string = "Department"
	PGPosition            string = "Position"
	PGLevel               string = "Level"
	PGCatetory            string = "Category"
	PGType                string = "Type"
	PGStatus              string = "Status"
	PGPlatform            string = "Platform"
	PGCompetences         string = "Competences"
	PGCategories          string = "Categories"
	PGEmployeeDepartment         = PGEmployee + "." + PGDepartment
	PGEmployeePosition           = PGEmployee + "." + PGPosition
	PGEmployeeLevel              = PGEmployee + "." + PGLevel
	PGCatetoryType               = PGCatetory + "." + PGType
	PGCompetencesCategory        = PGCompetences + "." + PGCatetory
)
