package query

import "github.com/CarosDrean/api-results.git/models"

var service = models.TableDB{
	Name: "dbo.service",
	Fields: []string{"v_ServiceId", "v_PersonId", "v_ProtocolId", "d_ServiceDate", "i_ServiceStatusId", "i_isDeleted",
		"i_AptitudeStatusId"},
}

var Service = models.QueryDB{
	"getPersonID": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[1] + " = '%s' order by " + service.Fields[3] + " desc;"},
	"getProtocol": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[2] +
		" = '%s' and d_ServiceDate is not null order by " + service.Fields[3] + " desc;"},
	"getProtocolFilter": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[2] +
		" = '%s' and CAST(" + service.Fields[3] + " as date) >= CAST('%s' as date) and CAST(" + service.Fields[3] + " as date) <= CAST('%s' as date) and " + service.Fields[3] +
		" is not null order by " + service.Fields[3] + " desc;"},
	"listDiseaseFilter": {Q: "select " + fieldStringPrefix(service.Fields, "s") + ", " + fieldStringPrefix(person.Fields, "pe") + ", " +
		fieldStringPrefix(protocol.Fields, "p") +
		", d.v_Name from dbo.service s " +
		"inner join person pe on s.v_PersonId = pe.v_PersonId " +
		"inner join protocol p on s.v_ProtocolId = p.v_ProtocolId " +
		"left join diagnosticrepository dr on s.v_ServiceId = dr.v_ServiceId " +
		"left join diseases d on dr.v_DiseasesId = d.v_DiseasesId where CAST(s." +
		service.Fields[3] + " as date) >= CAST('%s' as date) and CAST(s." + service.Fields[3] + " as date) <= CAST('%s' as date) " +
		"and s.i_ServiceStatusId = 3 and s." + service.Fields[3] +
		" is not null;"},
	"get": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[0] + " = '%s';"},
	"listDate": {Q: "select " + fieldStringPrefix(service.Fields, "s") + ", " + fieldStringPrefix(person.Fields, "pe") +
		", o.v_OrganizationId, o.v_Name, p.i_EsoTypeId from service s " +
		"inner join protocol p on s.v_ProtocolId = p.v_ProtocolId " +
		"inner join person pe on s.v_PersonId = pe.v_PersonId " +
		"inner join organization o on p.v_CustomerOrganizationId = o.v_OrganizationId " +
		" where CAST(s." + service.Fields[3] + " as date) >= CAST('%s' as date) and CAST(s." + service.Fields[3] +
		" as date) <= CAST('%s' as date) and s." + service.Fields[3] + " is not null;"},
}
