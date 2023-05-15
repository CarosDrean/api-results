package query

import "github.com/CarosDrean/api-results.git/models"

var ResultService = models.QueryDB{
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
		"where s.v_ServiceId = '%s' and pc.v_ComponentId = '%s' and scf.v_ComponentFieldId = '%s' "},

	"getCardio": {Q: "select p.v_FirstName, p.v_FirstLastName, p.v_SecondLastName, o.v_Name, " +
		" s.d_ServiceDate from service s " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"inner join calendar c on c.v_PersonId = p.v_PersonId " +
		"inner join protocol pr on s.v_ProtocolId = pr.v_ProtocolId and pr.i_IsDeleted = 0 " +
		"inner join organization o on pr.v_CustomerOrganizationId = o.v_OrganizationId " +
		"where  c.i_ServiceTypeId ='%s' and s.i_MasterServiceId = '%s' " +
		"group by p.v_FirstName, p.v_FirstLastName,  p.v_SecondLastName, o.v_Name, s.d_ServiceDate " +
		"order by s.d_ServiceDate "},

	"getStatusLiquid": {Q: "select i_StatusLiquidation = IIF(i_StatusLiquidation is null,0,i_StatusLiquidation)" +
		"from service where i_IsDeleted = 0 and v_ServiceId = '%s'"},
}
