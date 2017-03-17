package pagerduty

import pagerduty "github.com/PagerDuty/go-pagerduty"

// Expands configured slice into []pagerduty.EscalationRule
func expandEscalationRules(configured interface{}) []pagerduty.EscalationRule {
	var rules []pagerduty.EscalationRule

	rawRules := configured.([]interface{})

	for _, raw := range rawRules {
		rule := raw.(map[string]interface{})

		escRule := &pagerduty.EscalationRule{
			Delay:   uint(rule["escalation_delay_in_minutes"].(int)),
			Targets: expandEscalationRuleTargets(rule["target"]),
		}

		rules = append(rules, *escRule)
	}

	return rules
}

func expandEscalationRuleTargets(configured interface{}) []pagerduty.APIObject {
	var targets []pagerduty.APIObject

	rawTargets := configured.([]interface{})

	for _, raw := range rawTargets {
		target := raw.(map[string]interface{})
		targets = append(targets, pagerduty.APIObject{
			ID:   target["id"].(string),
			Type: target["type"].(string),
		})
	}

	return targets
}

func flattenEscalationRuleTargets(targets []pagerduty.APIObject) []map[string]interface{} {
	var result []map[string]interface{}

	for _, target := range targets {
		result = append(result, map[string]interface{}{
			"id":   target.ID,
			"type": target.Type,
		})
	}

	return result
}

// Flattens a slice of []pagerduty.EscalationRule into []map[string]interface{}
func flattenEscalationRules(rules []pagerduty.EscalationRule) []map[string]interface{} {
	var result []map[string]interface{}

	for _, rule := range rules {
		r := map[string]interface{}{
			"id": rule.ID,
			"escalation_delay_in_minutes": rule.Delay,
			"target":                      flattenEscalationRuleTargets(rule.Targets),
		}

		result = append(result, r)
	}

	return result
}

// Expands configured slice into []pagerduty.ScheduleLayer
func expandScheduleLayers(configured interface{}) []pagerduty.ScheduleLayer {
	var layers []pagerduty.ScheduleLayer

	rawLayers := configured.([]interface{})

	for _, raw := range rawLayers {
		layer := raw.(map[string]interface{})
		rawUsers := layer["users"].([]interface{})
		rawRestrictions := layer["restriction"].([]interface{})

		sLayer := &pagerduty.ScheduleLayer{
			Name:                      layer["name"].(string),
			Start:                     layer["start"].(string),
			End:                       layer["end"].(string),
			RotationVirtualStart:      layer["rotation_virtual_start"].(string),
			RotationTurnLengthSeconds: uint(layer["rotation_turn_length_seconds"].(int)),
		}

		if layer["id"] != "" {
			sLayer.ID = layer["id"].(string)
		}

		for _, rawUser := range rawUsers {
			sLayer.Users = append(
				sLayer.Users,
				pagerduty.UserReference{
					User: pagerduty.APIObject{
						ID:   rawUser.(string),
						Type: "user_reference",
					},
				},
			)
		}

		for _, rawRestriction := range rawRestrictions {
			restriction := rawRestriction.(map[string]interface{})
			sLayer.Restrictions = append(
				sLayer.Restrictions,
				pagerduty.Restriction{
					Type:            restriction["type"].(string),
					StartTimeOfDay:  restriction["start_time_of_day"].(string),
					StartDayOfWeek:  uint(restriction["start_day_of_week"].(int)),
					DurationSeconds: uint(restriction["duration_seconds"].(int)),
				},
			)
		}

		layers = append(layers, *sLayer)
	}

	return layers
}

// Expands configured slice into []pagerduty.APIReference
func expandTeams(configured interface{}) []pagerduty.APIReference {
	var teams []pagerduty.APIReference

	rawTeams := configured.([]interface{})

	for _, raw := range rawTeams {
		team := &pagerduty.APIReference{
			ID:   raw.(string),
			Type: "team_reference",
		}

		teams = append(result, *team)
	}

	return teams
}

// Flattens a slice of []pagerduty.ScheduleLayer into []map[string]interface{}
func flattenScheduleLayers(layers []pagerduty.ScheduleLayer) []map[string]interface{} {
	var result []map[string]interface{}

	for _, layer := range layers {
		l := map[string]interface{}{
			"id":    layer.ID,
			"name":  layer.Name,
			"start": layer.Start,
			"end":   layer.End,
			"rotation_virtual_start":       layer.RotationVirtualStart,
			"rotation_turn_length_seconds": layer.RotationTurnLengthSeconds,
		}

		var users []string
		for _, u := range layer.Users {
			users = append(users, u.User.ID)
		}
		l["users"] = users

		var restrictions []map[string]interface{}
		for _, r := range layer.Restrictions {
			restriction := map[string]interface{}{
				"duration_seconds":  r.DurationSeconds,
				"start_time_of_day": r.StartTimeOfDay,
				"type":              r.Type,
			}

			if r.StartDayOfWeek > 0 {
				restriction["start_day_of_week"] = r.StartDayOfWeek
			}

			restrictions = append(restrictions, restriction)
		}

		l["restriction"] = restrictions

		result = append(result, l)
	}

	// Reverse the final result and return it
	var resultReversed []map[string]interface{}

	for i := len(result) - 1; i >= 0; i-- {
		resultReversed = append(resultReversed, result[i])
	}

	return resultReversed
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

// Expands attribute slice to incident urgency rule, returns it and true if successful
func expandIncidentUrgencyRule(incidentUrgencyList interface{}) (*pagerduty.IncidentUrgencyRule, bool) {
	i, ok := incidentUrgencyList.([]interface{})
	if !ok {
		return nil, false
	}

	m, ok := i[0].(map[string]interface{})
	if !ok || len(m) == 0 {
		return nil, false
	}

	iur := pagerduty.IncidentUrgencyRule{}
	if val, ok := m["type"]; ok {
		iur.Type = val.(string)
	}
	if val, ok := m["urgency"]; ok {
		iur.Urgency = val.(string)
	}
	if val, ok := m["during_support_hours"]; ok {
		iur.DuringSupportHours = expandIncidentUrgencyType(val)
	}
	if val, ok := m["outside_support_hours"]; ok {
		iur.OutsideSupportHours = expandIncidentUrgencyType(val)
	}

	return &iur, true
}

// Expands attribute to inline model
func expandActionInlineModel(inlineModelVal interface{}) *pagerduty.InlineModel {
	inlineModel := pagerduty.InlineModel{}

	if slice, ok := inlineModelVal.([]interface{}); ok && len(slice) == 1 {
		m := slice[0].(map[string]interface{})

		if val, ok := m["type"]; ok {
			inlineModel.Type = val.(string)
		}
		if val, ok := m["name"]; ok {
			inlineModel.Name = val.(string)
		}
	}

	return &inlineModel
}

// Expands attribute into incident urgency type
func expandIncidentUrgencyType(attribute interface{}) *pagerduty.IncidentUrgencyType {
	ict := pagerduty.IncidentUrgencyType{}

	slice := attribute.([]interface{})
	if len(slice) != 1 {
		return &ict
	}

	m := slice[0].(map[string]interface{})

	if val, ok := m["type"]; ok {
		ict.Type = val.(string)
	}
	if val, ok := m["urgency"]; ok {
		ict.Urgency = val.(string)
	}

	return &ict
}

// Returns service's incident urgency rule as slice of length one and bool indicating success
func flattenIncidentUrgencyRule(service *pagerduty.Service) ([]interface{}, bool) {
	if service.IncidentUrgencyRule.Type == "" && service.IncidentUrgencyRule.Urgency == "" {
		return nil, false
	}

	m := map[string]interface{}{
		"type":    service.IncidentUrgencyRule.Type,
		"urgency": service.IncidentUrgencyRule.Urgency,
	}

	if dsh := service.IncidentUrgencyRule.DuringSupportHours; dsh != nil {
		m["during_support_hours"] = flattenIncidentUrgencyType(dsh)
	}
	if osh := service.IncidentUrgencyRule.OutsideSupportHours; osh != nil {
		m["outside_support_hours"] = flattenIncidentUrgencyType(osh)
	}

	return []interface{}{m}, true
}

func flattenIncidentUrgencyType(iut *pagerduty.IncidentUrgencyType) []interface{} {
	incidenUrgencyType := map[string]interface{}{
		"type":    iut.Type,
		"urgency": iut.Urgency,
	}
	return []interface{}{incidenUrgencyType}
}

// Expands attribute to support hours
func expandSupportHours(attribute interface{}) (sh *pagerduty.SupportHours) {
	if slice, ok := attribute.([]interface{}); ok && len(slice) >= 1 {
		m := slice[0].(map[string]interface{})
		sh = &pagerduty.SupportHours{}

		if val, ok := m["type"]; ok {
			sh.Type = val.(string)
		}
		if val, ok := m["time_zone"]; ok {
			sh.Timezone = val.(string)
		}
		if val, ok := m["start_time"]; ok {
			sh.StartTime = val.(string)
		}
		if val, ok := m["end_time"]; ok {
			sh.EndTime = val.(string)
		}
		if val, ok := m["days_of_week"]; ok {
			daysOfWeekInt := val.([]interface{})
			var daysOfWeek []uint

			for _, i := range daysOfWeekInt {
				daysOfWeek = append(daysOfWeek, uint(i.(int)))
			}

			sh.DaysOfWeek = daysOfWeek
		}
	}

	return
}

// Returns service's support hours as slice of length one
func flattenSupportHours(service *pagerduty.Service) []interface{} {
	if service.SupportHours == nil {
		return nil
	}

	m := map[string]interface{}{}

	if s := service.SupportHours; s != nil {
		m["type"] = s.Type
		m["time_zone"] = s.Timezone
		m["start_time"] = s.StartTime
		m["end_time"] = s.EndTime
		m["days_of_week"] = s.DaysOfWeek
	}

	return []interface{}{m}
}

// Expands attribute to scheduled action
func expandScheduledActions(input interface{}) (scheduledActions []pagerduty.ScheduledAction) {
	inputs := input.([]interface{})

	for _, i := range inputs {
		m := i.(map[string]interface{})
		sa := pagerduty.ScheduledAction{}

		if val, ok := m["type"]; ok {
			sa.Type = val.(string)
		}
		if val, ok := m["to_urgency"]; ok {
			sa.ToUrgency = val.(string)
		}
		if val, ok := m["at"]; ok {
			sa.At = *expandActionInlineModel(val)
		}

		scheduledActions = append(scheduledActions, sa)
	}

	return scheduledActions
}

// Returns service's scheduled actions
func flattenScheduledActions(service *pagerduty.Service) []interface{} {
	scheduledActions := []interface{}{}

	if sas := service.ScheduledActions; sas != nil {
		for _, sa := range sas {
			m := map[string]interface{}{}
			m["to_urgency"] = sa.ToUrgency
			m["type"] = sa.Type
			if at, ok := scheduledActionsAt(sa.At); ok {
				m["at"] = at
			}
			scheduledActions = append(scheduledActions, m)
		}
	}

	return scheduledActions
}

// Returns service's scheduled action's at attribute as slice of length one
func scheduledActionsAt(inlineModel pagerduty.InlineModel) ([]interface{}, bool) {
	if inlineModel.Type == "" || inlineModel.Name == "" {
		return nil, false
	}

	m := map[string]interface{}{"type": inlineModel.Type, "name": inlineModel.Name}
	return []interface{}{m}, true
}
