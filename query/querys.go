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
		"where s.v_ServiceId = '%s' and pc.v_ComponentId = '%s' and scf.v_ComponentFieldId = '%s'"},
}
