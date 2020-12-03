package db


func fieldString(fields []string) string {
	fieldString := ""
	for i, field := range fields {
		if i == 0 {
			fieldString = field
		} else {
			fieldString = fieldString + ", " + field
		}
	}
	return fieldString
}

func valuesString(fields []string) string {
	values := ""
	for i := range fields {
		if i == 0 {
			values = "?"
		} else {
			values = values + ", ?"
		}
	}
	return values
}

func updatesString(fields []string) string {
	values := ""
	for i, field := range fields {
		if i == 1 {
			values = field + " = @" + field
		} else if i != 0 {
			values = values + ", " + field + " = @" + field
		}
	}
	return values
}
