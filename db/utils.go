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

func fieldStringPrefix(fields []string, prefix string) string {
	fieldString := ""
	for i, field := range fields {
		if i == 0 {
			fieldString = prefix + "." + field
		} else {
			fieldString = fieldString + ", " + prefix + "." + field
		}
	}
	return fieldString
}

func valuesString(fields []string) string {
	values := ""
	for i, field := range fields {
		if i == 1 {
			values = "@" + field
		} else if i != 0 {
			values = values + ", @" + field
		}
	}
	return values
}

func valuesStringNoID(fields []string) string {
	values := ""
	for i, field := range fields {
		if i == 0 {
			values = "@" + field
		} else {
			values = values + ", @" + field
		}
	}
	return values
}

func fieldStringInsert(fields []string) string {
	fieldString := ""
	for i, field := range fields {
		if i == 1 {
			fieldString = field
		} else if i != 0 {
			fieldString = fieldString + ", " + field
		}
	}
	return fieldString
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

func updatesStringNoID(fields []string) string {
	values := ""
	for i, field := range fields {
		if i == 0 {
			values = field + " = @" + field
		} else {
			values = values + ", " + field + " = @" + field
		}
	}
	return values
}

