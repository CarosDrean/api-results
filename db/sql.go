package db

import (
	"database/sql"
)

// db es la base de datos global
var DB *sql.DB

type nameQuery string

type queryConfig struct {
	Name string
	Q    string
}

type TableDB struct {
	Name   string
	Fields []string
}

var queryResultService = map[string]*queryConfig{
	"getOld": {Q: "select scfv.v_Value1 from servicecomponent sc " +
		"inner join servicecomponentfields scf on sc.v_ServiceComponentId = scf.v_ServiceComponentId " +
		"inner join servicecomponentfieldvalues scfv on scf.v_ServiceComponentFieldsId = scfv.v_ServiceComponentFieldsId " +
		"where sc.v_ServiceId = '%s' and sc.v_ComponentId = '%s' and scf.v_ComponentFieldId = '%s'"},
	"get": {Q: "select scfv.v_Value1 from service s " +
		"inner join protocol p on s.v_ProtocolId = p.v_ProtocolId " +
		"inner join protocolcomponent pc on p.v_ProtocolId = pc.v_ProtocolId " +
		"inner join servicecomponent sc on s.v_ServiceId = sc.v_ServiceId " +
		"inner join servicecomponentfields scf on sc.v_ServiceComponentId = scf.v_ServiceComponentId " +
		"inner join servicecomponentfieldvalues scfv on scf.v_ServiceComponentFieldsId = scfv.v_ServiceComponentFieldsId " +
		"where s.v_ServiceId = '%s' and pc.v_ComponentId = '%s' and scf.v_ComponentFieldId = '%s'"},
}

var queryStatistics = map[string]*queryConfig{
	"getDisease": {Q: "select s." + service.Fields[0] + ", s." + service.Fields[3] + ", " +
		"p." + person.Fields[0] + ", pr." + protocol.Fields[0] + ", s." + service.Fields[6] + ", p." + person.Fields[1] + ", p." + person.Fields[3] +
		", p." + person.Fields[4] + ", p." + person.Fields[5] + ", p." + person.Fields[6] + ", p." + person.Fields[7] + ", p." + person.Fields[8] +
		", d.v_Name, c.v_Name, sp.v_Value1 from service s " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"left join protocol pr on s.v_ProtocolId = pr.v_ProtocolId " +
		"left join organization o on pr.v_CustomerOrganizationId = o.v_OrganizationId " +
		"inner join diagnosticrepository dr on s.v_ServiceId = dr.v_ServiceId " +
		"inner join component c on dr.v_ComponentId = c.v_ComponentId " +
		"inner join systemparameter sp on sp.i_GroupId = 116 and c.i_CategoryId = sp.i_ParameterId " +
		"inner join diseases d on dr.v_DiseasesId = d.v_DiseasesId " +
		"where dr.i_IsDeleted = 0 and s.i_ServiceStatusId = 3 and pr.v_ProtocolId = '%s' " +
		"and s.d_ServiceDate >= CAST('%s' as date) and s.d_ServiceDate <= CAST('%s' as date) " +
		"order by s.d_ServiceDate desc"},
	"getAllDiseaseDate": {Q: "select s." + service.Fields[0] + ", s." + service.Fields[3] + ", " +
		"p." + person.Fields[0] + ", pr." + protocol.Fields[0] + ", s." + service.Fields[6] + ", p." + person.Fields[1] + ", p." + person.Fields[3] +
		", p." + person.Fields[4] + ", p." + person.Fields[5] + ", p." + person.Fields[6] + ", p." + person.Fields[7] + ", p." + person.Fields[8] +
		", d.v_Name, c.v_Name, sp.v_Value1 from service s " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"left join protocol pr on s.v_ProtocolId = pr.v_ProtocolId " +
		"left join organization o on pr.v_CustomerOrganizationId = o.v_OrganizationId " +
		"inner join diagnosticrepository dr on s.v_ServiceId = dr.v_ServiceId " +
		"inner join component c on dr.v_ComponentId = c.v_ComponentId " +
		"inner join systemparameter sp on sp.i_GroupId = 116 and c.i_CategoryId = sp.i_ParameterId " +
		"inner join diseases d on dr.v_DiseasesId = d.v_DiseasesId " +
		"where dr.i_IsDeleted = 0 and s.i_ServiceStatusId = 3 " +
		"and s.d_ServiceDate >= CAST('%s' as date) and s.d_ServiceDate <= CAST('%s' as date) " +
		"order by s.d_ServiceDate desc"},
}

var systemParameter = TableDB{
	Name:   "dbo.systemparameter",
	Fields: []string{"i_GroupId", "i_ParameterId", "v_Value1"},
}

var querySystemParameter = map[string]*queryConfig{
	"getGroup": {Q: "select " + fieldString(systemParameter.Fields) + " from " + systemParameter.Name + " where " + calendar.Fields[0] + " = %s;"},
}

var component = TableDB{
	Name:   "dbo.component",
	Fields: []string{"v_ComponentId", "v_Name", "i_CategoryId", "i_IsDeleted"},
}

var queryComponent = map[string]*queryConfig{
	"getCategory": {Q: "select " + fieldString(component.Fields) + " from " + component.Name + " where " + calendar.Fields[2] + " = %s;"},
}

var organization = TableDB{
	Name:   "dbo.organization",
	Fields: []string{"v_OrganizationId", "v_Name"},
}

var calendar = TableDB{
	Name:   "dbo.calendar",
	Fields: []string{"v_CalendarId", "v_ServiceId", "i_CalendarStatusId"},
}

var QueryCalendar = map[string]*queryConfig{
	"getServiceID": {Q: "select " + fieldString(calendar.Fields) + " from " + calendar.Name + " where " + calendar.Fields[1] + " = '%s';"},
}

var service = TableDB{
	Name: "dbo.service",
	Fields: []string{"v_ServiceId", "v_PersonId", "v_ProtocolId", "d_ServiceDate", "i_ServiceStatusId", "i_isDeleted",
		"i_AptitudeStatusId"},
}

var QueryService = map[nameQuery]*queryConfig{
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
	"get":            {Q: "select " + fieldString(service.Fields) + " from " + service.Name + " where " + service.Fields[0] + " = '%s';"},
}

var protocol = TableDB{
	Name:   "dbo.protocol",
	Fields: []string{"v_ProtocolId", "v_Name", "v_CustomerOrganizationId", "v_EmployerLocationId", "i_IsDeleted", "i_EsoTypeId"},
}

var QueryProtocol = map[string]*queryConfig{
	"getLocation":     {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[3] + " = '%s';"},
	"getOrganization": {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[2] + " = '%s';"},
	"get":             {Q: "select " + fieldString(protocol.Fields) + " from " + protocol.Name + " where " + protocol.Fields[0] + " = '%s';"},
}

var QueryOrganization = map[string]*queryConfig{
	"list": {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " order by " + organization.Fields[1] + " asc;"},
	"get":  {Q: "select " + fieldString(organization.Fields) + " from " + organization.Name + " where " + organization.Fields[0] + " = '%s';"},
}

var protocolSystemUser = TableDB{
	Name:   "dbo.protocolsystemuser",
	Fields: []string{"v_ProtocolSystemUserId", "i_SystemUserId", "v_ProtocolId"},
}

var QueryProtocolSystemUser = map[string]*queryConfig{
	"getSystemUserID": {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[1] + " = '%s';"},
	"get":             {Q: "select " + fieldString(protocolSystemUser.Fields) + " from " + protocolSystemUser.Name + " where " + protocolSystemUser.Fields[0] + " = '%s';"},
	"insert":          {Q: "insert into " + protocolSystemUser.Name + " (" + fieldString(protocolSystemUser.Fields) + ") values (" + valuesStringNoID(protocolSystemUser.Fields) + ");"},
}

var location = TableDB{
	Name:   "dbo.location",
	Fields: []string{"v_LocationId", "v_OrganizationId", "v_Name", "i_IsDeleted"},
}

var queryLocation = map[string]*queryConfig{
	"getOrganizationID": {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[1] + " = '%s';"},
	"get":               {Q: "select " + fieldString(location.Fields) + " from " + location.Name + " where " + location.Fields[0] + " = '%s';"},
}

var sequential = TableDB{
	Name:   "dbo.secuential",
	Fields: []string{"i_NodeId", "i_TableId", "i_SecuentialId"},
}

var querySequential = map[string]*queryConfig{
	"get":    {Q: "select " + fieldString(sequential.Fields) + " from " + sequential.Name + " where " + sequential.Fields[0] + " = %d and " + sequential.Fields[1] + " = %d;"},
	"insert": {Q: "insert into " + sequential.Name + " (" + fieldString(sequential.Fields) + ") values (" + valuesStringNoID(sequential.Fields) + ");"},
	"update": {Q: "update " + sequential.Name + " set " + updatesStringNoID(sequential.Fields) + " where " + sequential.Fields[0] + " = %d and " + sequential.Fields[1] + " = %d;"},
}

var user = TableDB{
	Name:   "dbo.systemuser",
	Fields: []string{"i_SystemUserId", "v_PersonId", "v_UserName", "v_Password", "i_SystemUserTypeId", "i_IsDeleted"},
}

var QuerySystemUser = map[string]*queryConfig{
	"getUserName":    {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[2] + " = '%s';"},
	"get":            {Q: "select " + fieldString(user.Fields) + " from " + user.Name + " where " + user.Fields[0] + " = %s;"},
	"list":           {Q: "select " + fieldString(user.Fields) + " from " + user.Name + ";"},
	"insert":         {Q: "insert into " + user.Name + " (" + fieldString(user.Fields) + ") values (" + valuesStringNoID(user.Fields) + ");"},
	"updatePassword": {Q: "update " + user.Name + " set v_Password = @Password where " + user.Fields[0] + " = %s;"},
	"update":         {Q: "update " + user.Name + " set " + updatesString(user.Fields) + " where " + user.Fields[0] + " = @ID;"},
	"delete":         {Q: "update " + user.Name + " set " + user.Fields[5] + " = 1 where " + user.Fields[0] + " = @ID;"},
}

var person = TableDB{
	Name: "dbo.person",
	Fields: []string{"v_PersonId", "v_DocNumber", "v_Password", "v_FirstName", "v_FirstLastName", "v_SecondLastName",
		"v_Mail", "i_SexTypeId", "d_Birthdate", "i_IsDeleted"},
}

var QueryPerson = map[string]*queryConfig{
	"get":            {Q: "select " + fieldString(person.Fields) + " from " + person.Name + " where " + person.Fields[0] + " = '%s';"},
	"getDNI":         {Q: "select " + fieldString(person.Fields) + " from " + person.Name + " where " + person.Fields[1] + " = '%s';"},
	"list":           {Q: "select " + fieldString(person.Fields) + " from " + person.Name + ";"},
	"insert":         {Q: "insert into " + person.Name + " (" + fieldString(person.Fields) + ") values (" + valuesStringNoID(person.Fields) + ");"},
	"update":         {Q: "update " + person.Name + " set " + updatesString(person.Fields) + " where " + person.Fields[0] + " = '%s';"},
	"updatePassword": {Q: "update " + person.Name + " set v_Password = @Password where " + person.Fields[0] + " = '%s';"},
}
