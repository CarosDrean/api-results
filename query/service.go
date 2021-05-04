package query

import "github.com/CarosDrean/api-results.git/models"

var service = models.TableDB{
	Name: "service",
	Fields: []string{"v_ServiceId", "v_PersonId", "v_ProtocolId", "d_ServiceDate", "i_ServiceStatusId", "i_isDeleted",
		"i_AptitudeStatusId"},
}

var Service = models.QueryDB{
	"getPersonID": {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[1] + " = '%s' order by " + service.Fields[3] + " desc;"},
	"getProtocol": {Q: "select " + fieldStringPrefix(service.Fields, "s") + " from " + service.Name +
		" s inner join calendar c on s.v_ServiceId = c.v_ServiceId and c.i_CalendarStatusId != 4" +
		" where " + service.Fields[2] +
		" = '%s' and d_ServiceDate is not null order by " + service.Fields[3] + " desc;"},
	"getProtocolFilter": {Q: "select " + fieldStringPrefix(service.Fields, "s") + " from " + service.Name +
		" s inner join calendar c on s.v_ServiceId = c.v_ServiceId and c.i_CalendarStatusId != 4" +
		" inner join  protocol p on s.v_ProtocolId = p.v_ProtocolId " +
		" where s." + service.Fields[2] +
		" = '%s' and CAST(s." + service.Fields[3] + " as date) >= CAST('%s' as date) and CAST(s." + service.Fields[3] + " as date) <= CAST('%s' as date) and s." + service.Fields[3] +
		" is not null order by s." + service.Fields[3] + " desc;"},
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
		" as date) <= CAST('%s' as date) and s." + service.Fields[3] + " is not null" +
		" order by s.d_ServiceDate;"},
	// reestructurar query para nuevo model
	"listExamsDetailDate": {Q: "select s.v_ServiceId, pr.v_ProtocolId, l.v_LocationId,  o.v_OrganizationId, p.v_FirstLastName, p.v_SecondLastName, " +
		"p.v_FirstName, dh.v_Value1, p.v_DocNumber, o.v_Name, p.v_CurrentOccupation, s.d_ServiceDate, p.d_Birthdate, " +
		"pc.r_Price, c.v_Name, sc.r_Price, pr.v_Name, p.v_Mail, p.i_SexTypeId, s.i_AptitudeStatusId, pr.i_EsoTypeId " +
		"from service s " +
		"inner join servicecomponent sc on s.v_ServiceId = sc.v_ServiceId and sc.r_Price > 0 " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"inner join protocol pr on s.v_ProtocolId = pr.v_ProtocolId " +
		"inner join protocolcomponent pc on pr.v_ProtocolId = pc.v_ProtocolId and pc.r_Price > 0 " +
		"inner join component c on pc.v_ComponentId = c.v_ComponentId " +
		"inner join organization o on pr.v_CustomerOrganizationId = o.v_OrganizationId " +
		"inner join datahierarchy dh on p.i_DocTypeId = dh.i_ItemId and dh.i_GroupId = 106 " +
		"inner join location l on o.v_OrganizationId = l.v_OrganizationId " +
		"where CAST(s.d_ServiceDate as date) >= CAST('%s' as date) and CAST(s.d_ServiceDate as date) <= CAST('%s' as date)"},
	"getAllCovid": {Q: "select s.d_ServiceDate, p.v_FirstName, p.v_FirstLastName, p.v_SecondLastName, p.v_DocNumber, p.d_Birthdate, p.i_SexTypeId, g.v_Name, p.v_CurrentOccupation, c.v_Name, scfv.v_Value1 from service s " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"inner join protocol pr on s.v_ProtocolId = pr.v_ProtocolId and pr.i_IsDeleted = 0 " +
		"inner join groupoccupation g on pr.v_GroupOccupationId = g.v_GroupOccupationId " +
		"inner join protocolcomponent pc on s.v_ProtocolId = pc.v_ProtocolId and " +
		"(pc.v_ComponentId= 'N007-ME000000491' or pc.v_ComponentId= 'N009-ME000000567') and pc.i_IsDeleted = 0 " +
		"inner join component c on pc.v_ComponentId = c.v_ComponentId " +
		"inner join servicecomponent sc on s.v_ServiceId = sc.v_ServiceId " +
		"inner join servicecomponentfields scf on sc.v_ServiceComponentId = scf.v_ServiceComponentId and " +
		"(scf.v_ComponentFieldId = 'N007-MF000004612' or scf.v_ComponentFieldId = 'N009-MF000004572') " +
		"inner join servicecomponentfieldvalues scfv on scf.v_ServiceComponentFieldsId = scfv.v_ServiceComponentFieldsId " +
		"where p.v_DocNumber = '%s' " +
		"order by s.d_ServiceDate desc"},
	"getGesoFilter": {Q: "SELECT s.v_ServiceId, s.d_ServiceDate, p3.v_PersonId, p.v_ProtocolId, s.i_AptitudeStatusId, p3.v_DocNumber, " +
		"p3.v_FirstName, p3.v_FirstLastName, p3.v_SecondLastName,  p3.v_Mail, p3.i_SexTypeId, p3.d_Birthdate, " +
		"c.i_CalendarStatusId, c.d_CircuitStartDate, c.d_SalidaCM, gp.v_Name FROM organization o " +
		"INNER JOIN protocol p ON o.v_OrganizationId = p.v_EmployerOrganizationId " +
		"INNER JOIN service s ON p.v_ProtocolId = s.v_ProtocolId AND s.i_IsDeleted != 1 AND s.d_ServiceDate IS NOT NULL " +
		"INNER JOIN person p3 ON s.v_PersonId = p3.v_PersonId " +
		"INNER JOIN calendar c ON s.v_ServiceId = c.v_ServiceId AND c.i_CalendarStatusId != 4 " +
		"INNER JOIN groupoccupation gp ON p.v_GroupOccupationId = gp.v_GroupOccupationId " +
		"where o.v_OrganizationId = '%s' " +
		"and CAST(s." + service.Fields[3] + " as date) >= CAST('%s' as date) and CAST(s." + service.Fields[3] + " as date) <= CAST('%s' as date) and s." + service.Fields[3] +
		" is not null order by s." + service.Fields[3] + " desc;"},
}
