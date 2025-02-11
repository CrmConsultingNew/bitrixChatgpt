package reports

/*
// GetCallsInfo retrieves call information for specified users and date range.
func GetCallsInfo(users []string, start, finish string) ([]map[string]interface{}, error) {
	filter := map[string]interface{}{
		"CALL_FAILED_CODE": 200,
		">CALL_START_DATE": start,
		"<CALL_START_DATE": finish,
		"PORTAL_USER_ID":   users,
	}
	return GetCallList(filter)
}

func GetCallsVals(users []string, start, finish string) (map[string]map[string]map[string]interface{}, error) {
	vals := make(map[string]map[string]map[string]interface{})
	calls, err := GetCallsInfo(users, start, finish)
	if err != nil {
		return nil, err
	}

	for _, call := range calls {
		userID := fmt.Sprintf("%v", call["PORTAL_USER_ID"])
		callType := fmt.Sprintf("%v", call["CALL_TYPE"])
		duration := int(call["CALL_DURATION"].(float64))
		date, _ := time.Parse(time.RFC3339, call["CALL_START_DATE"].(string))
		dateFormatted := date.Format("15:04:05 02.01.2006")

		if _, ok := vals[userID]; !ok {
			vals[userID] = make(map[string]map[string]interface{})
		}
		if _, ok := vals[userID][callType]; !ok {
			vals[userID][callType] = map[string]interface{}{
				"count": 1,
				"dur":   duration,
				"first": dateFormatted,
				"last":  dateFormatted,
			}
		} else {
			vals[userID][callType]["count"] = vals[userID][callType]["count"].(int) + 1
			vals[userID][callType]["dur"] = vals[userID][callType]["dur"].(int) + duration
			vals[userID][callType]["last"] = dateFormatted
		}
	}
	return vals, nil
}

// GetLeadVals processes lead data and organizes it into a structured format.
func GetLeadVals(leads []map[string]interface{}) map[string]map[string]interface{} {
	vals := make(map[string]map[string]interface{})
	for _, lead := range leads {
		userID := fmt.Sprintf("%v", lead["ASSIGNED_BY_ID"])
		leadID := fmt.Sprintf("%v", lead["ID"])
		leadArr := map[string]string{
			"ID":    leadID,
			"link":  fmt.Sprintf("https://%s/crm/lead/details/%s/", B24_HOST, leadID),
			"title": fmt.Sprintf("%v", lead["TITLE"]),
		}

		if _, ok := vals[userID]; !ok {
			vals[userID] = map[string]interface{}{
				"count": 1,
				"leads": []map[string]string{leadArr},
			}
		} else {
			vals[userID]["count"] = vals[userID]["count"].(int) + 1
			vals[userID]["leads"] = append(vals[userID]["leads"].([]map[string]string), leadArr)
		}
	}
	return vals
}

// GetRegLeadInfo retrieves leads registered in the specified date range.
func GetRegLeadInfo(users []string, start, finish string) ([]map[string]interface{}, error) {
	params := map[string]interface{}{
		fmt.Sprintf(">%s", B24_LEAD_REGDATE_FIELD): start,
		fmt.Sprintf("<%s", B24_LEAD_REGDATE_FIELD): finish,
		"ASSIGNED_BY_ID": users,
	}
	order := map[string]string{
		B24_LEAD_REGDATE_FIELD: "ASC",
	}
	return GetLeadList(params, order)
}

// GetRegLeadVals processes and organizes registered lead data.
func GetRegLeadVals(users []string, start, finish string) (map[string]map[string]interface{}, error) {
	leads, err := GetRegLeadInfo(users, start, finish)
	if err != nil {
		return nil, err
	}
	return GetLeadVals(leads), nil
}

// GetDealVals processes and organizes deal data.
func GetDealVals(users []string, start, finish string) (map[string]map[string]map[string]interface{}, error) {
	vals := make(map[string]map[string]map[string]interface{})
	deals, err := GetDealInfo(users, start, finish)
	if err != nil {
		return nil, err
	}

	for _, deal := range deals {
		userID := fmt.Sprintf("%v", deal["ASSIGNED_BY_ID"])
		dealID := fmt.Sprintf("%v", deal["ID"])
		stage := fmt.Sprintf("%v", deal["STAGE_ID"])
		sum := deal["OPPORTUNITY"].(float64)
		dealArr := map[string]string{
			"ID":    dealID,
			"link":  fmt.Sprintf("https://%s/crm/deal/details/%s/", B24_HOST, dealID),
			"title": fmt.Sprintf("%v", deal["TITLE"]),
		}

		if vals[userID] == nil {
			vals[userID] = make(map[string]map[string]interface{})
		}
		if vals[userID][stage] == nil {
			vals[userID][stage] = map[string]interface{}{
				"count": 1,
				"sum":   sum,
				"deals": []map[string]string{dealArr},
			}
		} else {
			vals[userID][stage]["count"] = vals[userID][stage]["count"].(int) + 1
			vals[userID][stage]["sum"] = vals[userID][stage]["sum"].(float64) + sum
			vals[userID][stage]["deals"] = append(vals[userID][stage]["deals"].([]map[string]string), dealArr)
		}
	}
	return vals, nil
}

// GetTaskVals processes and organizes task data.
func GetTaskVals(role string, users []string, start, finish string) (map[string]map[string]map[string]interface{}, error) {
	vals := make(map[string]map[string]map[string]interface{})
	tasks, err := GetTaskInfo(role, users, start, finish)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		userID := fmt.Sprintf("%v", task[role])
		taskID := fmt.Sprintf("%v", task["ID"])
		tags := task["TAGS"].([]interface{})
		taskArr := map[string]string{
			"ID":    taskID,
			"link":  fmt.Sprintf("https://%s/company/personal/user/%s/tasks/task/view/%s/", B24_HOST, userID, taskID),
			"title": fmt.Sprintf("%v", task["TITLE"]),
		}
		taskType := "reg"
		if strings.HasPrefix(fmt.Sprintf("%v", task["TITLE"]), "Настройка") {
			taskType = "set"
		} else if strings.HasPrefix(fmt.Sprintf("%v", task["TITLE"]), "Размещение") {
			taskType = "pub"
		}
		tag := "reg"
		for _, t := range tags {
			if t == B24_TASK_SUCCTAG {
				tag = "succ"
				break
			}
		}

		if vals[userID] == nil {
			vals[userID] = make(map[string]map[string]interface{})
		}
		if vals[userID][taskType] == nil {
			vals[userID][taskType] = make(map[string]interface{})
		}
		if vals[userID][taskType][tag] == nil {
			vals[userID][taskType][tag] = map[string]interface{}{
				"count": 1,
				"tasks": []map[string]string{taskArr},
			}
		} else {
			vals[userID][taskType][tag]["count"] = vals[userID][taskType][tag]["count"].(int) + 1
			vals[userID][taskType][tag]["tasks"] = append(vals[userID][taskType][tag]["tasks"].([]map[string]string), taskArr)
		}
	}
	return vals, nil
}
*/
